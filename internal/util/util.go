package util

func UpperFirst(inp string) string {
	f := inp[0]
	if f < 0x61 || f > 0x7A {
		return inp
	}
	return string(f-0x20) + inp[1:]
}

type Queue[T any] struct {
	elements []T
}

func (q *Queue[T]) Push(el T) {
	q.elements = append(q.elements, el)
}

func (q *Queue[T]) Pop() (T, bool) {
	if len(q.elements) == 0 {
		var t T
		return t, false
	}
	r := q.elements[0]
	q.elements = q.elements[1:]
	return r, true
}
