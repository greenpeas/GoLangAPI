package modemData

import (
	"time"
)

type sealData struct {
	DevTime time.Time `json:"dev_time" db:"dev_time"`
	Seal    struct {
		Id     int    `json:"id"`
		Serial uint64 `json:"serial"`
	} `json:"seal"`
	StatusRaw            uint32 `json:"status_raw" db:"status_raw"`
	Status               uint32 `json:"status" db:"status"`
	ErrorsRaw            uint8  `json:"errors_raw" db:"errors_raw"`
	Errors               uint8  `json:"errors" db:"errors"`
	SensitivityRange     uint8  `json:"sensitivity_range" db:"sensitivity_range"`
	BatteryLevel         uint8  `json:"battery_level" db:"battery_level"`
	Rssi                 int16  `json:"rssi"`
	CountCommandsInQueue uint8  `json:"count_commands_in_queue" db:"count_commands_in_queue"`
}

type Db struct {
	DevTime                  time.Time  `json:"dev_time" db:"dev_time"`
	RegTime                  time.Time  `json:"reg_time" db:"reg_time"`
	Modem                    int        `json:"modem"`
	Model                    uint16     `json:"model"`
	Protocol                 uint8      `json:"protocol"`
	StatusRaw                int32      `json:"status_raw" db:"status_raw"`
	Status                   int32      `json:"status"`
	Iccid                    string     `json:"iccid"`
	Timezone                 int8       `json:"timezone"`
	Temperature              any        `json:"temperature"`
	BatteryVoltage           any        `json:"battery_voltage" db:"battery_voltage"`
	BatteryVoltageMin        any        `json:"battery_voltage_min" db:"battery_voltage_min"`
	Rssi                     int16      `json:"rssi"`
	Rsrp                     int16      `json:"rsrp" `
	Rsrq                     int8       `json:"rsrq"`
	Snr                      int8       `json:"snr"`
	Band                     uint8      `json:"band"`
	Network                  uint8      `json:"network"`
	SatellitesCount          int16      `json:"satellites_count" db:"satellites_count"`
	ModemErrorsCode          uint16     `json:"modem_errors_code" db:"modem_errors_code"`
	SealConnectPeriod        uint16     `json:"seal_connect_period" db:"seal_connect_period"`
	BatteryLevel             int16      `json:"battery_level" db:"battery_level"`
	ConnectPeriod            uint16     `json:"connect_period" db:"connect_period"`
	RetriesCount             uint8      `json:"retries_count" db:"retries_count"`
	SessionsCount            uint32     `json:"sessions_count" db:"sessions_count"`
	BatteryLevelMin          uint8      `json:"battery_level_min" db:"battery_level_min"`
	MaxRegTime               uint16     `json:"max_reg_time" db:"max_reg_time"`
	MaxSessionTime           uint16     `json:"max_session_time" db:"max_session_time"`
	CoordinatesPeriod        uint16     `json:"coordinates_period" db:"coordinates_period"`
	PositioningTime          time.Time  `json:"positioning_time" db:"positioning_time"`
	Latitude                 any        `json:"latitude"`
	Longitude                any        `json:"longitude"`
	Altitude                 int32      `json:"altitude"`
	SatellitesSearchPeriod   uint16     `json:"satellites_search_period" db:"satellites_search_period"`
	LowPowerTimeout          uint16     `json:"low_power_timeout" db:"low_power_timeout"`
	SealsData                []sealData `json:"seals_data" db:"seals_data"`
	Speed                    int32      `json:"speed"`
	StatusGpsModule          int16      `json:"status_gps_module" db:"status_gps_module"`
	Hdop                     uint8      `json:"hdop" db:"hdop"`
	SignalGps                int32      `json:"signal_gps" db:"signal_gps"`
	SignalGlonass            int32      `json:"signal_glonass" db:"signal_glonass"`
	ErrorsFlags              int32      `json:"errors_flags" db:"errors_flags"`
	SensitivityAccelerometer int16      `json:"sensitivity_accelerometer" db:"sensitivity_accelerometer"`
}

type Repo interface {
	List(params ListParams) ([]ModemData, error)
}

type Usecase interface {
	List(params ListParams) ([]ModemData, error)
}
