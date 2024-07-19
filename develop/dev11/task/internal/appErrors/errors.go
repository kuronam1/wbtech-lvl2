package appErrors

import (
	"errors"
	"fmt"
)

var (
	FreeDay  = errors.New("free day")
	NotSaved = errors.New("event is not saved")
	Existis  = errors.New("event already exists")
)

func WrapErr(functionName, errDescription string, err error) error {
	return fmt.Errorf("[%s] %s: %v", functionName, errDescription, err)
}
