package repository

import (
	"events/task/external/client"
	"events/task/internal/event"
	"log/slog"
	"time"
)

type Repository interface {
	CreateEvent(event event.Event) error
	UpdateEvent(newEvent event.Event) error
	DeleteEvent(event event.Event)
	DayEvents(date time.Time) ([]event.Event, error)
}

func NewRepository(client client.File, logger *slog.Logger) (Repository, error) {
	events, err := client.GetEventsFromFile()
	if err != nil {
		return nil, err
	}

	return &repository{
		cacheData: events,
		logger:    logger,
	}, nil
}
