package seal_model

import (
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"time"
)

type Db struct {
	Id        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	Title     string     `json:"title" validate:"required,max=50,min=5"`
}

type Repo interface {
	List(transport.QueryParams) (query.List[SealModel], error)
}

type Usecase interface {
	List(transport.QueryParams) (query.List[SealModel], error)
}
