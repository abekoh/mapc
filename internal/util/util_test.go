package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpperFirst(t *testing.T) {
	tests := []struct {
		x    string
		want string
	}{
		{x: "apple", want: "Apple"},
		{x: "Apple", want: "Apple"},
		{x: "zoo", want: "Zoo"},
		{x: "Zoo", want: "Zoo"},
		{x: "12345", want: "12345"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.x, tt.want), func(t *testing.T) {
			if got := UpperFirst(tt.x); got != tt.want {
				t.Errorf("UpperFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLowerFirst(t *testing.T) {
	tests := []struct {
		x    string
		want string
	}{
		{x: "Apple", want: "apple"},
		{x: "apple", want: "apple"},
		{x: "Zoo", want: "zoo"},
		{x: "zoo", want: "zoo"},
		{x: "12345", want: "12345"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.x, tt.want), func(t *testing.T) {
			if got := LowerFirst(tt.x); got != tt.want {
				t.Errorf("UpperFirst() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestPrepend(t *testing.T) {
	t.Run("[]string", func(t *testing.T) {
		strs := []string{"bar", "baz"}
		got := Prepend(strs, "foo")
		assert.Equal(t, []string{"foo", "bar", "baz"}, got)
	})
	t.Run("[]int", func(t *testing.T) {
		ints := []int{2, 3}
		got := Prepend(ints, 1)
		assert.Equal(t, []int{1, 2, 3}, got)
	})
}

func TestPkgNameFromPath(t *testing.T) {
	tests := []struct {
		pkgPath string
		want    string
	}{
		{pkgPath: "github.com/abekoh/sqlc", want: "sqlc"},
		{pkgPath: "github.com/abekoh/sqlc.test", want: "sqlc.test"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s -> %s", tt.pkgPath, tt.want), func(t *testing.T) {
			assert.Equalf(t, tt.want, PkgNameFromPath(tt.pkgPath), "PkgNameFromPath(%v)", tt.pkgPath)
		})
	}
}
