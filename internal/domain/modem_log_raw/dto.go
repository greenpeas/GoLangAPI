package modemLogRaw

import "time"

type ModemLogRaw struct {
	RegTime        time.Time `json:"reg_time" db:"reg_time"`
	Imei           uint64
	Src            int
	Hex            string
	Payload        any
	CmdName        string `json:"cmd_name" db:"cmd_name"`
	CmdDescription string `json:"cmd_description" db:"cmd_description"`
	RemotePort     int    `json:"remote_port" db:"remote_port"`
}

type ListParams struct {
	Imei      string
	From      time.Time
	To        time.Time
	Limit     int
	OrderDesc bool
}

type Telemetry struct {
	CurrentTime     int64     `db:"current_time"`
	RegTime         time.Time `db:"reg_time"`
	Status          int32     `db:"status"`
	ErrorsFlags     int32     `db:"errors_flags"`
	PositioningTime int64     `db:"positioning_time"`
	Latitude        any       `db:"latitude"`
	Longitude       any       `db:"longitude"`
	Altitude        int32     `db:"altitude"`
	SatellitesCount int16     `db:"satellites_count"`
	Speed           int32     `db:"speed"`
	StatusGpsModule int16     `db:"status_gps_module"`
	Rssi            int16     `db:"rssi"`
	BatteryLevel    int16     `db:"battery_level"`
	SignalGps       int32     `db:"signal_gps"`
	SignalGlonass   int32     `db:"signal_glonass"`
}
