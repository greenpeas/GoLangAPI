package tests

import (
	"seal/internal/app"
	"seal/internal/config"
	"seal/internal/tests/auth"
	"seal/internal/tests/custom"
	"seal/internal/tests/data"
	"seal/internal/tests/route"
	"seal/internal/tests/seal"
	"seal/internal/tests/seal_model"
	"seal/internal/tests/secret_area"
	"seal/internal/tests/shipping"
	transp "seal/internal/tests/transport"
	"seal/internal/tests/transport_type"
	"seal/internal/tests/user"
	"seal/internal/transport/rest"
	"testing"
	"time"
)

var appCore *app.App

// var jwt transport.JwtWithRefresh
var testData *data.TestData

func init() {
	cfg := config.Get("../../configs/app.yml")

	appCore = app.Get(cfg)
	//defer appCore.Clean()
	rest.Register(rest.Handler{
		Usecase:   appCore.Usecase,
		JwtWorker: appCore.JwtWorker,
		Logger:    appCore.Logger,
		Validator: appCore.Validator,
		Router:    appCore.Router,
	})

	testData = &data.TestData{App: appCore}
	testData.TimeStamp = time.Now().Unix()
}

func TestAuth(t *testing.T) {
	testData.Jwt = auth.Run(t, testData)
}

func TestUser(t *testing.T) {
	testData.User = user.Run(t, testData)
}

func TestCustom(t *testing.T) {
	testData.Custom = custom.Run(t, testData)
}

func TestRoute(t *testing.T) {
	testData.Route = route.Run(t, testData)
}

func TestTransportType(t *testing.T) {
	transport_type.Run(t, testData)
}

func TestTransport(t *testing.T) {
	testData.Transport = transp.Run(t, testData)
}

func TestSeal(t *testing.T) {
	testData.Seal = seal.Run(t, testData)
}

func TestSealModel(t *testing.T) {
	seal_model.Run(t, testData)
}

func TestSecretArea(t *testing.T) {
	secret_area.Run(t, testData)
}

func TestShipping(t *testing.T) {
	shipping.Run(t, testData)
}

func TestDelete(t *testing.T) {
	route.Delete(t, testData)
	seal.Delete(t, testData)
	custom.Delete(t, testData)
	transp.Delete(t, testData)
	user.Delete(t, testData)
}
