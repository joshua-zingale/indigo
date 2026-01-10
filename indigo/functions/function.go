package functions

import (
	"reflect"

	"github.com/joshua-zingale/indigo/indigo/functools"
	"github.com/joshua-zingale/indigo/indigo/interfaces"
	"github.com/joshua-zingale/indigo/indigo/internal"
)

type goFunction struct {
	function func(evaluator interfaces.IndigoEvaluator, namespace interfaces.NameSpace, args interfaces.List) (any, error)
}

func NewIndigoFunctionFromGo(function func(evaluator interfaces.IndigoEvaluator, namespace interfaces.NameSpace, args interfaces.List) (any, error)) interfaces.IndigoFunction {
	return &goFunction{
		function: function,
	}
}

// Input must be a function
func NewTypeCheckedIndigoFunctionFromGo(function any) interfaces.IndigoFunction {
	if reflect.ValueOf(function).Kind() != reflect.Func {
		panic("must be function")
	}

	parameterTypes := getFuncParameterTypes(function)

	return &goFunction{
		function: func(evaluator interfaces.IndigoEvaluator, namespace interfaces.NameSpace, args interfaces.List) (any, error) {
			evaluatedArgs, err := functools.MapWithError(func(a any) (any, error) {
				return evaluator.Eval(a, namespace)
			}, internal.ListToSlice(args))
			if err != nil {
				return nil, err
			}

			validatedArgs, err := internal.ValidateFunctionArgs(parameterTypes, evaluatedArgs)
			if err != nil {
				return nil, err
			}
			return invoke(function, validatedArgs)
		},
	}
}

func (gf *goFunction) Call(evaluator interfaces.IndigoEvaluator, namespace interfaces.NameSpace, args interfaces.List) (any, error) {
	return gf.function(evaluator, namespace, args)
}
