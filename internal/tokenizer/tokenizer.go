package tokenizer

type Tokenizer struct {
	tokenizers []TokenizeFunc
}

func NewTokenizer(tokenizers ...TokenizeFunc) *Tokenizer {
	return &Tokenizer{tokenizers: tokenizers}
}

type TokenizeFunc func(string) string

func (t Tokenizer) Tokenize(inp string) string {
	for _, tokenize := range t.tokenizers {
		inp = tokenize(inp)
	}
	return inp
}
