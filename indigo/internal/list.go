package internal

import "github.com/joshua-zingale/indigo/indigo/interfaces"

type List []any

func (l List) Car() any {
	return l[0]
}

func (l List) Cdr() any {
	return l[1:]
}

func (l List) Empty() bool {
	return len(l) == 0
}

func (l List) IsList() {}

func NewList(elements ...any) interfaces.List {
	return List(elements)
}
