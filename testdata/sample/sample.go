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
	ExternalType    uuid.UUID
	ExternalPointer *uuid.UUID
}

var ObjectFieldNames = []string{
	"Bool",
	"Int",
	"Int8",
	"Int16",
	"Int32",
	"Int64",
	"Uint",
	"Uint8",
	"Uint16",
	"Uint32",
	"Uint64",
	"Uintptr",
	"Float32",
	"Float64",
	"Complex64",
	"Complex128",
	"IntArray",
	"IntChan",
	"IntToStringFunc",
	"Interface",
	"StringIntMap",
	"IntPointer",
	"Slice",
	"String",
	"EmptyStruct",
	"ExternalType",
	"ExternalPointer",
}

// SrcUser is source for mapping test
type SrcUser struct {
	ID           uuid.UUID
	Name         string
	Age          int
	RegisteredAt time.Time
}

// DestUser is destination for mapping test
type DestUser struct {
	ID           uuid.UUID
	Name         string
	Age          int
	RegisteredAt time.Time
}

// ToDestUser is mapper from SrcUser into DestUser
func ToDestUser(x SrcUser) DestUser {
	return DestUser{
		ID:           x.ID,
		Name:         x.Name,
		Age:          x.Age,
		RegisteredAt: x.RegisteredAt,
	}
}
