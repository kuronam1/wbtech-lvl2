package event

import "time"

type Event struct {
	Date        time.Time
	Description Description
}

type Description struct {
	Name            string
	DescriptionText string
}
