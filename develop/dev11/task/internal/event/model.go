package event

import "time"

type Event struct {
	Date        time.Time   `json:"date"`
	Description Description `json:"description"`
}

type Description struct {
	Time time.Time `json:"time"`
	Name string    `json:"name"`
	Text string    `json:"text"`
}
