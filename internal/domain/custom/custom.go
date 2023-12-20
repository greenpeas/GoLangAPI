package custom

import (
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"time"
)

type Db struct {
	Id        int       `json:"id"`
	Title     string    `json:"title" validate:"required,max=50,min=1"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Repo interface {
	Create(custom Db) (Custom, error)
	Update(data Db) (Custom, error)
	GetById(id int) (Custom, error)
	GetDbById(id int) (Db, error)
	List(params transport.QueryParams) (query.List[Custom], error)
	Exists(id int) (bool, error)
	ExistsByUnique(id int, title string) (bool, error)
	DeleteById(id int) (bool, error)
}

type Usecase interface {
	Create(data CreateRequest) (Custom, error)
	Update(id int, data UpdateRequest) (Custom, error)
	GetById(id int) (Custom, error)
	GetDbById(id int) (Db, error)
	List(params transport.QueryParams) (query.List[Custom], error)
	Exists(id int) (bool, error)
	ExistsByUnique(id int, title string) (bool, error)
	DeleteById(id int) (bool, error)
}
