package user

import (
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"time"
)

type Db struct {
	Id        int       `json:"id"`
	Login     string    `json:"login" validate:"required,max=50,min=5"`
	Password  string    `json:"password" validate:"required,min=6"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Role      int       `json:"role"`
	Title     string    `json:"title" validate:"required,min=3,max=50"`
}

type Repo interface {
	Create(data Db) (User, error)
	Update(data Db) (User, error)
	GetByCredentials(login, password string) (User, error)
	GetById(id int) (User, error)
	GetDbById(id int) (Db, error)
	List(params transport.QueryParams) (query.List[User], error)
	Exists(id int) (bool, error)
	ExistsByUnique(id int, title string) (bool, error)
	DeleteById(id int) (bool, error)
}

type Usecase interface {
	Create(data CreateRequest) (User, error)
	Update(id int, data UpdateRequest) (User, error)
	GetByCredentials(login, password string) (User, error)
	GetById(id int) (User, error)
	GetDbById(id int) (Db, error)
	List(params transport.QueryParams) (query.List[User], error)
	Exists(id int) (bool, error)
	ExistsByUnique(id int, title string) (bool, error)
	DeleteById(id int) (bool, error)
}
