package repository

import (
	"events/task/internal/appErrors"
	"events/task/internal/event"
	"log/slog"
	"sort"
	"time"
)

type repository struct {
	cacheData map[time.Time][]event.Description
	logger    *slog.Logger
}

func (r *repository) CreateEvent(event event.Event) error {
	if value, inMap := r.cacheData[event.Date]; inMap {
		if _, ok := contains(value, event.Description); ok {

			return appErrors.Exists
		}
	}

	r.cacheData[event.Date] = append(r.cacheData[event.Date], event.Description)

	sort.Slice(r.cacheData[event.Date], func(i, j int) bool {
		return r.cacheData[event.Date][i].Time.Before(r.cacheData[event.Date][j].Time)
	})

	return nil
}

func (r *repository) UpdateEvent(newEvent event.Event) error {
	if currentEvent, inMap := r.cacheData[newEvent.Date]; inMap {
		for idx, description := range currentEvent {
			if description.Time.Equal(newEvent.Description.Time) {
				r.cacheData[newEvent.Date][idx] = newEvent.Description
				return nil
			}
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

func contains(eventsDesc []event.Description, event2 event.Description) (int, bool) {
	for idx, eventDescription := range eventsDesc {
		if eventDescription == event2 {
			return idx, true
		}
	}

	return -1, false
}
