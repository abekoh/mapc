package mapcstd

type Caster interface {
	Name() string
	PkgPath() string
	Type() CasterType
	ReturnType() ReturnType
}

type CasterType int

const (
	Nop CasterType = iota
	Caller
	Method
)

type NopCaster struct{}

func (n NopCaster) Name() string {
	return ""
}

func (n NopCaster) PkgPath() string {
	return ""
}

func (n NopCaster) Type() CasterType {
	return Nop
}

func (n NopCaster) ReturnType() ReturnType {
	return OnlyValue
}

type CallerCaster struct {
	pkgPath    string
	name       string
	callerType CallerType
	retType    ReturnType
}

type CallerType int

const (
	Unary CallerType = iota
	Type
	Func
)

func (c CallerCaster) Name() string {
	return c.name
}

func (c CallerCaster) PkgPath() string {
	return c.pkgPath
}

func (c CallerCaster) Type() CasterType {
	return Caller
}

func (c CallerCaster) ReturnType() ReturnType {
	return c.retType
}
