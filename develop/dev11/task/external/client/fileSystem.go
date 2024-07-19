package client

import (
	"events/task/internal/appErrors"
	"events/task/internal/event"
	"log/slog"
	"os"
)

type File interface {
	GetEventsFromFile() ([]event.Event, error)
	SaveEventsToFile(events []event.Event) (err error)
}

func GetFSClient(logger *slog.Logger) (File, error) {
	const function = "fs.NewClient"
	var (
		f   *os.File
		err error
	)

	f, err = os.Open("event.csv")
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create("event.csv")
			if err != nil {
				return nil, appErrors.WrapErr(function, "error file creating new file:", err)
			}
		}
		return nil, appErrors.WrapErr(function, "error while opening file: %s", err)
	}

	return &client{data: f, logger: logger}, nil
}
