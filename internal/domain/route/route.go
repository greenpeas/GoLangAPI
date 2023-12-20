package route

import (
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"time"
)

type Db struct {
	Id         int       `json:"id"`
	Title      string    `json:"title" validate:"required"`
	Points     []string  `json:"points" validate:"required"`
	Length     int       `json:"length" validate:"max=32767,min=0"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	TravelTime int       `json:"travel_time" db:"travel_time"`
}

func (r Db) getId() int {
	return r.Id
}

func (r Db) getTitle() string {
	return r.Title
}

type Repo interface {
	Create(data Route) (Route, error)
	Update(data Db) (Route, error)
	GetById(id int) (Route, error)
	GetDbById(id int) (Db, error)
	List(params transport.QueryParams) (query.List[Route], error)
	Exists(id int) (bool, error)
	ExistsByUnique(id int, title string) (bool, error)
	DeleteById(id int) (bool, error)
}

type Usecase interface {
	Create(data CreateRequest) (Route, error)
	Update(id int, data UpdateRequest) (Route, error)
	GetById(id int) (Route, error)
	GetDbById(id int) (Db, error)
	List(params transport.QueryParams) (query.List[Route], error)
	Exists(id int) (bool, error)
	ExistsByUnique(id int, title string) (bool, error)
	DeleteById(id int) (bool, error)
}
