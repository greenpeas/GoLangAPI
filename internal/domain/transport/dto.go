package transport

import (
	"seal/internal/domain/transport_type"
	"seal/internal/domain/user"
	"time"
)

type Transport struct {
	Id                 int                          `json:"id"`
	CreatedAt          *time.Time                   `json:"created_at" db:"created_at"`
	Title              string                       `json:"title"`
	Author             *user.Author                 `json:"author"`
	Type               transport_type.TransportType `json:"type"`
	RegistrationNumber string                       `json:"registration_number" db:"registration_number"`
}

type CreateRequest struct {
	Title              string `json:"title" validate:"required,max=50,min=5"`
	Type               int    `json:"type" validate:"required,max=3,min=1"`
	RegistrationNumber string `json:"registration_number"`
}

type UpdateRequest struct {
	Title              *string `json:"title,omitempty" validate:"max=50"`
	Type               *int    `json:"type,omitempty" validate:"max=3"`
	RegistrationNumber *string `json:"registration_number"`
}
