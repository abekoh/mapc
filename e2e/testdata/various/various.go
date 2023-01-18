package various

import "github.com/google/uuid"

type S struct {
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

// S2 is same with S
type S2 struct {
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

type SWithPointers struct {
	Bool            *bool
	Int             *int
	Int8            *int8
	Int16           *int16
	Int32           *int32
	Int64           *int64
	Uint            *uint
	Uint8           *uint8
	Uint16          *uint16
	Uint32          *uint32
	Uint64          *uint64
	Uintptr         *uintptr
	Float32         *float32
	Float64         *float64
	Complex64       *complex64
	Complex128      *complex128
	IntArray        *[5]int
	IntChan         *chan int
	IntToStringFunc *func(string) int
	Interface       *interface{}
	StringIntMap    *map[string]int
	IntPointer      **int
	Slice           *[]int
	String          *string
	EmptyStruct     *struct{}
	ExternalType    *uuid.UUID
	ExternalPointer **uuid.UUID
}
