package client

import (
	"encoding/csv"
	"events/task/internal/appErrors"
	"events/task/internal/event"
	"io"
	"log/slog"
	"os"
	"time"
)

type client struct {
	data   *os.File
	logger *slog.Logger
}

func (c *client) GetEventsFromFile() ([]event.Event, error) {
	const function = "client.GetEventsFromFile"
	events := make([]event.Event, 0, 10)

	reader := csv.NewReader(c.data)
	c.logger.Debug("starting reading loop")
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				c.logger.Debug("reached end of the file")
				break
			}
			return nil, appErrors.WrapErr(function, "error while reading csv line:", err)
		}

		date, err := time.Parse(time.DateOnly, record[0])
		if err != nil {
			return nil, appErrors.WrapErr(function, "error while parsing event date", err)
		}

		events = append(events, event.Event{
			Date: date,
			Description: event.Description{
				Name:            record[1],
				DescriptionText: record[2],
			},
		})
	}

	return events, nil
}

func (c *client) SaveEventsToFile(events []event.Event) (err error) {
	const functionName = "client.SaveEventsToFile"

	csvWriter := csv.NewWriter(c.data)
	c.logger.Info("starting writing loop")

	for _, element := range events {
		eventRecord := eventToSlice(element)

		if err = csvWriter.Write(eventRecord); err != nil {
			return appErrors.WrapErr(functionName, "error while saving event", err)
		}
	}
	csvWriter.Flush()

	if err = csvWriter.Error(); err != nil {
		return appErrors.WrapErr(functionName, "error: ", err)
	}

	return
}

func eventToSlice(event event.Event) []string {
	return []string{event.Date.String(), event.Description.Name, event.Description.DescriptionText}
}
