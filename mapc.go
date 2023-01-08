package mapc

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/abekoh/mapc/internal/code"
)

type MapC struct {
	pairs  []pair
	option Option
}

type FieldMapper func(string) string

type pair struct {
	from   any
	to     any
	option Option
}

type Option struct {
	FuncName     string
	OutPath      string
	FieldMappers []FieldMapper
}

type Generated struct {
	filePath string
	codeFile code.File
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

func (o Option) override(opts ...Option) Option {
	for _, opt := range opts {
		// TODO: to be more simple
		if len(opt.OutPath) != 0 {
			o.OutPath = opt.OutPath
		}
	}
	return o
}

func New() *MapC {
	return &MapC{
		pairs: []pair{},
		option: Option{
			FuncName:     "",
			OutPath:      "",
			FieldMappers: []FieldMapper{},
		},
	}
}

func (m *MapC) Global(opt Option) {
	m.option = opt
}

func (m *MapC) Register(from, to any, options ...*Option) {
	m.pairs = append(m.pairs, pair{from: from, to: to})
}

func (m MapC) Generate() (GeneratedList, []error) {
	return GeneratedList{}, []error{}
}
