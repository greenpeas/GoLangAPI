package v1

import (
	app_interface "seal/internal/app/interface"
	"seal/internal/app/usecase"
	"seal/internal/config"
	"seal/internal/transport"
	"seal/internal/transport/rest/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Usecase   *usecase.Usecase
	JwtWorker transport.JwtWorker
	Logger    app_interface.Logger
	Validator app_interface.Validator
	Cfg       *config.Config
}

func (h *Handler) Register(api *gin.RouterGroup) {
	h.Logger.Info("Register rest v1 handler")
	v1 := api.Group("/v1")
	{
		h.registerAuthHandler(v1)
		h.registerShippingTempHandler(v1)
		v1.Use(middleware.Auth(h.JwtWorker))
		h.registerCustomHandler(v1)
		h.registerRouteHandler(v1)
		h.registerSealHandler(v1)
		h.registerSealModelHandler(v1)
		h.registerSecretAreaHandler(v1)
		h.registerShippingHandler(v1)
		h.registerTransportHandler(v1)
		h.registerTransportTypeHandler(v1)
		h.registerUserHandler(v1)
		h.registerModemHandler(v1)

	}
}
