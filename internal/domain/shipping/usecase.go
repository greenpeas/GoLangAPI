package shipping

import (
	"fmt"
	"os"
	app_interface "seal/internal/app/interface"
	"seal/internal/domain/modem"
	modemData "seal/internal/domain/modem_data"
	"seal/internal/domain/route"
	"seal/internal/domain/seal"
	transp "seal/internal/domain/transport"
	"seal/internal/domain/user"
	"seal/internal/repository/pg/query"
	"seal/pkg/app_error"
	"seal/pkg/utils"
	"strconv"
	"time"
)

type CoreUseCase struct {
	User      user.Usecase
	Route     route.Usecase
	Seal      seal.Usecase
	Transport transp.Usecase
	Modem     modem.Usecase
	ModemData modemData.Usecase
}

type usecase struct {
	repo      Repo
	logger    app_interface.Logger
	validator app_interface.Validator
	usecase   CoreUseCase
	filesPath string
}

func NewUsecase(repo Repo, logger app_interface.Logger, validator app_interface.Validator, filesPath string, coreUsecase CoreUseCase) Usecase {
	return &usecase{repo, logger, validator, coreUsecase, filesPath}
}

func (s *usecase) Create(data CreateRequest, userId int) (Shipping, error) {
	var shipping Db

	if err := utils.BindFromStruct(data, &shipping); err != nil {
		return Shipping{}, app_error.InternalServerError(err)
	}

	shipping.Author = userId

	if errs, err := s.validate(shipping); err != nil {
		return Shipping{}, err
	} else if len(errs) > 0 {
		return Shipping{}, app_error.ValidationError(errs)
	}

	return s.repo.Create(shipping)
}

func (s *usecase) Update(id int, data UpdateRequest) (Shipping, error) {
	shipping, err := s.GetDbById(id)

	if err != nil {
		return Shipping{}, app_error.ErrNotFound
	}

	if err := utils.BindFromStruct(data, &shipping); err != nil {
		return Shipping{}, app_error.InternalServerError(err)
	}

	if errs, err := s.validate(shipping); err != nil {
		return Shipping{}, err
	} else if len(errs) > 0 {
		return Shipping{}, app_error.ValidationError(errs)
	}

	return s.repo.Update(shipping)
}

func (s *usecase) Start(shipping Db) (Shipping, error) {
	if shipping.Status != STATUS_NEW {
		return Shipping{}, app_error.ValidationError(`Начать можно только перевозку со статусом "Новый"`)
	}

	now := time.Now()
	shipping.Status = STATUS_ACTIVE
	shipping.TimeStart = &now

	if errs, err := s.validate(shipping); err != nil {
		return Shipping{}, err
	} else if len(errs) > 0 {
		return Shipping{}, app_error.ValidationError(errs)
	}

	return s.repo.Update(shipping)
}

func (s *usecase) End(shipping Db) (Shipping, error) {
	if shipping.Status != STATUS_ACTIVE {
		return Shipping{}, app_error.ValidationError(`Завершить можно только перевозку со статусом "Активная"`)
	}

	now := time.Now()
	shipping.Status = STATUS_END
	shipping.TimeEnd = &now

	if errs, err := s.validate(shipping); err != nil {
		return Shipping{}, err
	} else if len(errs) > 0 {
		return Shipping{}, app_error.ValidationError(errs)
	}

	return s.repo.Update(shipping)
}

func (s *usecase) GetById(id int) (Shipping, error) {
	return s.repo.GetById(id)
}

func (s *usecase) GetActiveByModemImei(imei uint64) (Shipping, error) {
	modem, err := s.usecase.Modem.GetByImei(imei)
	if err != nil {
		return Shipping{}, err
	}

	return s.repo.GetActiveByModemId(modem.Id)
}

func (s *usecase) GetDbById(id int) (Db, error) {
	return s.repo.GetDbById(id)
}

func (s *usecase) List(queryParams QueryParams) (query.List[ShippingForList], error) {
	if errs := s.validator.Struct(queryParams); errs != nil {
		return query.List[ShippingForList]{}, app_error.ValidationError(errs)
	}

	return s.repo.List(queryParams)
}

func (s *usecase) Exists(id int) (bool, error) {
	return s.repo.Exists(id)
}

func (s *usecase) ExistsByUnique(customNumber, createDate string, id, number int) (bool, error) {
	return s.repo.ExistsByUnique(customNumber, createDate, id, number)
}

func (s *usecase) DeleteById(id int) (bool, error) {
	filesDirectory := s.GetFilesDirectory(id)
	if err := os.RemoveAll(filesDirectory); err != nil {
		return false, err
	}

	return s.repo.DeleteById(id)
}

func (s *usecase) Route(id int) (route.Route, error) {
	if shipping, err := s.repo.GetById(id); err != nil {
		return route.Route{}, err
	} else {
		return s.usecase.Route.GetById(shipping.Route.Id)
	}
}

func (s *usecase) AddFiles(id int, files []File) (Shipping, error) {
	shipping, err := s.GetDbById(id)

	if err != nil {
		return Shipping{}, app_error.ErrNotFound
	}

	shipping.Files = append(shipping.Files, files...)

	return s.repo.Update(shipping)
}

func (s *usecase) RemoveFilesFromDisk(id int, files []File) error {
	var res error

	for _, file := range files {
		filesDirectory := s.GetFilesDirectory(id)
		filePath := fmt.Sprintf("%s/%s", filesDirectory, file.Name)
		if err := os.Remove(filePath); err != nil {
			res = err
		}
	}

	return res
}

func (s *usecase) GetFilesDirectory(id int) string {
	return fmt.Sprintf("%s/%d", s.filesPath, id)
}

func (s *usecase) GetFileInfo(id int, name string) (File, error) {
	shipping, err := s.GetDbById(id)

	if err != nil {
		return File{}, err
	}

	for _, file := range shipping.Files {
		if file.Name == name {
			filesDirectory := s.GetFilesDirectory(id)
			filePath := fmt.Sprintf("%s/%s", filesDirectory, name)
			if _, err := os.Stat(filePath); err == nil {
				return file, nil
			} else {
				return File{}, app_error.ErrNotFound
			}
		}
	}

	return File{}, app_error.ErrNotFound
}

func (s *usecase) RemoveFile(id int, name string) error {
	shipping, err := s.GetDbById(id)

	if err != nil {
		return err
	}

	for i, file := range shipping.Files {
		if file.Name == name {
			shipping.Files = append(shipping.Files[:i], shipping.Files[i+1:]...)
			if _, err := s.repo.Update(shipping); err != nil {
				return err
			}

			if err := s.RemoveFilesFromDisk(id, []File{file}); err != nil {
				return err
			}

			return nil
		}
	}

	return app_error.ErrNotFound
}

func (s *usecase) SetModemById(shippingId int, modemId int) (bool, error) {
	shipping, err := s.GetDbById(shippingId)

	if err != nil {
		return false, err
	}

	if shipping.Status != STATUS_NEW {
		return false, app_error.ValidationError(`Привязать модем можно только к перевозке со статусом "Новый"`)
	}

	modem, err := s.usecase.Modem.GetById(modemId)

	if err != nil {
		return false, err
	}

	if imei, err := strconv.ParseUint(modem.Imei, 10, 64); err != nil {
		return false, err
	} else if _, err := s.GetActiveByModemImei(imei); err == nil {
		return false, app_error.ValidationError(`Модем привязан к активной поездке`)
	}

	shipping.Modem = &modem.Id

	if _, err := s.repo.Update(shipping); err != nil {
		return false, err
	}

	return true, nil
}

func (s *usecase) SetModemByImei(shippingId int, imei uint64) (bool, error) {
	shipping, err := s.GetDbById(shippingId)

	if err != nil {
		return false, err
	}

	if shipping.Status != STATUS_NEW {
		return false, app_error.ValidationError(`Привязать модем можно только к перевозке со статусом "Новый"`)
	}

	modem, err := s.usecase.Modem.GetByImei(imei)

	if err != nil {
		return false, err
	}

	if _, err := s.GetActiveByModemImei(imei); err != nil {
		return false, app_error.ValidationError(`Модем привязан к активной поездке`)
	}

	shipping.Modem = &modem.Id

	if _, err := s.repo.Update(shipping); err != nil {
		return false, err
	}

	return true, nil
}

func (s *usecase) Coordinates(params TrackQueryParams) ([]trackResponseCoordinate, error) {
	var response []trackResponseCoordinate

	shipping, err := s.GetDbById(params.Id)

	if err != nil {
		return response, err
	}

	if shipping.Modem == nil || shipping.TimeStart == nil {
		return response, app_error.ValidationError(`Поездка не начата`)
	}

	timeFrom := *shipping.TimeStart
	if !params.From.IsZero() && params.From.After(timeFrom) {
		timeFrom = params.From
	}

	queryParams := modem.TrackQueryParams{
		Id:        *shipping.Modem,
		From:      timeFrom,
		Limit:     params.Limit,
		OrderDesc: params.OrderDesc,
	}

	if shipping.TimeEnd != nil {
		queryParams.To = *shipping.TimeEnd
	}

	coord, err := s.usecase.Modem.Track(queryParams)
	if err != nil {
		return response, err
	}

	for _, row := range coord.Coordinates {
		response = append(response, trackResponseCoordinate{
			DevTime:            row.DevTime,
			Latitude:           row.Latitude,
			Longitude:          row.Longitude,
			MinDistanceToRoute: row.MinDistanceToRoute,
		})
	}

	return response, nil
}

func (s *usecase) Telemetry(params TrackQueryParams) ([]trackResponseTelemetry, error) {
	var response []trackResponseTelemetry

	shipping, err := s.GetDbById(params.Id)

	if err != nil {
		return response, err
	}

	if shipping.Modem == nil || shipping.TimeStart == nil {
		return response, app_error.ValidationError(`Поездка не начата`)
	}

	timeFrom := *shipping.TimeStart
	if !params.From.IsZero() && params.From.After(timeFrom) {
		timeFrom = params.From
	}

	queryParams := modemData.ListParams{
		ModemId:   *shipping.Modem,
		TimeFrom:  timeFrom,
		Limit:     params.Limit,
		OrderDesc: params.OrderDesc,
	}

	if shipping.TimeEnd != nil {
		queryParams.TimeTo = *shipping.TimeEnd
	}

	arch, err := s.usecase.ModemData.List(queryParams)
	if err != nil {
		return response, err
	}

	for _, row := range arch {
		var sealsData []trackResponseSealData
		for _, data := range row.SealsData {
			if data.Status == 0 && data.Errors == 0 {
				continue
			}

			sealsData = append(sealsData, trackResponseSealData{
				DevTime: data.DevTime,
				Seal: trackResponseSeal{
					Id:     data.Seal.Id,
					Serial: data.Seal.Serial,
				},
				Status:       data.Status,
				Errors:       data.Errors,
				BatteryLevel: data.BatteryLevel,
			})
		}

		if row.Status == 0 && row.ErrorsFlags == 0 && len(sealsData) == 0 {
			continue
		}

		response = append(response, trackResponseTelemetry{
			DevTime:     row.DevTime,
			Status:      row.Status,
			Latitude:    row.Latitude,
			Longitude:   row.Longitude,
			Altitude:    row.Altitude,
			SealsData:   sealsData,
			ErrorsFlags: row.ErrorsFlags,
		})
	}

	return response, nil
}

func (s *usecase) UpdateFileInfo(id int, name string, data UpdateFileRequest) (Shipping, error) {
	shipping, err := s.GetDbById(id)

	if err != nil {
		return Shipping{}, err
	}

	for i, file := range shipping.Files {
		if file.Name == name {
			if err := utils.BindFromStruct(data, &file); err != nil {
				return Shipping{}, app_error.InternalServerError(err)
			}

			shipping.Files[i] = file

			return s.repo.Update(shipping)
		}
	}

	return Shipping{}, app_error.ErrNotFound
}
