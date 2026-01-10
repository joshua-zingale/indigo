package evaluation

import (
	"reflect"
)

func invoke(funcToCall any, args []any) (any, error) {
	inputs := make([]reflect.Value, len(args))
	for i, arg := range args {
		inputs[i] = reflect.ValueOf(arg)
	}

	values := reflect.ValueOf(funcToCall).Call(inputs)
	if len(values) != 2 {
		panic("all invoked functions must return two values, the primary value and a (nil-able) error")
	}
	var errorVal error = nil
	if secondReturnValue, ok := values[1].Interface().(error); ok {
		errorVal = secondReturnValue
	} else if values[1].Interface() != nil {
		panic("all invoked functions must return two values, the primary value and a (nil-able) error")
	}
	return values[0].Interface(), errorVal
}

func getFuncParameterTypes(function any) []reflect.Type {
	functionType := reflect.TypeOf(function)
	numParameters := functionType.NumIn()
	parameterTypes := make([]reflect.Type, numParameters)
	for i := range parameterTypes {
		parameterTypes[i] = functionType.In(i)
	}

	return parameterTypes
}
