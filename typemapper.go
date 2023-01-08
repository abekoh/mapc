package mapc

type Typ interface {
	PkgPath() string
	Name() string
}

type BuiltinTyp string

const (
	Int   BuiltinTyp = "int"
	Int64 BuiltinTyp = "int64"
)

func (bt BuiltinTyp) PkgPath() string {
	return ""
}

func (bt BuiltinTyp) Name() string {
	return string(bt)
}
