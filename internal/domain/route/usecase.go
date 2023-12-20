package route

import (
	app_interface "seal/internal/app/interface"
	"seal/internal/domain/custom"
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"seal/pkg/app_error"
	"seal/pkg/utils"
)

type usecase struct {
	repo          Repo
	logger        app_interface.Logger
	customUsecase custom.Usecase
	validator     app_interface.Validator
}

func NewUsecase(repo Repo, logger app_interface.Logger, customUsecase custom.Usecase, validator app_interface.Validator) Usecase {
	return &usecase{repo, logger, customUsecase, validator}
}

func (s *usecase) Create(data CreateRequest) (Route, error) {
	var route Route

	if err := utils.BindFromStruct(data, &route); err != nil {
		return Route{}, app_error.InternalServerError(err)
	}

	if errs, err := validate[Route](s, route); err != nil {
		return Route{}, err
	} else if len(errs) > 0 {
		return Route{}, app_error.ValidationError(errs)
	}

	return s.repo.Create(route)
}

func (s *usecase) Update(id int, data UpdateRequest) (Route, error) {
	route, err := s.GetDbById(id)

	if err != nil {
		return Route{}, app_error.ErrNotFound
	}

	if err := utils.BindFromStruct(data, &route); err != nil {
		return Route{}, app_error.InternalServerError(err)
	}

	if errs, err := validate[Db](s, route); err != nil {
		return Route{}, err
	} else if len(errs) > 0 {
		return Route{}, app_error.ValidationError(errs)
	}

	return s.repo.Update(route)
}

func (s *usecase) GetById(id int) (Route, error) {
	return s.repo.GetById(id)
}

func (s *usecase) GetDbById(id int) (Db, error) {
	return s.repo.GetDbById(id)
}

func (s *usecase) List(queryParams transport.QueryParams) (query.List[Route], error) {
	if errs := s.validator.Struct(queryParams); errs != nil {
		return query.List[Route]{}, app_error.ValidationError(errs)
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
