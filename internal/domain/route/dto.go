package route

import "time"

type Route struct {
	Id         int          `json:"id"`
	Title      string       `json:"title"`
	Points     []string     `json:"points"`
	Length     int          `json:"length"`
	CreatedAt  *time.Time   `json:"created_at" db:"created_at"`
	Coords     [][2]float32 `json:"coords" db:"coords"`
	TravelTime int          `json:"travel_time" db:"travel_time"`
}

func (r Route) getId() int {
	return r.Id
}

func (r Route) getTitle() string {
	return r.Title
}

type CreateRequest struct {
	Title      string       `json:"title" validate:"required,max=200,min=1"`
	Points     []string     `json:"points" validate:"required"`
	Length     int          `json:"length" validate:"max=32767,min=0"`
	Coords     [][2]float32 `json:"coords" db:"coords"`
	TravelTime int          `json:"travel_time"`
}

type UpdateRequest struct {
	Title      *string  `json:"title,omitempty"`
	Points     []string `json:"points,omitempty"`
	Length     *int     `json:"length,omitempty"`
	TravelTime *int     `json:"travel_time,omitempty"`
}
