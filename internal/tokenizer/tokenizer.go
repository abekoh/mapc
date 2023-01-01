package tokenizer

import "sync"

var (
	mu         sync.RWMutex
	tokenizers []TokenizeFunc
)

type TokenizeFunc func(string) string

func Initialize(funcs ...TokenizeFunc) {
	mu.Lock()
	defer mu.Unlock()
	tokenizers = funcs
}

func Tokenize(inp string) string {
	for _, tokenize := range tokenizers {
		inp = tokenize(inp)
	}
	return inp
}
