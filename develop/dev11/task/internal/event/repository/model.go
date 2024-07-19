package repository

import (
	"events/task/internal/appErrors"
	"events/task/internal/event"
	"log/slog"
	"time"
)

type repository struct {
	cacheData map[time.Time][]event.Description
	logger    *slog.Logger
}

func (r *repository) CreateEvent(event event.Event) error {
	if value, inMap := r.cacheData[event.Date]; inMap {
		if _, ok := contains(value, event.Description); ok {

			return appErrors.Existis
		}
	}

	r.cacheData[event.Date] = append(r.cacheData[event.Date], event.Description)

	return nil
}

func contains(eventsDesc []event.Description, event2 event.Description) (int, bool) {
	for idx, eventDescription := range eventsDesc {
		if eventDescription.Name == event2.Name ||
			eventDescription.DescriptionText == event2.DescriptionText {

			return idx, true
		}
	}

	return -1, false
}

func (r *repository) UpdateEvent(oldEvent, newEvent event.Event) error {
	if value, inMap := r.cacheData[oldEvent.Date]; inMap {
		if idx, ok := contains(value, oldEvent.Description); ok {
			r.cacheData[oldEvent.Date][idx] = newEvent.Description
			return nil
		}
	}

	return appErrors.NotSaved
}

func (r *repository) DeleteEvent(event event.Event) {
	if value, inMap := r.cacheData[event.Date]; inMap {
		if idx, ok := contains(value, event.Description); ok {
			r.cacheData[event.Date] = append(r.cacheData[event.Date][:idx], r.cacheData[event.Date][idx+1:]...)
		}
	}
}

func (r *repository) DayEvents(date time.Time) ([]event.Event, error) {
	if events, inMap := r.cacheData[date]; inMap {
		result := make([]event.Event, 0, len(events))

		for _, element := range events {
			result = append(result, event.Event{
				Date:        date,
				Description: element,
			})
		}

		return result, nil
	} else {
		return nil, appErrors.FreeDay
	}
}
