package typemapper

type Caster interface {
	PkgPath() string
	Fun() string
}

type NopCaster struct{}

func (n NopCaster) PkgPath() string {
	return ""
}

func (n NopCaster) Fun() string {
	return ""
}

type SimpleCaster struct {
	pkgPath string
	fun     string
}

func (s SimpleCaster) PkgPath() string {
	return s.pkgPath
}

func (s SimpleCaster) Fun() string {
	return s.fun
}
