package evaluation

import (
	"testing"

	reading "github.com/joshua-zingale/indigo/indigo/Reading"
	"github.com/joshua-zingale/indigo/indigo/interfaces"
	"github.com/joshua-zingale/indigo/indigo/internal"
)

func TestStandardEvaluator(t *testing.T) {
	pairs := map[string]any{
		"123":           123,
		"9.5":           9.5,
		"(+ 2 3)":       5,
		"Bob":           "WOW!",
		"(+ 1 (* 3 2))": 7,
	}

	evaluator := NewStandardEvaluator()
	namespace := internal.NewNameSpace()
	namespace.Set(interfaces.Symbol("+"), NewIndigoFunctionFromGoFunction(func(a int, b int) (int, error) {
		return a + b, nil
	}))
	namespace.Set(interfaces.Symbol("*"), NewIndigoFunctionFromGoFunction(func(a int, b int) (int, error) {
		return a * b, nil
	}))
	namespace.Set(interfaces.Symbol("Bob"), "WOW!")

	for source, expected := range pairs {

		syntax, err := reading.NewStandardReader(source).Read()
		if err != nil {
			t.Error(err)
		}

		found, err := evaluator.Eval(syntax, namespace)
		if err != nil {
			t.Error(err)
		}

		if !internal.IndigoEqual(found, expected) {
			t.Errorf("\nfound   : %v\nexpected: %v", found, expected)
		}

	}
}
