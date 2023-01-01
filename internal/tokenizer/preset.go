package tokenizer

import (
	"strings"
)

func UpperFirst(inp string) string {
	f := inp[0]
	if f < 0x61 || f > 0x7A {
		return inp
	}
	return string(f-0x20) + inp[1:]
}

func SnakeToCamel(inp string) string {
	var b strings.Builder
	sp := strings.Split(inp, "_")
	if len(sp) == 1 {
		return inp
	}
	for i, s := range sp {
		if i != 0 {
			s = UpperFirst(s)
		}
		b.WriteString(s)
	}
	return b.String()
}
