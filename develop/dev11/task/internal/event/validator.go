package event

import (
	"errors"
)

var (
	ZeroValueDescriptionName = errors.New("zero value description name")
	NameLengthOverflow       = errors.New("name length overflow")
	ZeroValueDescriptionText = errors.New("zero value description text")
	TextLengthOverflow       = errors.New("text length overflow")
	ZeroDate                 = errors.New("zero value date")
)

func (e *Event) Validate() error {
	if e.Date.IsZero() {
		return ZeroDate
	}

	if e.Description.Name == "" {
		return ZeroValueDescriptionName
	}

	if e.nameLen() > 20 {
		return NameLengthOverflow
	}

	if e.Description.Text == "" {
		return ZeroValueDescriptionText
	}

	if e.textLen() > 512 {
		return TextLengthOverflow
	}

	return nil
}

func (e *Event) nameLen() int {
	return len([]rune(e.Description.Name))
}

func (e *Event) textLen() int {
	return len([]rune(e.Description.Text))
}
