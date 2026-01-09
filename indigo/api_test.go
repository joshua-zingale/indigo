package indigo

import (
	"testing"

	"github.com/joshua-zingale/indigo/indigo/internal"
)

func TestReader(t *testing.T) {
	pairs := map[string]any{
		// "123":     123,
		// "9.5": 9.5,
		"(1 2 3)": NewList(1, 2, 3),
		// "Bob":     Symbol("Bob"),
		// "(+ 1 (* 3 2))": NewList(Symbol("+"), 1, NewList(Symbol("*"), 3, 2)),
	}

	for source, expected := range pairs {
		found, err := Read(source)
		if err != nil {
			t.Error(err)
		}
		if !internal.IndigoEqual(found, expected) {
			t.Errorf("\nfound   : %v\nexpected: %v", found, expected)
		}

	}
}
