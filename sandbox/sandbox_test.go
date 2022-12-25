package sandbox

import (
	"fmt"
	"reflect"
	"testing"
)

type Hoge struct {
}

func TestSand(t *testing.T) {
	h := Hoge{}
	fmt.Println(reflect.TypeOf(h).Kind())
}
