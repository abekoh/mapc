package fieldmapper

type FieldMapper interface {
	Map(string) string
}

var DefaultMappers = []FieldMapper{
	&Identify{},
}
