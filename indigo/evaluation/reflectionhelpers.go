package evaluation

import "reflect"

func invoke(funcToCall any, args []any) (any, error) {
	inputs := make([]reflect.Value, len(args))
	for i, arg := range args {
		inputs[i] = reflect.ValueOf(arg)
	}

	values := reflect.ValueOf(funcToCall).Call(inputs)
	if len(values) != 2 {
		panic("all invoked functions must return two values, the primary value and a (nil-able) error")
	}
	return values[0].Interface(), values[1].Interface().(error)
}

func getFuncParameters(function any) []reflect.Type {
	functionType := reflect.TypeOf(function)
	numParameters := functionType.NumIn()
	parameterTypes := make([]reflect.Type, numParameters)
	for i := range parameterTypes {
		parameterTypes = append(parameterTypes, functionType.In(i))
	}

	return parameterTypes
}

func valuesToTypes(values []any) []reflect.Type {
	types := make([]reflect.Type, len(values))
	for _, v := range values {
		types = append(types, reflect.TypeOf(v))
	}
	return types
}

func isGoFunc(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Func
}
