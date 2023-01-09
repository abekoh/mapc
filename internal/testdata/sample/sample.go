package sample

import (
	"time"

	"github.com/google/uuid"
)

type Object struct {
	Bool            bool
	Int             int
	Int8            int8
	Int16           int16
	Int32           int32
	Int64           int64
	Uint            uint
	Uint8           uint8
	Uint16          uint16
	Uint32          uint32
	Uint64          uint64
	Uintptr         uintptr
	Float32         float32
	Float64         float64
	Complex64       complex64
	Complex128      complex128
	IntArray        [5]int
	IntChan         chan int
	IntToStringFunc func(string) int
	Interface       interface{}
	StringIntMap    map[string]int
	IntPointer      *int
	Slice           []int
	String          string
	EmptyStruct     struct{}
	// TODO: unsafe pointer
}

// AUser is source for mapping test
type AUser struct {
	ID           uuid.UUID
	Name         string
	Age          int
	RegisteredAt time.Time
}

// BUser is destination for mapping test
type BUser struct {
	ID           uuid.UUID
	Name         string
	Age          int
	RegisteredAt time.Time
}

// ToBUser is mapper from AUser into BUser
func ToBUser(x AUser) BUser {
	return BUser{
		ID:           x.ID,
		Name:         x.Name,
		Age:          x.Age,
		RegisteredAt: x.RegisteredAt,
	}
}
