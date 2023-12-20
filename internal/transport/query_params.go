package transport

type QueryParams struct {
	Find     string `form:"find,omitempty" validate:"max=1024"`
	FindType uint8  `form:"find_type,omitempty" validate:"max=4"`
	Limit    int    `form:"limit,omitempty" validate:"max=100"`
	Offset   uint16 `form:"offset,omitempty" validate:"max=32767"`
}
