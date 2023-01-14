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

func upperFirst(x string) string {
	f := x[0]
	if f < 0x61 || f > 0x7A {
		return x
	}
	return string(f-0x20) + x[1:]
}

type LowerFirst struct{}

func (l LowerFirst) Map(x string) string {
	return lowerFirst(x)
}

func lowerFirst(x string) string {
	f := x[0]
	if f < 0x41 || f > 0x5A {
		return x
	}
	return string(f+0x20) + x[1:]
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

type HashMap struct {
	StringMap map[string]string
}

func (m HashMap) Map(x string) string {
	if r, ok := m.StringMap[x]; ok {
		return r
	}
	return ""
}
