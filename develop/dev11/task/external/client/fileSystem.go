package client

import (
	"events/task/internal/appErrors"
	"events/task/internal/event"
	"log/slog"
	"os"
	"time"
)

type File interface {
	GetEventsFromFile() (map[time.Time][]event.Description, error)
	SaveEventsToFile(events map[time.Time][]event.Description) (err error)
}

func GetFSClient(logger *slog.Logger, pathToData string) (File, error) {
	const function = "fs.NewClient"
	var (
		f   *os.File
		err error
	)

	f, err = os.Open(pathToData)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(pathToData)
			if err != nil {
				return nil, appErrors.WrapErr(function, "error file creating new file:", err)
			}
		}
		return nil, appErrors.WrapErr(function, "error while opening file: %s", err)
	}

	return &client{data: f, logger: logger}, nil
}
