package custom

type Custom = Db

type CreateRequest struct {
	Title string `json:"title" validate:"required,max=50,min=1"`
}

type UpdateRequest struct {
	Title string `json:"title,omitempty" validate:"max=50"`
}
