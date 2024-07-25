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

func (c *client) GetEventsFromFile() (map[time.Time][]event.Description, error) {
	const function = "client.GetEventsFromFile"
	events := make(map[time.Time][]event.Description)

	reader := csv.NewReader(c.data)
	c.logger.Debug("starting reading loop")

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				c.logger.Debug("reached end of the file")
				break
			} else {
				return nil, appErrors.WrapErr(function, "error while reading csv line:", err)
			}
		}

		date, err := time.Parse(time.DateOnly, record[0])
		if err != nil {
			return nil, appErrors.WrapErr(function, "error while parsing event date", err)
		}

		eventTime, err := time.Parse(time.TimeOnly, record[1])
		if err != nil {
			return nil, appErrors.WrapErr(function, "error while parsing event time", err)
		}

		events[date] = append(events[date], event.Description{
			Time: eventTime,
			Name: record[2],
			Text: record[3],
		})
	}

	return events, nil
}

func (c *client) SaveEventsToFile(events map[time.Time][]event.Description) (err error) {
	const functionName = "client.SaveEventsToFile"

	csvWriter := csv.NewWriter(c.data)
	c.logger.Info("starting writing loop")

	for key, element := range events {
		for _, eventDescription := range element {
			eventRecord := eventToSlice(key, eventDescription)

			if err = csvWriter.Write(eventRecord); err != nil {
				return appErrors.WrapErr(functionName, "error while saving event", err)
			}

			csvWriter.Flush()
		}
	}

	if err = csvWriter.Error(); err != nil {
		return appErrors.WrapErr(functionName, "error: ", err)
	}

	return nil
}

func eventToSlice(date time.Time, description event.Description) []string {
	return []string{date.String(), description.Time.String(), description.Name, description.Text}
}
