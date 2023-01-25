package mapcstd

type Caster interface {
	Caller() *Caller
}

type Caller struct {
	PkgPath    string
	Name       string
	CallerType CallerType
}

type NopCaster struct{}

func (n NopCaster) Caller() *Caller {
	return nil
}

type CallerType int

const (
	Unary CallerType = iota
	Typ
	Func
)

type SimpleCaster struct {
	caller *Caller
}

func (s SimpleCaster) Caller() *Caller {
	return s.caller
}
