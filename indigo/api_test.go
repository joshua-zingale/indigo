package indigo

import (
	"testing"

	"github.com/joshua-zingale/indigo/indigo/internal"
	"github.com/joshua-zingale/indigo/indigo/standard/library"
)

func TestRead(t *testing.T) {
	pairs := map[string]any{
		"123":           123,
		"9.5":           9.5,
		"(1 2 3)":       NewList(1, 2, 3),
		"Bob":           Symbol("Bob"),
		"(+ 1 (* 3 2))": NewList(Symbol("+"), 1, NewList(Symbol("*"), 3, 2)),
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

func TestEvaluate(t *testing.T) {
	pairs := map[string]any{
		"123":               123,
		"9.5":               9.5,
		"(+ 1 2 3 4)":       10,
		"(+ (+ 1 2) 3 4.5)": 10.5,
	}
	interpreter := NewStandardInterpreter()
	interpreter.LoadModule(library.IndigoCore)
	for source, expected := range pairs {
		read, err := Read(source)
		if err != nil {
			t.Error(err)
		}

		found, err := interpreter.Eval(read)
		if err != nil {
			t.Error(err)
		}

		if !internal.IndigoEqual(found, expected) {
			t.Errorf("\nfound   : %v\nexpected: %v", found, expected)
		}

	}
}
