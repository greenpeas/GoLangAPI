package shipping

import (
	"seal/internal/domain/user"
	"time"
)

type Shipping struct {
	Db
	Author    user.Author `json:"author"`
	Transport struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
		Type  struct {
			Id    int    `json:"id"`
			Title string `json:"title"`
		} `json:"type"`
		RegistrationNumber string `json:"registration_number"`
	} `json:"transport"`
	Route struct {
		Id         int      `json:"id"`
		Title      string   `json:"title"`
		Points     []string `json:"points"`
		Length     int      `json:"length"`
		TravelTime int      `json:"travel_time"`
	} `json:"route"`
	Files []File `json:"files"`
	Modem *struct {
		Id     int `json:"id"`
		Imei   int `json:"imei"`
		Serial int `json:"serial"`
		Last   *struct {
			Status       uint16 `json:"status"`
			BatteryLevel int16  `json:"battery_level"`
			Rssi         int16  `json:"rssi"`
		} `json:"last"`
	} `json:"modem"`
	Seals []struct {
		Id     int    `json:"id"`
		Serial uint64 `json:"serial"`
		Last   struct {
			Status       int32 `json:"status"`
			BatteryLevel int16 `json:"battery_level"`
			Rssi         int16 `json:"rssi"`
		} `json:"last"`
	} `json:"seals"`
	EstimatedArrivalTime time.Time `json:"estimated_arrival_time" db:"estimated_arrival_time"`
}

type ShippingForList struct {
	Id           int    `json:"id"`
	CustomNumber string `json:"custom_number" db:"custom_number"`
	CreateDate   string `json:"create_date" db:"create_date"`
	Number       int    `json:"number"`
	Status       int    `json:"status"`
	Transport    struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
		Type  struct {
			Id    int    `json:"id"`
			Title string `json:"title"`
		} `json:"type"`
		RegistrationNumber string `json:"registration_number"`
	} `json:"transport"`
	Route struct {
		Id         int      `json:"id"`
		Title      string   `json:"title"`
		Points     []string `json:"points"`
		Length     int      `json:"length"`
		TravelTime int      `json:"travel_time"`
	} `json:"route"`
	Files                []File    `json:"files"`
	EstimatedArrivalTime time.Time `json:"estimated_arrival_time" db:"estimated_arrival_time"`
}

type CreateRequest struct {
	CustomNumber string `json:"custom_number" db:"custom_number" validate:"required,max=8,min=8"`
	CreateDate   string `json:"create_date" db:"create_date" validate:"required,max=6,min=6"`
	Number       int    `json:"number" validate:"required,max=2147483647,min=1"`
	Transport    int    `json:"transport" validate:"required,max=2147483647,min=1"`
	Route        int    `json:"route" validate:"required,max=2147483647,min=1"`
}

type UpdateRequest struct {
	CustomNumber string `json:"custom_number,omitempty" db:"custom_number" validate:"max=8,min=0"`
	CreateDate   string `json:"create_date,omitempty" db:"create_date" validate:"max=6,min=0"`
	Number       int    `json:"number,omitempty" validate:"max=2147483647,min=0"`
	Transport    int    `json:"transport,omitempty" validate:"max=2147483647,min=0"`
	Route        int    `json:"route,omitempty" validate:"max=2147483647,min=0"`
}

type trackResponseSeal struct {
	Id     int    `json:"id"`
	Serial uint64 `json:"serial"`
}

type trackResponseSealData struct {
	DevTime      time.Time         `json:"dev_time"`
	Seal         trackResponseSeal `json:"seal"`
	Status       uint32            `json:"status"`
	Errors       uint8             `json:"errors"`
	BatteryLevel uint8             `json:"battery_level"`
}

type trackResponseTelemetry struct {
	DevTime     time.Time               `json:"dev_time"`
	Status      int32                   `json:"status"`
	Latitude    any                     `json:"latitude"`
	Longitude   any                     `json:"longitude"`
	Altitude    int32                   `json:"altitude"`
	SealsData   []trackResponseSealData `json:"sealsData"`
	ErrorsFlags int32                   `json:"errors_flags"`
}

type trackResponseCoordinate struct {
	DevTime   time.Time `json:"dev_time"`
	Latitude  float32   `json:"latitude"`
	Longitude float32   `json:"longitude"`
	//Altitude  int32     `json:"altitude"`
	MinDistanceToRoute *int `json:"min_distance_to_route"`
}

type TrackQueryParams struct {
	Id        int       `validate:"required"`
	From      time.Time `form:"from"`
	To        time.Time `form:"to"`
	Limit     int       `form:"limit" validate:"max=500000"`
	OrderDesc bool      `form:"order_desc"`
}

type AddFileRequest struct {
	Comment string `form:"comment,omitempty"`
	SealId  int    `form:"seal_id,omitempty"`
}

type UpdateFileRequest struct {
	Comment *string `json:"comment,omitempty"`
	Title   *string `json:"title,omitempty"`
	SealId  *int    `json:"seal_id,omitempty"`
}
