package fieldmapper

import (
	"strings"

	"github.com/abekoh/mapc/internal/util"
)

type Identify struct{}

func (i Identify) Map(x string) string {
	return x
}

type UpperFirst struct{}

func (u UpperFirst) Map(x string) string {
	return util.UpperFirst(x)
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
			s = util.UpperFirst(s)
		}
		b.WriteString(s)
	}
	return b.String()
}
