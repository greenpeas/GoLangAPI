package command

import (
	"time"
)

type Command struct {
	Id           int32      `json:"id"`
	Serial       string     `json:"serial"`
	Name         string     `json:"name"`
	Params       any        `json:"params"`
	Author       string     `json:"author"`
	Dateon       time.Time  `json:"dateon"`
	TryNumber    int32      `json:"try_number"`
	TryDate      *time.Time `json:"try_date"`
	ResponseDate *time.Time `json:"response_date"`
	AbortDate    *time.Time `json:"abort_date"`
	Response     *string    `json:"response"`
	RawRequest   *string    `json:"raw_request"`
	RawResponse  *string    `json:"raw_response"`
}
