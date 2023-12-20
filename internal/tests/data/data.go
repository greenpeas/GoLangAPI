package data

import (
	"seal/internal/app"
	"seal/internal/domain/custom"
	"seal/internal/domain/route"
	"seal/internal/domain/seal"
	"seal/internal/domain/transport"
	"seal/internal/domain/user"
	transp "seal/internal/transport"
)

type TestData struct {
	App       *app.App
	Jwt       transp.JwtWithRefresh
	Route     route.Route
	Seal      seal.Seal
	Transport transport.Transport
	User      user.User
	Custom    custom.Custom
	TimeStamp int64
}
