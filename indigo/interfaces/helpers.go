package interfaces

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func ValidateFunctionArgs(function IndigoFunction, args []any) error {
	if len(args) != len(function.ParameterTypes()) {
		return fmt.Errorf("expected %d arguments but found %d", len(function.ParameterTypes()), len(args))
	}

	var typeErrors []string
	for i := range args {
		argType := reflect.TypeOf(args[i])
		paramType := function.ParameterTypes()[i]

		if argType != paramType {
			typeErrors = append(typeErrors, fmt.Sprintf("argument %d should be of type %v but found type %v", i, paramType, argType))
		}
	}

	if len(typeErrors) > 0 {
		return errors.New(strings.Join(typeErrors, "; "))
	}
	return nil
}
