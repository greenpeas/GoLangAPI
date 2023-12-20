package transport_type

import (
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"time"
)

type Db struct {
	Id        int       `json:"id"`
	Title     string    `json:"title" validate:"required,max=50,min=5"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type TransportType struct {
	Id        int        `json:"id"`
	Title     string     `json:"title"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
}
type Repo interface {
	List(transport.QueryParams) (query.List[TransportType], error)
}

type Usecase interface {
	List(transport.QueryParams) (query.List[TransportType], error)
}
