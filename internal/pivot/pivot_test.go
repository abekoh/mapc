package pivot

import (
	"testing"
)

type From struct {
	Int int
}

type To struct {
	Int int
}

func loadField(t *testing.T, str any, fieldName string) Var {
	t.Helper()

	s := loadStruct(t, str)
	f := s.Var(fieldName)
	if f == nil {
		t.Fatalf("field '%s' is not found", fieldName)
	}
	return Var{
		v: f,
	}
}

func Test_match(t *testing.T) {
	type args struct {
		from            *Struct
		toTokenFieldMap tokenFieldMap
	}
	tests := []struct {
		name         string
		args         args
		wantElements int
	}{
		{
			name: "int -> int",
			args: args{
				from:            loadStruct(t, From{}),
				toTokenFieldMap: tokenFieldMap{"Int": loadField(t, To{}, "Int")},
			},
			wantElements: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := match(tt.args.from, tt.args.toTokenFieldMap); len(got) != tt.wantElements {
				t.Errorf("match() = %v, wantElements = %v", got, tt.wantElements)
			}
		})
	}
}
