package modem

import (
	"fmt"
	app_interface "seal/internal/app/interface"
	"seal/internal/domain/command"
	modemData "seal/internal/domain/modem_data"
	modemLogRaw "seal/internal/domain/modem_log_raw"
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"seal/pkg/app_error"
	"seal/pkg/utils"
	"strconv"
	"time"
)

type cmds interface {
	Send(serial string, name string, params any, author string) (bool, error)
	List(serial string) (query.List[command.Command], error)
}

type CoreUseCase struct {
	Commands    cmds
	ModemData   modemData.Usecase
	ModemLogRaw modemLogRaw.Usecase
}

type usecase struct {
	repo      Repo
	logger    app_interface.Logger
	validator app_interface.Validator
	usecase   CoreUseCase
}

func NewUsecase(repo Repo, logger app_interface.Logger, validator app_interface.Validator, coreUseCase CoreUseCase) Usecase {
	return &usecase{repo, logger, validator, coreUseCase}
}

func (s *usecase) GetById(id int) (Modem, error) {
	return s.repo.GetById(id)
}

func (s *usecase) GetDbById(id int) (Db, error) {
	return s.repo.GetDbById(id)
}

func (s *usecase) Update(id int, data UpdateRequest) (Modem, error) {
	modem, err := s.GetDbById(id)

	if err != nil {
		return Modem{}, err
	}

	if err := utils.BindFromStruct(data, &modem); err != nil {
		return Modem{}, app_error.InternalServerError(err)
	}

	return s.repo.Update(modem)
}

func (s *usecase) GetByImei(imei uint64) (Modem, error) {
	return s.repo.GetByImei(imei)
}

func (s *usecase) List(queryParams transport.QueryParams) (query.List[ModemForList], error) {
	if errs := s.validator.Struct(queryParams); errs != nil {
		return query.List[ModemForList]{}, app_error.ValidationError(errs)
	}

	return s.repo.List(queryParams)
}

func (s *usecase) ListShippingReady(queryParams transport.QueryParams) (query.List[ModemForListShippingReady], error) {
	if errs := s.validator.Struct(queryParams); errs != nil {
		return query.List[ModemForListShippingReady]{}, app_error.ValidationError(errs)
	}

	return s.repo.ListShippingReady(queryParams)
}

func (s *usecase) SendCommand(id int, data SendCommandRequest, author string) (bool, error) {

	modem, err := s.GetById(id)

	if err != nil {
		return false, err
	}

	imei := fmt.Sprintf("%v", modem.Imei)

	return s.usecase.Commands.Send(imei, data.Name, data.Params, author)
}

func (s *usecase) CommandsList(id int) (query.List[command.Command], error) {

	modem, err := s.GetById(id)

	if err != nil {
		return query.List[command.Command]{}, err
	}

	imei := fmt.Sprintf("%v", modem.Imei)

	return s.usecase.Commands.List(imei)
}

func (s *usecase) Archive(params ArchiveQueryParams) ([]ArchiveModemData, error) {
	if errs := s.validator.Struct(params); errs != nil {
		return []ArchiveModemData{}, app_error.ValidationError(errs)
	}

	dataFromRepo, err := s.usecase.ModemData.List(modemData.ListParams{
		ModemId:   params.Id,
		TimeFrom:  params.From,
		TimeTo:    params.To,
		Limit:     params.Limit,
		OrderDesc: params.OrderDesc,
	})

	if err != nil {
		return []ArchiveModemData{}, err
	}

	archive := []ArchiveModemData{}

	for _, data := range dataFromRepo {
		archiveModemData := ArchiveModemData{
			DevTime:         data.DevTime,
			RegTime:         data.RegTime,
			Status:          data.Status,
			ErrorsFlags:     data.ErrorsFlags,
			PositioningTime: data.PositioningTime,
			Latitude:        data.Latitude,
			Longitude:       data.Longitude,
			Altitude:        data.Altitude,
			SatellitesCount: data.SatellitesCount,
			Speed:           data.Speed,
			StatusGpsModule: data.StatusGpsModule,
			Rssi:            data.Rssi,
			BatteryLevel:    data.BatteryLevel,
			SignalGps:       data.SignalGps,
			SignalGlonass:   data.SignalGlonass,
		}

		if data.CoordinatesLbs != nil {
			archiveModemData.CoordinatesLbs = &CoordinateLbs{
				PositioningTime: data.CoordinatesLbs.PositioningTime,
				Latitude:        data.CoordinatesLbs.Latitude,
				Longitude:       data.CoordinatesLbs.Longitude,
				Precision:       data.CoordinatesLbs.Precision,
			}
		}

		archive = append(archive, archiveModemData)
	}

	return archive, nil
}

func (s *usecase) LogRawTelemetry(params ArchiveQueryParams) ([]ArchiveModemData, error) {
	if errs := s.validator.Struct(params); errs != nil {
		return []ArchiveModemData{}, app_error.ValidationError(errs)
	}

	modem, err := s.GetDbById(params.Id)
	if err != nil {
		return []ArchiveModemData{}, err
	}

	dataFromRepo, err := s.usecase.ModemLogRaw.ListTelemetry(modemLogRaw.ListParams{
		Imei:      strconv.FormatUint(modem.Imei, 10),
		From:      params.From,
		To:        params.To,
		Limit:     params.Limit,
		OrderDesc: params.OrderDesc,
	})

	if err != nil {
		return []ArchiveModemData{}, err
	}

	archive := []ArchiveModemData{}

	for _, data := range dataFromRepo {
		archiveModemData := ArchiveModemData{
			DevTime:         time.Unix(data.CurrentTime, 0),
			RegTime:         data.RegTime,
			Status:          data.Status,
			ErrorsFlags:     data.ErrorsFlags,
			PositioningTime: time.Unix(data.PositioningTime, 0),
			Latitude:        data.Latitude,
			Longitude:       data.Longitude,
			Altitude:        data.Altitude,
			SatellitesCount: data.SatellitesCount,
			Speed:           data.Speed,
			StatusGpsModule: data.StatusGpsModule,
			Rssi:            data.Rssi,
			BatteryLevel:    data.BatteryLevel,
			SignalGps:       data.SignalGps,
			SignalGlonass:   data.SignalGlonass,
		}

		archive = append(archive, archiveModemData)
	}

	return archive, nil
}

func (s *usecase) Track(params TrackQueryParams) (TrackResponse, error) {
	if errs := s.validator.Struct(params); errs != nil {
		return TrackResponse{}, app_error.ValidationError(errs)
	}

	coords, err := s.repo.Track(params)

	return TrackResponse{Coordinates: coords}, err
}

func (s *usecase) TrackLbs(params TrackQueryParams) ([]CoordinateLbs, error) {
	if errs := s.validator.Struct(params); errs != nil {
		return []CoordinateLbs{}, app_error.ValidationError(errs)
	}

	return s.repo.TrackLbs(params)
}

func (s *usecase) Log(params LogQueryParams) ([]modemLogRaw.ModemLogRaw, error) {
	if errs := s.validator.Struct(params); errs != nil {
		return []modemLogRaw.ModemLogRaw{}, app_error.ValidationError(errs)
	}

	modem, err := s.GetById(params.Id)

	if err != nil {
		return []modemLogRaw.ModemLogRaw{}, err
	}

	return s.usecase.ModemLogRaw.List(modemLogRaw.ListParams{
		Imei:      modem.Imei,
		From:      params.From,
		To:        params.To,
		Limit:     params.Limit,
		OrderDesc: params.OrderDesc,
	})

}
