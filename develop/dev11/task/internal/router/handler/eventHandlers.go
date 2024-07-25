package handler

import (
	"encoding/json"
	"errors"
	"events/task/internal/appErrors"
	"events/task/internal/event"
	"events/task/internal/event/repository"
	"events/task/internal/parser"
	"events/task/internal/router/middlewear"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type errResponse struct {
	Error string `json:"error"`
}

type OkResponse struct {
	Result any `json:"result"`
}

var (
	unknownMethod = errors.New("unknown method")
	unknownPath   = errors.New("not found")
)

type Handler struct {
	logger *slog.Logger
	repo   repository.Repository
}

func New(logger *slog.Logger, repository repository.Repository) *Handler {
	return &Handler{
		logger: logger,
		repo:   repository,
	}
}

func (h *Handler) Register(router *http.ServeMux) {
	router.Handle("/create_event", middlewear.LoggingRequest(h.logger, h.Creation))
	router.Handle("/update_event", middlewear.LoggingRequest(h.logger, h.Update))
	router.Handle("/delete_event", middlewear.LoggingRequest(h.logger, h.Delete))
	router.Handle("/events_for_day", middlewear.LoggingRequest(h.logger, h.FindDayEvents))
	router.Handle("/events_for_week", middlewear.LoggingRequest(h.logger, h.FindWeekEvents))
	router.Handle("/events_for_month", middlewear.LoggingRequest(h.logger, h.FindMonthEvents))
}

func (h *Handler) Creation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONResponse(w, errResponse{Error: unknownMethod.Error()}, http.StatusBadRequest)
		return
	}

	if r.URL.Path != "/create_event" {
		h.writeJSONResponse(w, errResponse{Error: unknownPath.Error()}, http.StatusNotFound)
		return
	}

	dayEvent, err := parser.ParseRequest(r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		h.writeJSONResponse(w, errResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	defer func() {
		if err = r.Body.Close(); err != nil {
			h.logger.Error("error while closing request body", "error", err.Error())
		}
	}()

	if err = h.repo.CreateEvent(dayEvent); err != nil {
		h.logger.Error(err.Error())
		h.writeJSONResponse(w, errResponse{Error: err.Error()}, http.StatusInternalServerError)
	}

	h.logger.Info("new event created", "event", dayEvent)
	h.writeJSONResponse(w, OkResponse{Result: "new event is created"}, http.StatusOK)
}

func (h *Handler) FindDayEvents(w http.ResponseWriter, r *http.Request) {
	const function = "handler.findHandler"

	if r.Method != http.MethodGet {
		h.writeJSONResponse(w, errResponse{Error: unknownMethod.Error()}, http.StatusBadRequest)
		return
	}

	if r.URL.Path != "/events_for_day" {
		h.writeJSONResponse(w, errResponse{Error: unknownPath.Error()}, http.StatusNotFound)
		return
	}

	date, err := h.getDateFromURL(r.URL.Query())
	if err != nil {
		h.logger.Error("error while getting date from url", "error", err.Error())
		h.writeJSONResponse(w, errResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	events, err := h.repo.DayEvents(date)
	if err != nil {
		h.logger.Error("error while getting day events", "error", err.Error())
		h.writeJSONResponse(w, errResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	h.logger.Info("found events", "events", events)
	h.writeJSONResponse(w, OkResponse{Result: events}, http.StatusOK)
}

func (h *Handler) FindWeekEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONResponse(w, errResponse{Error: unknownMethod.Error()}, http.StatusBadRequest)
		return
	}

	if r.URL.Path != "/events_for_week" {
		h.writeJSONResponse(w, errResponse{Error: unknownPath.Error()}, http.StatusNotFound)
		return
	}

	date, err := h.getDateFromURL(r.URL.Query())
	if err != nil {
		h.logger.Error("error while getting date from url", "error", err.Error())
		h.writeJSONResponse(w, errResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	defer func() {
		if err = r.Body.Close(); err != nil {
			h.logger.Error("error while closing request body", "error", err.Error())
		}
	}()

	offset := int(time.Sunday - date.Weekday())
	date = date.AddDate(0, 0, offset)
	start := date.Weekday().String()

	events := make([]event.Event, 0)
	for date.Weekday() != 6 {
		dayEvents, err := h.repo.DayEvents(date)
		if err != nil {
			continue
		}

		events = append(events, dayEvents...)
		date = date.Add(24 * time.Hour)
	}

	if len(events) == 0 {
		h.logger.Info(fmt.Sprintf("no events for the week from %s to %s", start, date.Weekday().String()))
		h.writeJSONResponse(w, OkResponse{Result: "Free week!"}, 200)
	} else {
		h.logger.Info(fmt.Sprintf("got all events from %s to %s ", start, date.Weekday().String()),
			"events", events)
		h.writeJSONResponse(w, OkResponse{Result: events}, 200)
	}
}

func (h *Handler) FindMonthEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONResponse(w, errResponse{Error: unknownMethod.Error()}, http.StatusNotFound)
		return
	}

	if r.URL.Path != "/events_for_month" {
		h.writeJSONResponse(w, errResponse{Error: unknownPath.Error()}, http.StatusNotFound)
		return
	}

	date, err := h.getDateFromURL(r.URL.Query())
	if err != nil {
		h.logger.Error("error while getting date from url", "error", err.Error())
		h.writeJSONResponse(w, errResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	defer func() {
		if err = r.Body.Close(); err != nil {
			h.logger.Error("error while closing request body", "error", err.Error())
		}
	}()

	searchMonth := date.Month()
	currentDay := time.Date(date.Year(), searchMonth, 1, 0, 0, 0, 0, date.Location())

	var (
		events    []event.Event
		dayEvents []event.Event
	)
	for currentDay.Month() == searchMonth {
		dayEvents, err = h.repo.DayEvents(currentDay)
		if err != nil {
			if errors.Is(err, appErrors.FreeDay) {
				continue
			} else {
				h.logger.Error("error while getting day events", "error", err.Error())
				h.writeJSONResponse(w, errResponse{Error: err.Error()}, http.StatusInternalServerError)
				return
			}
		}

		events = append(events, dayEvents...)
		currentDay = currentDay.Add(24 * time.Hour)
	}

	if len(events) == 0 {
		h.logger.Info(fmt.Sprintf("no events for the month %s", searchMonth.String()))
		h.writeJSONResponse(w, OkResponse{Result: "Free month!"}, 200)
	} else {
		h.logger.Info(fmt.Sprintf("got all events for the month %s", searchMonth.String()),
			"events", events)
		h.writeJSONResponse(w, OkResponse{Result: events}, 200)
	}
}

func (h *Handler) getDateFromURL(parameters url.Values) (time.Time, error) {
	const function = "handler.getDateFromURL"

	if !parameters.Has("date") {
		h.logger.Warn("no date in parameters")
		return time.Now(), nil
	}

	date, err := time.Parse(time.DateOnly, parameters.Get("date"))
	if err != nil {
		h.logger.Error("error while parsing date", "function", function, "error", err.Error())
		return time.Time{}, appErrors.WrapErr(function, "error while parsing date", err)
	}

	return date, nil
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONResponse(w, errResponse{Error: unknownMethod.Error()}, http.StatusBadRequest)
		return
	}

	if r.URL.Path != "/delete_event" {
		h.writeJSONResponse(w, errResponse{Error: unknownPath.Error()}, http.StatusNotFound)
		return
	}

	deleteEvent, err := parser.ParseRequest(r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		h.writeJSONResponse(w, errResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	defer func() {
		if err = r.Body.Close(); err != nil {
			h.logger.Error("error while closing request body", "error", err.Error())
		}
	}()

	h.repo.DeleteEvent(deleteEvent)

	h.logger.Info("event has been deleted", "deleted event", deleteEvent)
	h.writeJSONResponse(w, OkResponse{
		Result: fmt.Sprintf("event %v has been deleted", deleteEvent),
	}, http.StatusOK)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONResponse(w, errResponse{Error: unknownMethod.Error()}, http.StatusBadRequest)
		return
	}

	if r.URL.Path != "/update_event" {
		h.writeJSONResponse(w, errResponse{Error: unknownPath.Error()}, http.StatusNotFound)
		return
	}

	updateEvent, err := parser.ParseRequest(r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		h.writeJSONResponse(w, errResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	defer func() {
		if err = r.Body.Close(); err != nil {
			h.logger.Error("error while closing request body", "error", err.Error())
		}
	}()

	if err = h.repo.UpdateEvent(updateEvent); err != nil {
		h.logger.Error("error while updating event", "error", err.Error())
		h.writeJSONResponse(w, errResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	h.logger.Info("event has been updated. New event", "event", updateEvent)
	h.writeJSONResponse(w, OkResponse{Result: "event has been updated"}, http.StatusOK)
}

func (h *Handler) writeJSONResponse(w http.ResponseWriter, data any, code int) {
	const functionName = "handler.writeJSONResponse"

	response, err := json.Marshal(data)
	if err != nil {
		h.logger.Error(appErrors.WrapErr(functionName, "error while marshaling response to json", err).Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err = w.Write(response); err != nil {
		h.logger.Error("error while writing error response", "error", err.Error())
	}
}
