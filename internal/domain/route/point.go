package route

import (
	"time"
)

type PointDb struct {
	Route     int       `json:"route" validate:"required,max=32767,min=1"`
	Number    int       `json:"number" validate:"required,max=32767,min=0"`
	Latitude  int       `json:"latitude" db:"latitude" validate:"required,max=90,min=-90"`
	Longitude int       `json:"longitude" db:"longitude" validate:"required,max=180,min=-180"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Point = Db

type Points struct {
	Route  int `json:"route" validate:"required,max=32767,min=1"`
	Coords []struct {
		Latitude  float32 `json:"latitude" db:"latitude" validate:"required,max=90,min=-90"`
		Longitude float32 `json:"longitude" db:"longitude" validate:"required,max=180,min=-180"`
	} `json:"coords"`
	CustomDeparture   int `json:"custom_departure"`
	CustomDestination int `json:"custom_destination"`
}
