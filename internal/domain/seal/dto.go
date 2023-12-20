package seal

import (
	"time"
)

type Data struct {
	DevTime              time.Time `json:"dev_time"`
	Seal                 int       `json:"seal"`
	ModemTime            time.Time `json:"modem_time"`
	Modem                int       `json:"modem"`
	Status               int64     `json:"status"`
	Errors               int16     `json:"errors"`
	SensitivityRange     int16     `json:"sensitivity_range"`
	BatteryLevel         int16     `json:"battery_level"`
	Rssi                 int16     `json:"rssi"`
	Temperature          int16     `json:"temperature"`
	SensitivityCable     int16     `json:"sensitivity_cable"`
	BuildVersion         int32     `json:"build_version"`
	CountCommandsInQueue int16     `json:"count_commands_in_queue"`
}

type ModemData struct {
	DevTime               time.Time `json:"dev_time"`
	RegTime               time.Time `json:"reg_time"`
	Modem                 int       `json:"modem"`
	Model                 int32     `json:"model"`
	Protocol              int16     `json:"protocol"`
	StatusRaw             int32     `json:"status_raw"`
	Status                int32     `json:"status"`
	Iccid                 string    `json:"iccid"`
	Timezone              int16     `json:"timezone"`
	Temperature           any       `json:"temperature"`
	BatteryVoltage        any       `json:"battery_voltage"`
	BatteryVoltageMin     any       `json:"battery_voltage_min"`
	Rssi                  int16     `json:"rssi"`
	Rspr                  int16     `json:"rspr"`
	Rsrq                  int16     `json:"rsrq"`
	Snr                   int16     `json:"snr"`
	Band                  int16     `json:"band"`
	Network               int16     `json:"network"`
	SatellitesCount       int16     `json:"satellites_count"`
	ModemErrorsCode       int32     `json:"modem_errors_code"`
	SealConnectPeriod     int32     `json:"seal_connect_period"`
	BatteryLevel          int16     `json:"battery_level"`
	ConnectPeriod         int32     `json:"connect_period"`
	RetriesCount          int16     `json:"retries_count"`
	SessionsCount         int64     `json:"sessions_count"`
	BatteryLevelMin       int16     `json:"battery_level_min"`
	MaxRegTime            int32     `json:"max_reg_time"`
	MaxSessionTime        int32     `json:"max_session_time"`
	SatelliteSearchPeriod int32     `json:"satellite_search_period"`
	PositioningTime       time.Time `json:"positioning_time"`
	Latitude              any       `json:"latitude"`
	Longitude             any       `json:"longitude"`
	Altitude              int32     `json:"altitude"`
}

type Modem struct {
	Id     int        `json:"id"`
	Serial uint64     `json:"serial"`
	Linked bool       `json:"linked"` //TODO delete later
	Last   *ModemData `json:"last"`
}

type Seal struct {
	Id      int     `json:"id"`
	Serial  uint64  `json:"serial"`
	Last    *Data   `json:"last"`
	Comment string  `json:"comment"`
	Modems  []Modem `json:"modems"`
}

type ArchiveQueryParams struct {
	Id        int
	From      time.Time `form:"from"`
	To        time.Time `form:"to"`
	Limit     int       `form:"limit"`
	OrderDesc bool      `form:"order_desc"`
}

type UpdateRequest struct {
	Comment *string `json:"comment,omitempty"`
}

type lastForList struct {
	Rssi                 int8      `json:"rssi"`
	Temperature          int8      `json:"temperature"`
	BatteryLevel         uint8     `json:"battery_level" db:"battery_level"`
	DevTime              time.Time `json:"dev_time" db:"dev_time"`
	Status               uint32    `json:"status"`
	CountCommandsInQueue uint8     `json:"count_commands_in_queue" db:"count_commands_in_queue"`
	BuildVersion         int32     `json:"build_version" db:"build_version"`
}
type SealForList struct {
	Id     int `json:"id"`
	Serial any `json:"serial"`
	Modems []struct {
		Id     int    `json:"id"`
		Serial uint64 `json:"serial"`
	} `json:"modems"`
	Last *lastForList `json:"last"`
}

type ArchiveSealDataModem struct {
	Id     int    `json:"id"`
	Serial uint64 `json:"serial"`
}
type ArchiveSealData struct {
	DevTime              time.Time            `json:"dev_time" db:"dev_time"`
	Modem                ArchiveSealDataModem `json:"modem"`
	Status               int64                `json:"status" db:"status"`
	Errors               int16                `json:"errors" db:"errors"`
	SensitivityRange     int16                `json:"sensitivity_range" db:"sensitivity_range"`
	BatteryLevel         int16                `json:"battery_level" db:"battery_level"`
	Rssi                 int16                `json:"rssi"`
	Temperature          int16                `json:"temperature"`
	CountCommandsInQueue int16                `json:"count_commands_in_queue" db:"count_commands_in_queue"`
}
