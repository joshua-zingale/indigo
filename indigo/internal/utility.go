package internal

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

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

func ValidateFunctionArgs(parameterTypes []reflect.Type, args []any) ([]any, error) {
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
