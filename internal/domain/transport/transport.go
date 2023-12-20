package transport

import (
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"time"
)

type Db struct {
	Id                 int       `json:"id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	Author             int       `json:"author" validate:"required"`
	Title              string    `json:"title" validate:"required,max=50,min=5"`
	Type               int       `json:"type" validate:"required,max=3,min=1"`
	RegistrationNumber string    `json:"registration_number"`
}

type Repo interface {
	Create(data Db) (Transport, error)
	Update(data Db) (Transport, error)
	GetById(id int) (Transport, error)
	GetDbById(id int) (Db, error)
	List(params transport.QueryParams) (query.List[Transport], error)
	Exists(id int) (bool, error)
	ExistsByUnique(id int, title string) (bool, error)
	DeleteById(id int) (bool, error)
}

type Usecase interface {
	Create(data CreateRequest, userId int) (Transport, error)
	Update(id int, data UpdateRequest) (Transport, error)
	GetById(id int) (Transport, error)
	GetDbById(id int) (Db, error)
	List(params transport.QueryParams) (query.List[Transport], error)
	Exists(id int) (bool, error)
	ExistsByUnique(id int, title string) (bool, error)
	DeleteById(id int) (bool, error)
}
