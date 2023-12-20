package seal

import (
	app_interface "seal/internal/app/interface"
	sealData "seal/internal/domain/seal_data"
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"seal/pkg/app_error"
	"seal/pkg/utils"
)

type CoreUseCase struct {
	SealData sealData.Usecase
}

type usecase struct {
	repo      Repo
	logger    app_interface.Logger
	validator app_interface.Validator
	usecase   CoreUseCase
}

func NewUsecase(repo Repo, logger app_interface.Logger, validator app_interface.Validator, coreUsecase CoreUseCase) Usecase {
	return &usecase{repo, logger, validator, coreUsecase}
}

func (s *usecase) GetById(id int) (Seal, error) {
	return s.repo.GetById(id)
}

func (s *usecase) GetDbById(id int) (Db, error) {
	return s.repo.GetDbById(id)
}

func (s *usecase) Update(id int, data UpdateRequest) (Seal, error) {
	seal, err := s.GetDbById(id)

	if err != nil {
		return Seal{}, err
	}

	if err := utils.BindFromStruct(data, &seal); err != nil {
		return Seal{}, app_error.InternalServerError(err)
	}

	return s.repo.Update(seal)
}

func (s *usecase) List(queryParams transport.QueryParams) (query.List[SealForList], error) {
	if errs := s.validator.Struct(queryParams); errs != nil {
		return query.List[SealForList]{}, app_error.ValidationError(errs)
	}

	return s.repo.List(queryParams)
}

func (s *usecase) Exists(id int) (bool, error) {
	return s.repo.Exists(id)
}

func (s *usecase) Archive(params ArchiveQueryParams) ([]ArchiveSealData, error) {
	if errs := s.validator.Struct(params); errs != nil {
		return nil, app_error.ValidationError(errs)
	}

	dataFromRepo, err := s.usecase.SealData.List(sealData.ListParams{
		SealId:    params.Id,
		TimeFrom:  params.From,
		TimeTo:    params.To,
		Limit:     params.Limit,
		OrderDesc: params.OrderDesc,
	})

	if err != nil {
		return []ArchiveSealData{}, err
	}

	archive := []ArchiveSealData{}

	for _, data := range dataFromRepo {
		archiveSealData := ArchiveSealData{
			DevTime: data.DevTime,
			Modem: ArchiveSealDataModem{
				Id:     data.Modem.Id,
				Serial: data.Modem.Serial,
			},
			Status:               data.Status,
			Errors:               data.Errors,
			SensitivityRange:     data.SensitivityRange,
			BatteryLevel:         data.BatteryLevel,
			Rssi:                 data.Rssi,
			Temperature:          data.Temperature,
			CountCommandsInQueue: data.CountCommandsInQueue,
		}

		archive = append(archive, archiveSealData)
	}

	return archive, nil
}
