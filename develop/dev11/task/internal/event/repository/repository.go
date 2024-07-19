package repository

import (
	"events/task/external/client"
	"events/task/internal/event"
	"log/slog"
	"time"
)

type Repository interface {
	CreateEvent(event event.Event) error
	UpdateEvent(oldEvent, newEvent event.Event) error
	DeleteEvent(event event.Event)
	DayEvents(date time.Time) ([]event.Event, error)
}

func NewRepository(client client.File, logger *slog.Logger) (Repository, error) {
	eventsArray, err := client.GetEventsFromFile()
	if err != nil {
		return nil, err
	}

	events := make(map[time.Time][]event.Description)
	for _, element := range eventsArray {
		events[element.Date] = append(events[element.Date], element.Description)
	}

	return &repository{
		cacheData: events,
		logger:    logger,
	}, nil
}
