package model

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"strings"
)

type TmplParam struct {
	Package        string
	ImportPackages []string
	Funcs          []Func
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

const mapperTmpl = `package {{ .Package }}

import (
{{ range .ImportPackages }}
	"{{ . }}"
{{ end }}
)

{{ range .Funcs }}
func {{ .Name }}(inp {{ .FromType }}) {{ .ToType }} {
	return {{ .ToType }}{
		{{ range .Fields }}
		{{ .From }}: inp.{{ .To }},
		{{ end }}
	}
}
{{ end }}
`

func NewTmplParam(m *Model, dstPkgName string) TmplParam {
	fields := make([]Field, 0, len(m.maps))
	for _, mp := range m.maps {
		fields = append(fields, Field{
			From: mp.from.Name(),
			To:   mp.to.Name(),
		})
	}
	fc := Func{
		Name:     funcName(m),
		FromType: fmt.Sprintf("%s.%s", m.from.param.Pkg, m.from.param.Struct),
		ToType:   fmt.Sprintf("%s.%s", m.to.param.Pkg, m.to.param.Struct),
		Fields:   fields,
	}
	// FIXME
	pkgs := []string{
		"example.com/hoge",
	}
	return TmplParam{
		Package:        dstPkgName,
		ImportPackages: pkgs,
		Funcs:          []Func{fc},
	}
}

func funcName(m *Model) string {
	var b strings.Builder
	b.WriteString("to")
	b.WriteString(camelize(m.to.param.Pkg))
	b.WriteString(camelize(m.to.param.Struct))
	return b.String()
}

func camelize(s string) string {
	f := s[0]
	if f < 0x61 || f > 0x7A {
		return s
	}
	return string(f+0x20) + s[1:]
}

func Write(w io.Writer, param TmplParam) error {
	var buf bytes.Buffer
	tmpl := template.Must(template.New("mapper").Parse(mapperTmpl))
	if err := tmpl.Execute(w, param); err != nil {
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
