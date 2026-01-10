package evaluation

import (
	"testing"

	"github.com/joshua-zingale/indigo/indigo/interfaces"
	"github.com/joshua-zingale/indigo/indigo/internal"
	"github.com/joshua-zingale/indigo/indigo/reading"
)

func TestStandardEvaluator(t *testing.T) {
	pairs := map[string]any{
		"123":                    123,
		"9.5":                    9.5,
		"(+ 2 3)":                5.0,
		"(+ 2.5 3)":              5.5,
		"Bob":                    "WOW!",
		"(if true -3 2)":         -3,
		"(if false -3 2)":        2,
		"(if false undefined 2)": 2,
	}

	evaluator := NewStandardEvaluator()
	namespace := internal.NewNameSpace()
	namespace.Set(interfaces.Symbol("+"), NewTypeCheckedIndigoFunctionFromGo(func(a float64, b float64) (float64, error) {
		return a + b, nil
	}))
	namespace.Set(interfaces.Symbol("*"), NewTypeCheckedIndigoFunctionFromGo(func(a float64, b float64) (float64, error) {
		return a * b, nil
	}))
	namespace.Set(interfaces.Symbol("if"), NewIndigoFunctionFromGo(func(evaluator interfaces.IndigoEvaluator, namespace interfaces.NameSpace, args interfaces.List) (any, error) {
		argSlice := internal.ListToSlice(args)
		if len(argSlice) != 3 {
			panic("invalid num args for if")
		}
		condition, err := evaluator.Eval(argSlice[0], namespace)
		if err != nil {
			panic(err)
		}
		veracity := condition.(bool)
		if veracity {
			return evaluator.Eval(argSlice[1], namespace)
		} else {
			return evaluator.Eval(argSlice[2], namespace)
		}
	}))
	namespace.Set(interfaces.Symbol("Bob"), "WOW!")
	namespace.Set(interfaces.Symbol("true"), true)
	namespace.Set(interfaces.Symbol("false"), false)

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
