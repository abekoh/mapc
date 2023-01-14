package mapc

import (
	"bytes"
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

func New() *MapC {
	return &MapC{
		inputs: []input{},
		Option: &Option{
			FieldMappers: fieldmapper.Default,
			TypeMappers:  typemapper.Default,
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
	for _, pair := range m.inputs {
		if pair.option == nil {
			errs = append(errs, fmt.Errorf("option is nil. from=%T, to=%T", pair.from, pair.to))
			continue
		}
		// TODO: validate with Option's method
		if pair.option.OutPath == "" {
			errs = append(errs, fmt.Errorf("OutPath is empty"))
			continue
		}
		mapper := mapping.Mapper{
			FieldMappers: pair.option.FieldMappers,
			TypeMappers:  pair.option.TypeMappers,
		}
		mp, err := mapper.NewMapping(pair.from, pair.to)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to map: %w", err))
			continue
		}
		pkgPath := pkgPathFromRelativePath(pair.option.OutPath)
		// TODO: cache file
		var f *code.File
		if existed, err := code.LoadFile(pair.option.OutPath, pkgPath); err == nil {
			f = existed
		} else {
			f = code.NewFile(pkgPath)
		}
		fn := code.NewFuncFromMapping(mp, &code.FuncOption{})
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
			filePath: pair.option.OutPath,
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
	OutPath      string
	FieldMappers []fieldmapper.FieldMapper
	TypeMappers  []typemapper.TypeMapper
}

func (o *Option) copy() *Option {
	fms := make([]fieldmapper.FieldMapper, len(o.FieldMappers))
	copy(fms, o.FieldMappers)
	tms := make([]typemapper.TypeMapper, len(o.TypeMappers))
	copy(tms, o.TypeMappers)
	return &Option{
		OutPath:      o.OutPath,
		FieldMappers: fms,
		TypeMappers:  tms,
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

func (g Generated) Save() error {
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
