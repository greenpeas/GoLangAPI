package main

import (
	"seal/internal/app"
	"seal/internal/config"
)

// @title           Seal API
// @version         0.1
// @description     Описание протокола обмена сообщениями

// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	cfg := config.Get("configs/app.yml")

	app := app.Get(cfg)
	defer app.Clean()

	app.Run()

}
