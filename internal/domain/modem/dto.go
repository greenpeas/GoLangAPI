package modem

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type Data struct {
	DevTime                  time.Time      `json:"dev_time"`
	RegTime                  time.Time      `json:"reg_time"`
	Modem                    int            `json:"modem"`
	Model                    uint16         `json:"model"`
	Protocol                 uint8          `json:"protocol"`
	StatusRaw                uint16         `json:"status_raw"`
	Status                   uint16         `json:"status"`
	Iccid                    string         `json:"iccid"`
	Timezone                 int8           `json:"timezone"`
	Temperature              any            `json:"temperature"`
	BatteryVoltage           any            `json:"battery_voltage"`
	BatteryVoltageMin        any            `json:"battery_voltage_min"`
	Rssi                     int16          `json:"rssi"`
	Rsrp                     int16          `json:"rsrp"`
	Rsrq                     int8           `json:"rsrq"`
	Snr                      int8           `json:"snr"`
	Band                     uint8          `json:"band"`
	Network                  uint8          `json:"network"`
	SatellitesCount          uint8          `json:"satellites_count"`
	ModemErrorsCode          uint16         `json:"modem_errors_code"`
	SealConnectPeriod        uint16         `json:"seal_connect_period"`
	BatteryLevel             int16          `json:"battery_level"`
	ConnectPeriod            uint16         `json:"connect_period"`
	RetriesCount             uint8          `json:"retries_count"`
	SessionsCount            uint32         `json:"sessions_count"`
	BatteryLevelMin          uint8          `json:"battery_level_min"`
	MaxRegTime               uint16         `json:"max_reg_time"`
	MaxSessionTime           uint16         `json:"max_session_time"`
	CoordinatesPeriod        uint16         `json:"coordinates_period"`
	PositioningTime          time.Time      `json:"positioning_time"`
	Latitude                 any            `json:"latitude"`
	Longitude                any            `json:"longitude"`
	Altitude                 int32          `json:"altitude"`
	Speed                    uint16         `json:"speed"`
	StatusGpsModule          uint8          `json:"status_gps_module"`
	Hdop                     uint8          `json:"hdop"`
	SignalGps                int32          `json:"signal_gps"`
	SignalGlonass            int32          `json:"signal_glonass"`
	ErrorsFlags              uint16         `json:"errors_flags"`
	SatellitesSearchPeriod   uint16         `json:"satellites_search_period"`
	LowPowerTimeout          uint16         `json:"low_power_timeout"`
	CoordinatesLbs           *CoordinateLbs `json:"coordinate_lbs"`
	SensitivityAccelerometer int16          `json:"sensitivity_accelerometer"`
}

type Extra struct {
	Network          string    `json:"network"`
	Host             string    `json:"host"`
	Apn              string    `json:"apn"`
	Dns1             string    `json:"dns_1"`
	Dns2             string    `json:"dns_2"`
	Ip               string    `json:"ip"`
	SoftwareVersion  string    `json:"software_version"`
	HardwareRevision string    `json:"hardware_revision"`
	ModemRevision    string    `json:"modem_revision"`
	GpsRevision      string    `json:"gps_revision"`
	RfRevision       string    `json:"rf_revision"`
	ListOperators    string    `json:"list_operators"`
	UpdatedTime      time.Time `json:"updated_time"`
}

type SealData struct {
	Rssi                 int8      `json:"rssi"`
	Temperature          int8      `json:"temperature"`
	BatteryLevel         uint8     `json:"battery_level" db:"battery_level"`
	DevTime              time.Time `json:"dev_time" db:"dev_time"`
	Status               uint32    `json:"status"`
	CountCommandsInQueue uint8     `json:"count_commands_in_queue" db:"count_commands_in_queue"`
	BuildVersion         int32     `json:"build_version" db:"build_version"`
}
type Seal struct {
	Id     int       `json:"id"`
	Serial uint64    `json:"serial"`
	Last   *SealData `json:"last"`
}

type Modem struct {
	Id             int                `json:"id"`
	Imei           string             `json:"imei"`
	Serial         uint64             `json:"serial"`
	Iccid          string             `json:"iccid"`
	LastDevTime    pgtype.Timestamptz `json:"last_dev_time" db:"last_dev_time" swaggertype:"string"`
	Extra          *Extra             `json:"extra"`
	Last           *Data              `json:"last"`
	SerialsOfSeals []uint64           `json:"-" db:"serials_of_seals"`
	Seals          []Seal             `json:"seals"`
	LastChargeTime pgtype.Timestamptz `json:"last_charge_time" db:"last_charge_time" swaggertype:"string"`
	Comment        string             `json:"comment"`
	Msisdn         *string            `json:"msisdn"`
	LastCoordinate *Coordinate        `json:"last_coordinate" db:"last_coordinate"`
}

type SendCommandRequest struct {
	Name   string `json:"name" validate:"required,max=50,min=5"`
	Params any    `json:"params" validate:"required"`
}

type ArchiveQueryParams struct {
	Id        int
	From      time.Time `form:"from"`
	To        time.Time `form:"to"`
	Limit     int       `form:"limit"`
	OrderDesc bool      `form:"order_desc"`
}

type TrackQueryParams struct {
	Id        int
	From      time.Time `form:"from"`
	To        time.Time `form:"to"`
	Limit     int       `form:"limit"`
	OrderDesc bool      `form:"order_desc"`
}

type Coordinate struct {
	Latitude           float32   `db:"latitude" json:"latitude"`
	Longitude          float32   `db:"longitude" json:"longitude"`
	Altitude           int32     `db:"altitude" json:"altitude"`
	DevTime            time.Time `db:"dev_time" json:"dev_time"`
	SatellitesCount    uint8     `db:"satellites_count" json:"satellites_count"`
	Speed              uint16    `db:"speed" json:"speed"`
	StatusGpsModule    uint8     `db:"status_gps_module" json:"status_gps_module"`
	Hdop               uint8     `db:"hdop" json:"hdop"`
	SignalGps          int32     `db:"signal_gps" json:"signal_gps"`
	SignalGlonass      int32     `db:"signal_glonass" json:"signal_glonass"`
	MinDistanceToRoute *int      `db:"min_distance_to_route" json:"min_distance_to_route"`
}

type CoordinateLbs struct {
	DevTime         time.Time `db:"dev_time" json:"-"`
	Modem           int       `db:"modem" json:"-"`
	PositioningTime time.Time `db:"positioning_time" json:"positioning_time"`
	Latitude        float32   `db:"latitude" json:"latitude"`
	Longitude       float32   `db:"longitude" json:"longitude"`
	Precision       int       `db:"precision" json:"precision"`
}

type TrackResponse struct {
	Coordinates []Coordinate `json:"coordinates"`
}

type LogQueryParams struct {
	Id        int       `validate:"required"`
	From      time.Time `form:"from" validate:"required"`
	To        time.Time `form:"to"`
	Limit     int       `form:"limit"`
	OrderDesc bool      `form:"order_desc"`
}

type UpdateRequest struct {
	Comment *string `json:"comment,omitempty"`
}

type lastForList struct {
	RegTime       time.Time `json:"reg_time" db:"reg_time"`
	Status        int32     `json:"status"`
	ErrorsFlags   int32     `json:"errors_flags"`
	Rssi          int16     `json:"rssi"`
	ConnectPeriod int32     `json:"connect_period" db:"connect_period"`
	BatteryLevel  int16     `json:"battery_level" db:"battery_level"`
	BuildVersion  int32     `json:"build_version" db:"build_version"`
}
type coordinateForList struct {
	DevTime   time.Time `json:"dev_time"`
	Latitude  float32   `json:"latitude"`
	Longitude float32   `json:"longitude"`
}
type ModemForList struct {
	Id               int                `json:"id"`
	Imei             string             `json:"imei"`
	Serial           uint64             `json:"serial"`
	Iccid            string             `json:"iccid"`
	LastDevTime      pgtype.Timestamptz `json:"last_dev_time" db:"last_dev_time" swaggertype:"string"`
	Last             *lastForList       `json:"last"`
	Comment          string             `json:"comment"`
	SoftwareVersion  *string            `json:"software_version" db:"software_version"`
	HardwareRevision *string            `json:"hardware_revision" db:"hardware_revision"`
	LastCoordinate   *coordinateForList `json:"last_coordinate" db:"last_coordinate"`
}

type ModemForListShippingReady struct {
	Id     int `json:"id"`
	Serial any `json:"serial,string"`
	Seals  []struct {
		Id     int `json:"id"`
		Serial any `json:"serial,string"`
		Last   struct {
			Status       int32 `json:"status"`
			BatteryLevel int16 `json:"battery_level"`
		} `json:"last"`
	} `json:"seals"`
}

type ArchiveModemData struct {
	DevTime         time.Time      `json:"dev_time"`
	RegTime         time.Time      `json:"reg_time"`
	Status          int32          `json:"status"`
	ErrorsFlags     int32          `json:"errors_flags"`
	PositioningTime time.Time      `json:"positioning_time"`
	Latitude        any            `json:"latitude"`
	Longitude       any            `json:"longitude"`
	Altitude        int32          `json:"altitude"`
	SatellitesCount int16          `json:"satellites_count"`
	Speed           int32          `json:"speed"`
	StatusGpsModule int16          `json:"status_gps_module"`
	Rssi            int16          `json:"rssi"`
	BatteryLevel    int16          `json:"battery_level"`
	CoordinatesLbs  *CoordinateLbs `json:"coordinate_lbs"`
	SignalGps       int32          `json:"signal_gps"`
	SignalGlonass   int32          `json:"signal_glonass"`
}
