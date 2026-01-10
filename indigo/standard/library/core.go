package library

import (
	"fmt"
	"reflect"

	"github.com/joshua-zingale/indigo/indigo/functions"
	"github.com/joshua-zingale/indigo/indigo/functools"
	"github.com/joshua-zingale/indigo/indigo/interfaces"
	"github.com/joshua-zingale/indigo/indigo/internal"
)

var IndigoCore interfaces.IndigoModule = internal.NewModule("indigo-core", internal.NewNameSpaceFromMap(map[string]any{
	"+": functions.NewIndigoFunctionFromGo(add),
}))

func add(evaluator interfaces.IndigoEvaluator, namespace interfaces.NameSpace, args interfaces.List) (any, error) {
	argSlice := internal.ListToSlice(args)
	argSliceEvaluated, err := functools.MapShortCircuit(func(a any) (any, error) {
		return evaluator.Eval(a, namespace)
	}, argSlice)
	if err != nil {
		return nil, err
	}
	if xs, err := functools.MapShortCircuit(func(a any) (int, error) {
		v, ok := a.(int)
		if !ok {
			return 0, fmt.Errorf("")
		}
		return v, nil
	}, argSliceEvaluated); err == nil {
		return functools.Reduce(func(v int, x int) int {
			return v + x
		}, 0, xs), nil
	}

	floats, err := functools.MapWithError(func(a any) (float64, error) {
		float, err := internal.ValidateType(reflect.TypeFor[float64](), a)
		if err != nil {
			return 0, err
		}
		return float.(float64), nil
	}, argSliceEvaluated)
	if err != nil {
		return 0, err
	}

	return functools.Reduce(func(v float64, x float64) float64 {
		return v + x
	}, 0, floats), nil
}
