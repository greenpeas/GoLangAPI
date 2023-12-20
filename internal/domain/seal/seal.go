package seal

import (
	"github.com/jackc/pgx/v5/pgtype"
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
)

type Db struct {
	Id          int                `json:"id"`
	Serial      uint64             `json:"serial"`
	LastDevTime pgtype.Timestamptz `json:"last_dev_time" db:"last_dev_time"`
	Comment     string             `json:"comment"`
}

type Repo interface {
	GetById(id int) (Seal, error)
	GetDbById(id int) (Db, error)
	List(params transport.QueryParams) (query.List[SealForList], error)
	Exists(id int) (bool, error)
	Update(data Db) (Seal, error)
}

type Usecase interface {
	GetById(id int) (Seal, error)
	GetDbById(id int) (Db, error)
	List(params transport.QueryParams) (query.List[SealForList], error)
	Exists(id int) (bool, error)
	Archive(params ArchiveQueryParams) ([]ArchiveSealData, error)
	Update(id int, data UpdateRequest) (Seal, error)
}
