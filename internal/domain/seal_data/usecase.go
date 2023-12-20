package sealData

import (
	app_interface "seal/internal/app/interface"
)

type usecase struct {
	repo      Repo
	logger    app_interface.Logger
	validator app_interface.Validator
}

func NewUsecase(repo Repo, logger app_interface.Logger, validator app_interface.Validator) Usecase {
	return &usecase{repo, logger, validator}
}

func (s *usecase) List(params ListParams) ([]SealData, error) {
	return s.repo.List(params)
}
