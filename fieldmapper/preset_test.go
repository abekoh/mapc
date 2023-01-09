package fieldmapper

import (
	"fmt"
	"testing"
)

func TestIdentify_Map(t *testing.T) {
	tests := []struct {
		x    string
		want string
	}{
		{x: "a", want: "a"},
		{x: "Foo", want: "Foo"},
		{x: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s->%s", tt.x, tt.want), func(t *testing.T) {
			if got := (Identify{}).Map(tt.x); got != tt.want {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
