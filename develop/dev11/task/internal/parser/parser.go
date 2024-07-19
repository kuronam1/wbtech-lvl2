package parser

import (
	"encoding/json"
	"events/task/internal/appErrors"
	"events/task/internal/event"
	"events/task/internal/parser/validator"
	"io"
)

func ParseRequest(body io.ReadCloser) (event.Event, error) {
	const functionName = "parser.ParseRequest"
	result := event.Event{}

	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return result,
			appErrors.WrapErr(functionName, "error while parsing request:", err)
	}

	if err := validator.ValidateStruct(result); err != nil {
		return result, appErrors.WrapErr(functionName, "error while validating request:", err)
	}

	return result, nil
}
