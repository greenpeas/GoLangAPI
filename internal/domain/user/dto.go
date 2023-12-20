package user

import "time"

type User struct {
	Id        int       `json:"id"`
	Login     string    `json:"login"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Role      int       `json:"role"`
	Title     string    `json:"title"`
}

type Author struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

type CreateRequest struct {
	Login    string `json:"login" validate:"required,max=50,min=5"`
	Password string `json:"password" validate:"required,min=6"`
	Role     int    `json:"role" validate:"required,max=6,min=0"`
	Title    string `json:"title" validate:"required,max=50,min=3"`
}

type UpdateRequest struct {
	Login    string `json:"login,omitempty" validate:"max=50"`
	Password string `json:"password,omitempty" validate:"max=6"`
	Role     int    `json:"role,omitempty" validate:"max=6"`
	Title    string `json:"title,omitempty" validate:"max=50"`
}
