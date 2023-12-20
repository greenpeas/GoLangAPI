package seal_model

import (
	app_interface "seal/internal/app/interface"
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"seal/pkg/app_error"
)

type usecase struct {
	repo      Repo
	logger    app_interface.Logger
	validator app_interface.Validator
}

func NewUsecase(repo Repo, logger app_interface.Logger, validator app_interface.Validator) Usecase {
	return &usecase{repo, logger, validator}
}

func (s *usecase) List(queryParams transport.QueryParams) (query.List[SealModel], error) {
	if errs := s.validator.Struct(queryParams); errs != nil {
		return query.List[SealModel]{}, app_error.ValidationError(errs)
	}

	return s.repo.List(queryParams)
}
