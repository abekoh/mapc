package fieldmapper

type FieldMapper interface {
	Map(string) string
}

var Default = []FieldMapper{
	&Identify{},
}
