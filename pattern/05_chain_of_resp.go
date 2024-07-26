package main

import "fmt"

type Handler interface {
	SetNext(handler Handler)
	HandleRequest(level int, message string)
}

type BaseHandler struct {
	next Handler
}

func (b *BaseHandler) SetNext(handler Handler) {
	b.next = handler
}

func (b *BaseHandler) PassToNext(level int, message string) {
	if b.next != nil {
		b.next.HandleRequest(level, message)
	}
}

const (
	INFO    = 1
	WARNING = 2
	ERROR   = 3
)

type InfoHandler struct {
	BaseHandler
}

func (h *InfoHandler) HandleRequest(level int, message string) {
	if level == INFO {
		fmt.Println("Info: ", message)
	} else {
		h.PassToNext(level, message)
	}
}

type WarningHandler struct {
	BaseHandler
}

func (h *WarningHandler) HandleRequest(level int, message string) {
	if level == WARNING {
		fmt.Println("Warning: ", message)
	} else {
		h.PassToNext(level, message)
	}
}

type ErrorHandler struct {
	BaseHandler
}

func (h *ErrorHandler) HandleRequest(level int, message string) {
	if level == ERROR {
		fmt.Println("Error: ", message)
	} else {
		h.PassToNext(level, message)
	}
}

func main() {
	infoHandler := &InfoHandler{}
	warningHandler := &WarningHandler{}
	errorHandler := &ErrorHandler{}

	infoHandler.SetNext(warningHandler)
	warningHandler.SetNext(errorHandler)

	fmt.Println("Processing requests:")
	infoHandler.HandleRequest(INFO, "This is an info message.")
	infoHandler.HandleRequest(WARNING, "This is a warning message.")
	infoHandler.HandleRequest(ERROR, "This is an error message.")
}
