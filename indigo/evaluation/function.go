package evaluation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/joshua-zingale/indigo/indigo/functools"
	"github.com/joshua-zingale/indigo/indigo/interfaces"
)

type goFunction struct {
	function func(evaluator interfaces.IndigoEvaluator, namespace interfaces.NameSpace, args []any) (any, error)
}

func NewIndigoFunctionFromGo(function func(evaluator interfaces.IndigoEvaluator, namespace interfaces.NameSpace, args []any) (any, error)) interfaces.IndigoFunction {
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
		function: func(evaluator interfaces.IndigoEvaluator, namespace interfaces.NameSpace, args []any) (any, error) {
			evaluatedArgs, err := functools.MapWithError(func(a any) (any, error) {
				return evaluator.Eval(a, namespace)
			}, args)
			if err != nil {
				return nil, err
			}

			if err := validateFunctionArgs(parameterTypes, evaluatedArgs); err != nil {
				return nil, err
			}
			return invoke(function, evaluatedArgs)
		},
	}
}

func (gf *goFunction) Call(evaluator interfaces.IndigoEvaluator, namespace interfaces.NameSpace, args []any) (any, error) {
	return gf.function(evaluator, namespace, args)
}

func validateFunctionArgs(parameterTypes []reflect.Type, args []any) error {
	if len(args) != len(parameterTypes) {
		return fmt.Errorf("expected %d arguments but found %d", len(parameterTypes), len(args))
	}

	var typeErrors []string
	for i := range args {
		argType := reflect.TypeOf(args[i])
		paramType := parameterTypes[i]

		if argType != paramType {
			typeErrors = append(typeErrors, fmt.Sprintf("argument %d should be of type %v but found type %v", i, paramType, argType))
		}
	}

	if len(typeErrors) > 0 {
		return errors.New(strings.Join(typeErrors, "; "))
	}
	return nil
}
