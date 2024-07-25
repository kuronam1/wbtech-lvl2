package appErrors

import (
	"errors"
	"fmt"
)

var (
	FreeDay  = errors.New("free day")
	NotSaved = errors.New("event is not saved")
	Exists   = errors.New("event already exists")
)

type LogError struct {
	Function string
	err      error
}

func WrapErr(functionName, errDescription string, err error) error {
	if err == nil {
		return fmt.Errorf("[%s] %s", functionName, errDescription)
	}

	return fmt.Errorf("[%s] %s: %v", functionName, errDescription, err)
}
