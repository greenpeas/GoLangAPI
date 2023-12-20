package modem

import (
	"github.com/jackc/pgx/v5/pgtype"
	"seal/internal/domain/command"
	modemLogRaw "seal/internal/domain/modem_log_raw"
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
)

type Db struct {
	Id             int                `json:"id"`
	Imei           uint64             `json:"imei,string"`
	Serial         uint64             `json:"serial"`
	Iccid          string             `json:"iccid"`
	LastDevTime    pgtype.Timestamptz `json:"last_dev_time" db:"last_dev_time"`
	Extra          *Extra             `json:"extra"`
	SerialsOfSeals []uint64           `json:"-" db:"serials_of_seals"`
	Comment        string             `json:"comment"`
	Msisdn         *string            `json:"msisdn"`
}

type Repo interface {
	GetById(id int) (Modem, error)
	GetDbById(id int) (Db, error)
	GetByImei(imei uint64) (Modem, error)
	List(params transport.QueryParams) (query.List[ModemForList], error)
	ListShippingReady(params transport.QueryParams) (query.List[ModemForListShippingReady], error)
	Track(params TrackQueryParams) ([]Coordinate, error)
	TrackLbs(params TrackQueryParams) ([]CoordinateLbs, error)
	Update(data Db) (Modem, error)
}

type Usecase interface {
	GetById(id int) (Modem, error)
	GetDbById(id int) (Db, error)
	GetByImei(imei uint64) (Modem, error)
	List(params transport.QueryParams) (query.List[ModemForList], error)
	ListShippingReady(params transport.QueryParams) (query.List[ModemForListShippingReady], error)
	SendCommand(sealId int, data SendCommandRequest, author string) (bool, error)
	CommandsList(sealId int) (query.List[command.Command], error)
	Archive(params ArchiveQueryParams) ([]ArchiveModemData, error)
	LogRawTelemetry(params ArchiveQueryParams) ([]ArchiveModemData, error)
	Track(params TrackQueryParams) (TrackResponse, error)
	TrackLbs(params TrackQueryParams) ([]CoordinateLbs, error)
	Log(params LogQueryParams) ([]modemLogRaw.ModemLogRaw, error)
	Update(id int, data UpdateRequest) (Modem, error)
}
