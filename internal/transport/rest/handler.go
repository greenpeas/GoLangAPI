package rest

import (
	app_interface "seal/internal/app/interface"
	"seal/internal/app/usecase"
	"seal/internal/config"
	"seal/internal/transport"
	"seal/internal/transport/rest/middleware"
	v1 "seal/internal/transport/rest/v1"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Usecase   *usecase.Usecase
	JwtWorker transport.JwtWorker
	Logger    app_interface.Logger
	Validator app_interface.Validator
	Router    *gin.Engine
	Cfg       *config.Config
}

type Params struct {
	Handler
	Router *gin.Engine
}

func Register(handler Handler) *Handler {
	handler.register(handler.Router)

	return &handler
}

func (h *Handler) register(router *gin.Engine) {
	h.Logger.Info("Register rest handler")
	router.Use(middleware.Error())
	handlerV1 := v1.Handler{
		Usecase:   h.Usecase,
		JwtWorker: h.JwtWorker,
		Logger:    h.Logger,
		Validator: h.Validator,
		Cfg:       h.Cfg,
	}

	api := router.Group("/api")
	{
		handlerV1.Register(api)
	}
}
