package secret_area

import (
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"time"
)

type Db struct {
	Id          int          `json:"id"`
	CreatedAt   *time.Time   `json:"created_at" db:"created_at"`
	Author      int          `json:"author"`
	Title       string       `json:"title" validate:"required,max=127,min=5"`
	Area        [][2]float32 `json:"area" validate:"required"`
	Description string       `json:"description" validate:"max=127"`
}

type QueryParams struct {
	transport.QueryParams
}

type Repo interface {
	Create(data Db) (SecretArea, error)
	Update(data Db) (SecretArea, error)
	GetDbById(id int) (Db, error)
	GetById(id int) (SecretArea, error)
	List(params QueryParams) (query.List[SecretArea], error)
	ExistsByUnique(id int, title string) (bool, error)
	DeleteById(id int) (bool, error)
}

type Usecase interface {
	Create(data CreateRequest, userId int) (SecretArea, error)
	Update(id int, data UpdateRequest) (SecretArea, error)
	GetDbById(id int) (Db, error)
	GetById(id int) (SecretArea, error)
	List(params QueryParams) (query.List[SecretArea], error)
	// Exists(id int) (bool, error)
	ExistsByUnique(id int, title string) (bool, error)
	DeleteById(id int) (bool, error)
}
