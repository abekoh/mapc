package fieldmapper

import (
	"strings"
)

type Identify struct{}

func (i Identify) Map(x string) string {
	return x
}

type UpperFirst struct{}

func (u UpperFirst) Map(x string) string {
	return upperFirst(x)
}

func upperFirst(inp string) string {
	f := inp[0]
	if f < 0x61 || f > 0x7A {
		return inp
	}
	return string(f-0x20) + inp[1:]
}

type SnakeToCamel struct{}

func (s SnakeToCamel) Map(x string) string {
	var b strings.Builder
	sp := strings.Split(x, "_")
	if len(sp) == 1 {
		return x
	}
	for i, s := range sp {
		if i != 0 {
			s = upperFirst(s)
		}
		b.WriteString(s)
	}
	return b.String()
}
