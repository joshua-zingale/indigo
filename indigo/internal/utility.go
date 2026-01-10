package internal

import (
	"fmt"

	"github.com/joshua-zingale/indigo/indigo/interfaces"
)

func ListToSlice(l interfaces.List) []any {
	switch slice := l.(type) {
	case List:
		return slice
	}
	return listToSlice(l)
}

func listToSlice(l interfaces.List) []any {
	var slice []any

	for !l.Empty() {
		slice = append(slice, l.Car())
		l = l.Cdr().(interfaces.List)
	}
	return slice
}

type VerifiedList struct {
	interfaces.Cons
}

func (vl *VerifiedList) IsList() {}

func ValidateList(l interfaces.Cons) (interfaces.List, error) {
	for !l.Empty() {
		next := l.Cdr()

		if cons, ok := next.(interfaces.Cons); ok {
			l = cons
		} else if next == nil {
			l = nil
		} else {
			return nil, fmt.Errorf("Cons is not a List")
		}
	}

	return &VerifiedList{Cons: l}, nil
}
