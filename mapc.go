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

func (m *MapC) Generate() (res GeneratedList, errs []error) {
	for _, input := range m.inputs {
		if input.option == nil {
			errs = append(errs, fmt.Errorf("option is nil. src=%T, dest=%T", input.src, input.dest))
			continue
		}
		mapper := mapping.Mapper{
			FieldMappers: input.option.FieldMappers,
			TypeMappers:  input.option.TypeMappers,
		}
		pkgPath := input.option.OutPkgPath
		if pkgPath == "" {
			pkgPath = pkgPathFromRelativePath(input.option.OutPath)
		}
		mp, err := mapper.NewMapping(input.src, input.dest, pkgPath)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to map: %w", err))
			continue
		}
		// TODO: cache file
		var f *code.File
		if existed, err := code.LoadFile(input.option.OutPath, pkgPath); err == nil {
			f = existed
		} else {
			f = code.NewFile(pkgPath)
		}
		fn := code.NewFuncFromMapping(mp, &code.FuncOption{
			Name:                 input.option.FuncName,
			FuncComment:          input.option.FuncComment,
			NoMapperFieldComment: input.option.NoMapperFieldComment,
		})
		if existedFn, ok := f.FindFunc(fn.Name()); ok {
			if err := fn.AppendNotSetExprs(existedFn); err != nil {
				errs = append(errs, fmt.Errorf("failed to append not set exprs: %w", err))
				continue
			}
		}
		err = f.Attach(fn)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to apply: %w", err))
			continue
		}
		g := &Generated{
			filePath: input.option.OutPath,
			codeFile: f,
		}
		res = append(res, g)
	}
	return
}

type input struct {
	src    any
	dest   any
	option *Option
}

type Option struct {
	OutPath              string
	OutPkgPath           string
	FuncName             string
	FuncComment          bool
	NoMapperFieldComment bool
	FieldMappers         []mapcstd.FieldMapper
	TypeMappers          []mapcstd.TypeMapper
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

type GeneratedList []*Generated

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
