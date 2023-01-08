package mapc

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/abekoh/mapc/internal/code"
	"github.com/abekoh/mapc/internal/code/function"
	"github.com/abekoh/mapc/internal/mapping"
	"github.com/abekoh/mapc/internal/util"
	"github.com/abekoh/mapc/pkg/fieldmapper"
	"github.com/abekoh/mapc/pkg/typemapper"
)

type MapC struct {
	pairs  []pair
	option *Option
}

func New() *MapC {
	return &MapC{
		pairs: []pair{},
		option: &Option{
			FieldMappers: fieldmapper.Default,
			TypeMappers:  typemapper.Default,
		},
	}
}

func (m *MapC) Register(from, to any, options ...*Option) {
	m.pairs = append(m.pairs, pair{
		from:   from,
		to:     to,
		option: (m.option).override(options...),
	})
}

func (m MapC) NewGroup() *Group {
	return &Group{
		MapC:   &m,
		parent: nil,
		option: nil,
	}
}

func (m MapC) Generate() (res GeneratedList, errs []error) {
	for _, pair := range m.pairs {
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
		mapping, err := mapper.NewMapping(pair.from, pair.to)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to map: %w", err))
			continue
		}
		pkgPath := util.PkgPathFromFilePath(pair.option.OutPath)
		// TODO: cache file
		f := code.NewFile(pkgPath)
		fn := function.NewFromMapping(mapping)
		err = f.Apply(fn)
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

type pair struct {
	from   any
	to     any
	option *Option
}

type Option struct {
	OutPath      string
	FieldMappers []fieldmapper.FieldMapper
	TypeMappers  []typemapper.TypeMapper
}

func (o Option) override(opts ...*Option) *Option {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		// TODO: simplify
		if len(opt.OutPath) != 0 {
			o.OutPath = opt.OutPath
		}
		o.FieldMappers = append(o.FieldMappers, opt.FieldMappers...)
		o.TypeMappers = append(o.TypeMappers, opt.TypeMappers...)
	}
	return &o
}

type Group struct {
	*MapC
	parent *Group
	option *Option
}

func (g Group) NewGroup() *Group {
	return &Group{
		MapC:   g.MapC,
		parent: &g,
		option: &Option{},
	}
}

func (g *Group) Register(from, to any, options ...*Option) {
	opts := g.extendOptions()
	opts = append(opts, options...)
	g.MapC.pairs = append(g.MapC.pairs, pair{
		from:   from,
		to:     to,
		option: (&Option{}).override(opts...),
	})
}

func (g Group) extendOptions() []*Option {
	var opts []*Option
	targetG := &g
	for {
		if targetG == nil {
			break
		}
		util.Prepend(opts, targetG.option)
		targetG = targetG.parent
	}
	util.Prepend(opts, g.MapC.option)
	return opts
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
