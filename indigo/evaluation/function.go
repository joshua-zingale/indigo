package evaluation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

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
			validatedArgs, err := validateFunctionArgs(parameterTypes, evaluatedArgs)
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

func validateFunctionArgs(parameterTypes []reflect.Type, args []any) ([]any, error) {
	if len(args) != len(parameterTypes) {
		return nil, fmt.Errorf("expected %d arguments but found %d", len(parameterTypes), len(args))
	}

	var typeErrors []string
	convertedArgs := make([]any, len(args))
	for i := range args {
		argVal := reflect.ValueOf(args[i])
		argType := argVal.Type()
		paramType := parameterTypes[i]

		if argType.AssignableTo(paramType) {
			convertedArgs[i] = args[i]
		} else if argType.ConvertibleTo(paramType) {
			convertedValue := argVal.Convert(paramType)
			convertedArgs[i] = convertedValue.Interface()
		} else {
			typeErrors = append(typeErrors, fmt.Sprintf("argument %d should be of type %v but found type %v", i, paramType, argType))
		}
	}

	if len(typeErrors) > 0 {
		return nil, errors.New(strings.Join(typeErrors, "; "))
	}
	return convertedArgs, nil
}
