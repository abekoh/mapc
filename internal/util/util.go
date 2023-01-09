package util

import (
	"strings"
)

func UpperFirst(x string) string {
	f := x[0]
	if f < 0x61 || f > 0x7A {
		return x
	}
	return string(f-0x20) + x[1:]
}

func Prepend[T any](x []T, y T) []T {
	var zero T
	x = append(x, zero)
	copy(x[1:], x)
	x[0] = y
	return x
}

func PkgNameFromPath(pkgPath string) string {
	sp := strings.Split(pkgPath, "/")
	if len(sp) == 0 {
		return ""
	}
	return sp[len(sp)-1]
}
