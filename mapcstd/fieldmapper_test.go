package mapcstd

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
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s->%s", tt.x, tt.want), func(t *testing.T) {
			if got := (Identify{}).Map(tt.x); got != tt.want {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpperFirst_Map(t *testing.T) {
	tests := []struct {
		x    string
		want string
	}{
		{x: "a", want: "A"},
		{x: "foo", want: "Foo"},
		{x: "zero", want: "Zero"},
		{x: "mapFunc", want: "MapFunc"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s->%s", tt.x, tt.want), func(t *testing.T) {
			if got := (UpperFirst{}).Map(tt.x); got != tt.want {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLowerFirst_Map(t *testing.T) {
	tests := []struct {
		x    string
		want string
	}{
		{x: "A", want: "a"},
		{x: "Foo", want: "foo"},
		{x: "Zero", want: "zero"},
		{x: "MapFunc", want: "mapFunc"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s->%s", tt.x, tt.want), func(t *testing.T) {
			if got := (LowerFirst{}).Map(tt.x); got != tt.want {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnakeToCamel_Map(t *testing.T) {
	tests := []struct {
		x    string
		want string
	}{
		{x: "snake_case", want: "snakeCase"},
		{x: "Snake_case", want: "SnakeCase"},
		{x: "Snake_Case", want: "SnakeCase"},
		{x: "camelCase", want: "camelCase"},
		{x: "CamelCase", want: "CamelCase"},
		{x: "kebab-case", want: "kebab-case"},
		{x: "a", want: "a"},
		{x: "foo", want: "foo"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s->%s", tt.x, tt.want), func(t *testing.T) {
			if got := (SnakeToCamel{}).Map(tt.x); got != tt.want {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashMap_Map(t *testing.T) {
	tests := []struct {
		hashMap HashMap
		x       string
		want    string
	}{
		{
			hashMap: HashMap{
				"Foo": "Bar",
			},
			x:    "Foo",
			want: "Bar",
		},
		{
			hashMap: HashMap{
				"Foo": "Bar",
			},
			x:    "Baz",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s->%s", tt.x, tt.want), func(t *testing.T) {
			if got := tt.hashMap.Map(tt.x); got != tt.want {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
