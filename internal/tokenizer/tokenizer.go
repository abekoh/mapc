package tokenizer

import "sync"

var (
	mu         sync.RWMutex
	tokenizers []TokenizeFunc
)

type TokenizeFunc func(string) string

func Register(inps ...TokenizeFunc) {
	mu.Lock()
	defer mu.Unlock()
	for _, inp := range inps {
		tokenizers = append(tokenizers, inp)
	}
}

func Tokenize(inp string) string {
	for _, tokenize := range tokenizers {
		inp = tokenize(inp)
	}
	return inp
}
