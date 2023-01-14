package util

import (
	"strings"
)

func isLower(c byte) bool {
	if c >= 0x61 && c <= 0x7A {
		return true
	}
	return false
}

func isUpper(c byte) bool {
	if c >= 0x41 && c <= 0x5A {
		return true
	}
	return false
}

func UpperFirst(x string) string {
	f := x[0]
	if !isLower(f) {
		return x
	}
	return string(f-0x20) + x[1:]
}

func LowerFirst(x string) string {
	f := x[0]
	if !isUpper(f) {
		return x
	}
	return string(f+0x20) + x[1:]
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
	return sp[len(sp)-1]
}

func IsPrivate(name string) bool {
	if name == "" {
		return false
	}
	return isLower(name[0])
}
