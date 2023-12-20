package secret_area

import (
	"seal/internal/domain/user"
	"time"
)

type SecretArea struct {
	Id          int          `json:"id"`
	CreatedAt   *time.Time   `json:"created_at" db:"created_at"`
	Author      *user.Author `json:"author"`
	Title       string       `json:"title"`
	Area        [][2]float32 `json:"area"`
	Description string       `json:"description"`
}

type CreateRequest struct {
	Title       string       `json:"title"  validate:"required,max=127,min=5"`
	Description string       `json:"description" validate:"max=127"`
	Area        [][2]float32 `json:"area" validate:"required"`
}

type UpdateRequest struct {
	Title       string       `json:"title,omitempty"  validate:"max=127"`
	Description string       `json:"description,omitempty" validate:"max=127"`
	Area        [][2]float32 `json:"area,omitempty"`
}
