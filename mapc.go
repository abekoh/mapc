package mapc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/abekoh/mapc/fieldmapper"
	"github.com/abekoh/mapc/internal/code"
	"github.com/abekoh/mapc/internal/mapping"
	"github.com/abekoh/mapc/typemapper"
)

type MapC struct {
	Option *Option
	inputs []input
}

type Registerer interface {
	Register(from, to any, optFns ...func(option *Option))
	Group(optFns ...func(*Option)) *Group
}

func New() *MapC {
	return &MapC{
		inputs: []input{},
		Option: &Option{
			WithFuncComment: true,
			FieldMappers:    fieldmapper.DefaultMappers,
			TypeMappers:     typemapper.Defaults,
		},
	}
}

func (m *MapC) Register(from, to any, optFns ...func(*Option)) {
	m.inputs = append(m.inputs, input{
		from:   from,
		to:     to,
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
			errs = append(errs, fmt.Errorf("option is nil. from=%T, to=%T", input.from, input.to))
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
		mp, err := mapper.NewMapping(input.from, input.to, pkgPath)
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
			WithFuncComment: input.option.WithFuncComment,
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
	from   any
	to     any
	option *Option
}

type Option struct {
	OutPath         string
	OutPkgPath      string
	WithFuncComment bool
	FieldMappers    []fieldmapper.FieldMapper
	TypeMappers     []typemapper.TypeMapper
}

func (o *Option) copy() *Option {
	fms := make([]fieldmapper.FieldMapper, len(o.FieldMappers))
	copy(fms, o.FieldMappers)
	tms := make([]typemapper.TypeMapper, len(o.TypeMappers))
	copy(tms, o.TypeMappers)
	return &Option{
		OutPath:         o.OutPath,
		OutPkgPath:      o.OutPkgPath,
		WithFuncComment: o.WithFuncComment,
		FieldMappers:    fms,
		TypeMappers:     tms,
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

func (g *Group) Register(from, to any, optFns ...func(*Option)) {
	g.MapC.inputs = append(g.MapC.inputs, input{
		from:   from,
		to:     to,
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
