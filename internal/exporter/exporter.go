package exporter

import (
	"bytes"
	"fmt"
	"github.com/abekoh/mapc/internal/pivot"
	"go/format"
	"html/template"
	"io"
	"strings"
)

type TmplParam struct {
	Pkg        string
	ImportPkgs []string
	Funcs      []Func
}

type Func struct {
	Name     string
	FromType string
	ToType   string
	Fields   []Field
}

type Field struct {
	From string
	To   string
}

const mapperTmpl = `package {{ .Pkg }}

import (
{{- range .ImportPkgs }}
	"{{- . }}"
{{- end }}
)

{{- range .Funcs }}
func {{ .Name }}(inp {{ .FromType }}) {{ .ToType }} {
	return {{ .ToType }}{
		{{- range .Fields }}
		{{ .From }}: inp.{{ .To }},
		{{- end }}
	}
}
{{- end }}
`

func NewTmplParam(p *pivot.Pivot, dstPkgName string) TmplParam {
	fields := make([]Field, 0, len(p.Maps))
	for _, mp := range p.Maps {
		fields = append(fields, Field{
			From: mp.From.Name(),
			To:   mp.To.Name(),
		})
	}
	fc := Func{
		Name:     funcName(p),
		FromType: fmt.Sprintf("%s.%s", p.From.PkgName(), p.From.StructName()),
		ToType:   fmt.Sprintf("%s.%s", p.To.PkgName(), p.To.StructName()),
		Fields:   fields,
	}
	// FIXME
	pkgs := []string{
		"github.com/abekoh/mapc/internal/cmd/testdata/a",
		"github.com/abekoh/mapc/internal/cmd/testdata/b",
	}
	return TmplParam{
		Pkg:        dstPkgName,
		ImportPkgs: pkgs,
		Funcs:      []Func{fc},
	}
}

func funcName(p *pivot.Pivot) string {
	var b strings.Builder
	b.WriteString("To")
	b.WriteString(camelizeFirst(p.To.PkgName()))
	b.WriteString(camelizeFirst(p.To.StructName()))
	return b.String()
}

func camelizeFirst(s string) string {
	f := s[0]
	if f < 0x61 || f > 0x7A {
		return s
	}
	return string(f-0x20) + s[1:]
}

func Run(w io.Writer, param TmplParam) error {
	var buf bytes.Buffer
	tmpl := template.Must(template.New("mapper").Parse(mapperTmpl))
	if err := tmpl.Execute(&buf, param); err != nil {
		return fmt.Errorf("failed to write mapper: %w", err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format: %w", err)
	}
	if _, err := w.Write(formatted); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}
