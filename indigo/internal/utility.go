package internal

import (
	"fmt"
	"reflect"

	"github.com/joshua-zingale/indigo/indigo/functools"
	"github.com/joshua-zingale/indigo/indigo/interfaces"
)

func ListToSlice(l interfaces.List) []any {
	switch slice := l.(type) {
	case List:
		return slice
	}
	return listToSlice(l)
}

func listToSlice(l interfaces.List) []any {
	var slice []any

	for !l.Empty() {
		slice = append(slice, l.Car())
		l = l.Cdr().(interfaces.List)
	}
	return slice
}

type VerifiedList struct {
	interfaces.Cons
}

func (vl *VerifiedList) IsList() {}

func ValidateList(l interfaces.Cons) (interfaces.List, error) {
	for !l.Empty() {
		next := l.Cdr()

		if cons, ok := next.(interfaces.Cons); ok {
			l = cons
		} else if next == nil {
			l = nil
		} else {
			return nil, fmt.Errorf("Cons is not a List")
		}
	}

	return &VerifiedList{Cons: l}, nil
}

func ValidateType(desiredType reflect.Type, value any) (any, error) {
	valueVal := reflect.ValueOf(value)
	valueType := valueVal.Type()

	if valueType.AssignableTo(desiredType) {
		return value, nil
	} else if valueType.ConvertibleTo(desiredType) {
		convertedValue := valueVal.Convert(desiredType)
		return convertedValue.Interface(), nil
	} else {
		return nil, fmt.Errorf("cannot use %v of type %v as %v", value, valueType, desiredType)
	}
}

func ValidateFunctionArgs(parameterTypes []reflect.Type, args []any) ([]any, error) {
	if len(args) != len(parameterTypes) {
		return nil, fmt.Errorf("expected %d arguments but found %d", len(parameterTypes), len(args))
	}
	return functools.MapWithError(func(pair functools.Pair[reflect.Type, any]) (any, error) {
		return ValidateType(pair.First, pair.Second)
	}, functools.Zip(parameterTypes, args))
}
