package usecase

import (
	"context"
	app_interface "seal/internal/app/interface"
	"seal/internal/config"
	"seal/internal/domain/custom"
	"seal/internal/domain/modem"
	modemData "seal/internal/domain/modem_data"
	modemLogRaw "seal/internal/domain/modem_log_raw"
	"seal/internal/domain/route"
	"seal/internal/domain/seal"
	sealData "seal/internal/domain/seal_data"
	"seal/internal/domain/seal_model"
	"seal/internal/domain/secret_area"
	"seal/internal/domain/shipping"
	"seal/internal/domain/transport"
	"seal/internal/domain/transport_type"
	"seal/internal/domain/user"
	"seal/internal/transport/commands"
	"seal/pkg/hash"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Usecase struct {
	User          user.Usecase
	Route         route.Usecase
	Custom        custom.Usecase
	Seal          seal.Usecase
	SealData      sealData.Usecase
	SealModel     seal_model.Usecase
	SecretArea    secret_area.Usecase
	Shipping      shipping.Usecase
	Transport     transport.Usecase
	TransportType transport_type.Usecase
	Modem         modem.Usecase
	ModemData     modemData.Usecase
	ModemLogRaw   modemLogRaw.Usecase
}

type Params struct {
	Ctx       context.Context
	Db        *pgxpool.Pool
	Cfg       *config.Config
	Logger    app_interface.Logger
	Validator app_interface.Validator
}

func GetUsecase(params Params) *Usecase {
	var usecase Usecase

	userRepo := user.NewRepo(params.Ctx, params.Db, params.Logger)
	usecase.User = user.NewUsecase(userRepo, hash.NewSHA1Hasher(params.Cfg.Salt), params.Logger, params.Validator)

	customRepo := custom.NewRepo(params.Ctx, params.Db, params.Logger)
	usecase.Custom = custom.NewUsecase(customRepo, params.Logger, params.Validator)

	routeRepo := route.NewRepo(params.Ctx, params.Db, params.Logger)
	usecase.Route = route.NewUsecase(routeRepo, params.Logger, usecase.Custom, params.Validator)

	transportRepo := transport.NewRepo(params.Ctx, params.Db, params.Logger)
	usecase.Transport = transport.NewUsecase(transportRepo, params.Logger, params.Validator)

	transportTypeRepo := transport_type.NewRepo(params.Ctx, params.Db, params.Logger)
	usecase.TransportType = transport_type.NewUsecase(transportTypeRepo, params.Logger, params.Validator)

	sealModelRepo := seal_model.NewRepo(params.Ctx, params.Db, params.Logger)
	usecase.SealModel = seal_model.NewUsecase(sealModelRepo, params.Logger, params.Validator)

	secretAreaRepo := secret_area.NewRepo(params.Ctx, params.Db, params.Logger)
	usecase.SecretArea = secret_area.NewUsecase(secretAreaRepo, params.Logger, params.Validator, secret_area.CoreUseCase{User: usecase.User})

	modemDataRepo := modemData.NewRepo(params.Ctx, params.Db, params.Logger)
	usecase.ModemData = modemData.NewUsecase(modemDataRepo, params.Logger, params.Validator)

	modemLogRawRepo := modemLogRaw.NewRepo(params.Ctx, params.Db, params.Logger)
	usecase.ModemLogRaw = modemLogRaw.NewUsecase(modemLogRawRepo, params.Logger, params.Validator)

	cmd := commands.NewGRPC(params.Ctx, params.Cfg.GRPCCommands.Addr, params.Cfg.GRPCCommands.Timeout, params.Logger)

	modemRepo := modem.NewRepo(params.Ctx, params.Db, params.Logger)
	usecase.Modem = modem.NewUsecase(modemRepo, params.Logger, params.Validator, modem.CoreUseCase{
		Commands:    cmd,
		ModemData:   usecase.ModemData,
		ModemLogRaw: usecase.ModemLogRaw,
	})

	shippingRepo := shipping.NewRepo(params.Ctx, params.Db, params.Logger)
	usecase.Shipping = shipping.NewUsecase(shippingRepo, params.Logger, params.Validator, params.Cfg.ShippingFilesPath, shipping.CoreUseCase{
		User:      usecase.User,
		Route:     usecase.Route,
		Seal:      usecase.Seal,
		Transport: usecase.Transport,
		Modem:     usecase.Modem,
		ModemData: usecase.ModemData,
	})

	sealDataRepo := sealData.NewRepo(params.Ctx, params.Db, params.Logger)
	usecase.SealData = sealData.NewUsecase(sealDataRepo, params.Logger, params.Validator)

	sealRepo := seal.NewRepo(params.Ctx, params.Db, params.Logger)
	usecase.Seal = seal.NewUsecase(sealRepo, params.Logger, params.Validator, seal.CoreUseCase{SealData: usecase.SealData})

	return &usecase
}
