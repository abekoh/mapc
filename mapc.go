package mapc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/abekoh/mapc/internal/code"
	"github.com/abekoh/mapc/internal/mapping"
	"github.com/abekoh/mapc/mapcstd"
)

type MapC struct {
	Option *Option
	inputs []input
}

type Registerer interface {
	Register(src, dest any, optFns ...func(option *Option))
	Group(optFns ...func(*Option)) *Group
}

func New() *MapC {
	return &MapC{
		inputs: []input{},
		Option: &Option{
			FuncComment:          true,
			NoMapperFieldComment: true,
			Mode:                 PrioritizeGenerated,
			FieldMappers:         mapcstd.DefaultFieldMappers,
			TypeMappers:          mapcstd.DefaultTypeMappers,
		},
	}
}

func (m *MapC) Register(src, dest any, optFns ...func(*Option)) {
	m.inputs = append(m.inputs, input{
		src:    src,
		dest:   dest,
		option: (m.Option).overwrite(optFns...),
	})
}

func (m *MapC) Group(optFns ...func(*Option)) *Group {
	return &Group{
		MapC:   m,
		parent: nil,
		Option: (m.Option).overwrite(optFns...),
	}
}

func (m *MapC) Generate() (res []*Generated, errs []error) {
	store := NewGeneratedStore()
	for _, input := range m.inputs {
		pkgPath := input.option.OutPkgPath
		if pkgPath == "" {
			pkgPath = pkgPathFromRelativePath(input.option.OutPath)
		}
		generated := store.Get(input.option.OutPath, pkgPath)

		mapper := mapping.Mapper{
			FieldMappers: input.option.FieldMappers,
			TypeMappers:  input.option.TypeMappers,
		}

		mp, err := mapper.NewMapping(input.src, input.dest, pkgPath)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to map: %w", err))
			continue
		}

		fn := code.NewFuncFromMapping(mp, &code.FuncOption{
			Name:                 input.option.FuncName,
			FuncComment:          input.option.FuncComment,
			NoMapperFieldComment: input.option.NoMapperFieldComment,
			Editable:             input.option.Mode.Editable(),
		})
		//if existedFn, ok := generated.codeFile.FindFunc(fn.Name()); ok {
		//	if err := fn.AppendNotSetExprs(existedFn); err != nil {
		//		errs = append(errs, fmt.Errorf("failed to append not set exprs: %w", err))
		//		continue
		//	}
		//}
		err = generated.codeFile.Attach(fn, code.Mode(input.option.Mode))
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to apply: %w", err))
			continue
		}

		store.Set(input.option.OutPath, generated)
	}
	res = store.List()
	return
}

type input struct {
	src    any
	dest   any
	option *Option
}

type Option struct {
	Mode                 Mode
	OutPath              string
	OutPkgPath           string
	FuncName             string
	FuncComment          bool
	NoMapperFieldComment bool
	FieldMappers         []mapcstd.FieldMapper
	TypeMappers          []mapcstd.TypeMapper
}

type Mode int

const (
	PrioritizeGenerated Mode = iota
	PrioritizeExisted
	Deterministic
)

func (m Mode) Editable() bool {
	switch m {
	case Deterministic:
		return false
	default:
		return true
	}
}

func (o *Option) copy() *Option {
	fms := make([]mapcstd.FieldMapper, len(o.FieldMappers))
	copy(fms, o.FieldMappers)
	tms := make([]mapcstd.TypeMapper, len(o.TypeMappers))
	copy(tms, o.TypeMappers)
	return &Option{
		OutPath:              o.OutPath,
		OutPkgPath:           o.OutPkgPath,
		FuncName:             o.FuncName,
		FuncComment:          o.FuncComment,
		NoMapperFieldComment: o.NoMapperFieldComment,
		Mode:                 o.Mode,
		FieldMappers:         fms,
		TypeMappers:          tms,
	}
}

func (o *Option) overwrite(optFns ...func(*Option)) *Option {
	res := o.copy()
	for _, fn := range optFns {
		fn(res)
	}
	return res
}

type Group struct {
	*MapC
	Option *Option
	parent *Group
}

func (g *Group) Group(optFns ...func(*Option)) *Group {
	return &Group{
		MapC:   g.MapC,
		parent: g,
		Option: (g.Option).overwrite(optFns...),
	}
}

func (g *Group) Register(src, dest any, optFns ...func(*Option)) {
	g.MapC.inputs = append(g.MapC.inputs, input{
		src:    src,
		dest:   dest,
		option: (g.Option).overwrite(optFns...),
	})
}

type Generated struct {
	filePath string
	codeFile *code.File
}

func (g Generated) Write(w io.Writer) error {
	return g.codeFile.Write(w)
}

func (g Generated) Sprint() (string, error) {
	var buf bytes.Buffer
	err := g.Write(&buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (g Generated) Save() error {
	if g.filePath == "" {
		return errors.New("filepath must not be empty")
	}
	f, err := os.OpenFile(g.filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	defer f.Close()
	var buf bytes.Buffer
	if err := g.codeFile.Write(&buf); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	return nil
}

// GeneratedStore stores Generated elements.
// tempMap has elements that have OutPath. Generated elements are shared with same OutPath value.
// tempList has elements that doesn't have OutPath.
type GeneratedStore struct {
	tempMap  map[string]*Generated
	tempList []*Generated
}

func NewGeneratedStore() *GeneratedStore {
	return &GeneratedStore{
		tempMap:  make(map[string]*Generated),
		tempList: []*Generated{},
	}
}

func (gs *GeneratedStore) Get(outPath, pkgPath string) *Generated {
	if g, ok := gs.tempMap[outPath]; ok {
		return g
	} else {
		g := &Generated{filePath: outPath}
		if existed, err := code.LoadFile(outPath, pkgPath); err == nil {
			g.codeFile = existed
		} else {
			g.codeFile = code.NewFile(pkgPath)
		}
		return g
	}
}

func (gs *GeneratedStore) Set(outPath string, g *Generated) {
	if outPath == "" {
		gs.tempList = append(gs.tempList, g)
	} else {
		gs.tempMap[outPath] = g
	}
}

func (gs *GeneratedStore) List() []*Generated {
	res := make([]*Generated, len(gs.tempList), len(gs.tempList)+len(gs.tempMap))
	copy(res, gs.tempList)
	for _, g := range gs.tempMap {
		res = append(res, g)
	}
	return res
}
