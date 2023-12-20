package sealData

import (
	"time"
)

type Db struct {
	DevTime              time.Time `json:"dev_time" db:"dev_time"`
	Seal                 int       `json:"seal"`
	ModemTime            time.Time `json:"modem_time" db:"modem_time"`
	Modem                int       `json:"modem"`
	StatusRaw            int64     `json:"status_raw" db:"status_raw"`
	Status               int64     `json:"status" db:"status"`
	ErrorsRaw            int16     `json:"errors_raw" db:"errors_raw"`
	Errors               int16     `json:"errors" db:"errors"`
	SensitivityRange     int16     `json:"sensitivity_range" db:"sensitivity_range"`
	BatteryLevel         int16     `json:"battery_level" db:"battery_level"`
	Rssi                 int16     `json:"rssi"`
	Temperature          int16     `json:"temperature"`
	SensitivityCable     int16     `json:"sensitivity_cable" db:"sensitivity_cable"`
	BuildVersion         int32     `json:"build_version" db:"build_version"`
	CountCommandsInQueue int16     `json:"count_commands_in_queue"`
}

type Repo interface {
	List(params ListParams) ([]SealData, error)
}

type Usecase interface {
	List(params ListParams) ([]SealData, error)
}
