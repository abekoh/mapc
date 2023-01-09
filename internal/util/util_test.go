package util

import (
	"fmt"
	"testing"
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
