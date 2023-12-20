package transport

import (
	app_interface "seal/internal/app/interface"
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"seal/pkg/app_error"
	"seal/pkg/utils"
)

type usecase struct {
	repo      Repo
	logger    app_interface.Logger
	validator app_interface.Validator
}

func NewUsecase(repo Repo, logger app_interface.Logger, validator app_interface.Validator) Usecase {
	return &usecase{repo, logger, validator}
}

func (s *usecase) Create(data CreateRequest, userId int) (Transport, error) {
	var transport Db

	if err := utils.BindFromStruct(data, &transport); err != nil {
		return Transport{}, app_error.InternalServerError(err)
	}

	transport.Author = userId

	if errs, err := s.validate(transport); err != nil {
		return Transport{}, err
	} else if len(errs) > 0 {
		return Transport{}, app_error.ValidationError(errs)
	}

	return s.repo.Create(transport)
}

func (s *usecase) Update(id int, data UpdateRequest) (Transport, error) {
	transport, err := s.GetDbById(id)

	if err != nil {
		return Transport{}, app_error.ErrNotFound
	}

	if err := utils.BindFromStruct(data, &transport); err != nil {
		return Transport{}, app_error.InternalServerError(err)
	}

	if errs, err := s.validate(transport); err != nil {
		return Transport{}, err
	} else if len(errs) > 0 {
		return Transport{}, app_error.ValidationError(errs)
	}

	return s.repo.Update(transport)
}

func (s *usecase) GetById(id int) (Transport, error) {
	return s.repo.GetById(id)
}

func (s *usecase) GetDbById(id int) (Db, error) {
	return s.repo.GetDbById(id)
}

func (s *usecase) List(queryParams transport.QueryParams) (query.List[Transport], error) {
	if errs := s.validator.Struct(queryParams); errs != nil {
		return query.List[Transport]{}, app_error.ValidationError(errs)
	}

	return s.repo.List(queryParams)
}

func (s *usecase) Exists(id int) (bool, error) {
	return s.repo.Exists(id)
}

func (s *usecase) ExistsByUnique(id int, title string) (bool, error) {
	return s.repo.ExistsByUnique(id, title)
}

func (s *usecase) DeleteById(id int) (bool, error) {
	return s.repo.DeleteById(id)
}
