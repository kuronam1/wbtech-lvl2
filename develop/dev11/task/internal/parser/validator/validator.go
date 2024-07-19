package validator

import (
	"errors"
	"events/task/internal/event"
)

var (
	ZeroValue      = errors.New("zero value description")
	ZeroDate       = errors.New("zero value date")
	LengthOverflow = errors.New("length overflow")
)

func ValidateStruct(event event.Event) error {
	if event.Date.IsZero() {
		return ZeroDate
	} else if event.Description.Name == "" || event.Description.DescriptionText == "" {
		return ZeroValue
	}

	if len([]rune(event.Description.Name)) > 20 || len([]rune(event.Description.DescriptionText)) > 512 {
		return LengthOverflow
	}

	return nil
}
