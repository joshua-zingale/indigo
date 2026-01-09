package internal

import "github.com/joshua-zingale/indigo/indigo/interfaces"

type Cons struct {
	car any
	cdr any
}

func (c *Cons) Car() any {
	return c.car
}

func (c *Cons) Cdr() any {
	return c.cdr
}

func (c *Cons) Empty() bool {
	return c == nil
}

func NewCons(car any, cdr any) interfaces.Cons {
	return &Cons{car: car, cdr: cdr}
}

func NewList(elements ...any) interfaces.Cons {
	var head *Cons = nil
	for i := len(elements) - 1; i >= 0; i -= 1 {
		head = &Cons{car: elements[i], cdr: head}
	}
	return head
}
