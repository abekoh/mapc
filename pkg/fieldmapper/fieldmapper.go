package fieldmapper

import (
	"strings"
)

type FieldMapper func(string) string

var Default = []FieldMapper{
	Identify,
}

func Identify(x string) string {
	return x
}

func UpperFirst(x string) string {
	f := x[0]
	if f < 0x61 || f > 0x7A {
		return x
	}
	return string(f-0x20) + x[1:]
}

func SnakeToCamel(x string) string {
	var b strings.Builder
	sp := strings.Split(x, "_")
	if len(sp) == 1 {
		return x
	}
	for i, s := range sp {
		if i != 0 {
			s = UpperFirst(s)
		}
		b.WriteString(s)
	}
	return b.String()
}
