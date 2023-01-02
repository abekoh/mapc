package tokenizer

import "testing"

func TestTokenizer_Tokenize(t1 *testing.T) {
	type fields struct {
		tokenizers []TokenizeFunc
	}
	type args struct {
		inp string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "no tokenizer",
			fields: fields{
				tokenizers: []TokenizeFunc{},
			},
			args: args{
				inp: "input",
			},
			want: "input",
		},
		{
			name: "upperFirst",
			fields: fields{
				tokenizers: []TokenizeFunc{
					UpperFirst,
				},
			},
			args: args{
				inp: "input",
			},
			want: "Input",
		},
		{
			name: "snakeToCamel",
			fields: fields{
				tokenizers: []TokenizeFunc{
					SnakeToCamel,
				},
			},
			args: args{
				inp: "snake_case",
			},
			want: "snakeCase",
		},
		{
			name: "snakeToCamel, upperFirst",
			fields: fields{
				tokenizers: []TokenizeFunc{
					SnakeToCamel,
					UpperFirst,
				},
			},
			args: args{
				inp: "snake_case",
			},
			want: "SnakeCase",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := Tokenizer{
				tokenizers: tt.fields.tokenizers,
			}
			if got := t.Tokenize(tt.args.inp); got != tt.want {
				t1.Errorf("Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
