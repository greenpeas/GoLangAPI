package shipping

import (
	"seal/internal/domain/route"
	"seal/internal/repository/pg/query"
	transp "seal/internal/transport"
	"time"
)

type File struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Comment  string `json:"comment,omitempty"`
	SealId   int    `json:"seal_id,omitempty"`
	Type     int    `json:"type"`
	Checksum string `json:"checksum"`
}

type Db struct {
	Id           int        `json:"id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	TimeStart    *time.Time `json:"time_start" db:"time_start"`
	TimeEnd      *time.Time `json:"time_end" db:"time_end"`
	Author       int        `json:"author"`
	CustomNumber string     `json:"custom_number" db:"custom_number" validate:"required,max=8,min=8"`
	CreateDate   string     `json:"create_date" db:"create_date" validate:"required,max=6,min=6"`
	Number       int        `json:"number" validate:"required,max=2147483647,min=1"`
	Transport    int        `json:"transport" validate:"required,max=2147483647,min=1"`
	Route        int        `json:"route" validate:"required,max=2147483647,min=1"`
	Status       int        `json:"status" validate:"max=2,min=0"`
	Files        []File     `json:"files" db:"files"`
	Modem        *int       `json:"modem" db:"modem"`
}

const STATUS_NEW = 0
const STATUS_ACTIVE = 1
const STATUS_END = 2

type QueryParams struct {
	transp.QueryParams
	Status []int `form:"status"`
}

type TelemetryQueryParams struct {
	Id    int
	From  time.Time `form:"from"`
	To    time.Time `form:"to"`
	Limit int       `form:"limit"`
}

type Repo interface {
	Create(data Db) (Shipping, error)
	Update(data Db) (Shipping, error)
	GetDbById(id int) (Db, error)
	GetById(id int) (Shipping, error)
	GetActiveByModemId(id int) (Shipping, error)
	List(params QueryParams) (query.List[ShippingForList], error)
	Exists(id int) (bool, error)
	ExistsByUnique(string, string, int, int) (bool, error)
	DeleteById(id int) (bool, error)
}

type Usecase interface {
	Create(data CreateRequest, userId int) (Shipping, error)
	Update(id int, data UpdateRequest) (Shipping, error)
	GetDbById(id int) (Db, error)
	GetById(id int) (Shipping, error)
	GetActiveByModemImei(imei uint64) (Shipping, error)
	List(params QueryParams) (query.List[ShippingForList], error)
	Exists(id int) (bool, error)
	ExistsByUnique(string, string, int, int) (bool, error)
	DeleteById(id int) (bool, error)
	Route(id int) (route.Route, error)
	Start(data Db) (Shipping, error)
	End(data Db) (Shipping, error)
	AddFiles(id int, files []File) (Shipping, error)
	RemoveFile(id int, name string) error
	RemoveFilesFromDisk(id int, files []File) error
	GetFilesDirectory(id int) string
	GetFileInfo(id int, name string) (File, error)
	UpdateFileInfo(id int, name string, data UpdateFileRequest) (Shipping, error)
	SetModemById(shippingId int, modemId int) (bool, error)
	SetModemByImei(shippingId int, imei uint64) (bool, error)
	Coordinates(params TrackQueryParams) ([]trackResponseCoordinate, error)
	Telemetry(params TrackQueryParams) ([]trackResponseTelemetry, error)
}
