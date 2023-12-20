package modemData

import "time"

type ModemData struct {
	Db
	CoordinatesLbs *struct {
		DevTime         time.Time `db:"dev_time" json:"-"`
		Modem           int       `db:"modem" json:"-"`
		PositioningTime time.Time `db:"positioning_time" json:"positioning_time"`
		Latitude        float32   `db:"latitude" json:"latitude"`
		Longitude       float32   `db:"longitude" json:"longitude"`
		Precision       int       `db:"precision" json:"precision"`
	} `json:"coordinate_lbs" db:"coordinate_lbs"`
}

type ListParams struct {
	ModemId   int
	TimeFrom  time.Time
	TimeTo    time.Time
	Limit     int
	OrderDesc bool
}
