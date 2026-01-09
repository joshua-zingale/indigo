package indigo

import (
	"reflect"
	"testing"
)

func TestReader(t *testing.T) {
	pairs := map[string]any{
		"123":           123,
		"9.5":           9.5,
		"(1 2 3)":       NewCons(1, NewCons(2, NewCons(3, nil))),
		"Bob":           Symbol("Bob"),
		"(+ 1 (* 3 2))": NewCons(Symbol("+"), NewCons(1, NewCons(NewCons(Symbol("*"), NewCons(3, NewCons(2, nil))), nil))),
	}

	for source, expected := range pairs {
		found, err := Read(source)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(found, expected) {
			t.Errorf("\nfound   : %v\nexpected: %v", found, expected)
		}

	}
}
