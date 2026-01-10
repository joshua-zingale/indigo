package internal

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
