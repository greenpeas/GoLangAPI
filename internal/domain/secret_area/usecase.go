package secret_area

import (
	app_interface "seal/internal/app/interface"
	"seal/internal/domain/user"
	"seal/internal/repository/pg/query"
	"seal/pkg/app_error"
	"seal/pkg/utils"
)

type CoreUseCase struct {
	User user.Usecase
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

func (s *usecase) Create(data CreateRequest, userId int) (SecretArea, error) {
	var db Db

	if err := utils.BindFromStruct(data, &db); err != nil {
		return SecretArea{}, app_error.InternalServerError(err)
	}

	db.Author = userId

	if errs, err := s.validate(db); err != nil {
		return SecretArea{}, err
	} else if len(errs) > 0 {
		return SecretArea{}, app_error.ValidationError(errs)
	}

	return s.repo.Create(db)
}

func (s *usecase) Update(id int, data UpdateRequest) (SecretArea, error) {
	secretArea, err := s.GetDbById(id)

	if err != nil {
		return SecretArea{}, app_error.ErrNotFound
	}

	if err := utils.BindFromStruct(data, &secretArea); err != nil {
		return SecretArea{}, app_error.InternalServerError(err)
	}

	if errs, err := s.validate(secretArea); err != nil {
		return SecretArea{}, err
	} else if len(errs) > 0 {
		return SecretArea{}, app_error.ValidationError(errs)
	}

	return s.repo.Update(secretArea)
}

func (s *usecase) GetById(id int) (SecretArea, error) {
	return s.repo.GetById(id)
}

func (s *usecase) GetDbById(id int) (Db, error) {
	return s.repo.GetDbById(id)
}

func (s *usecase) List(queryParams QueryParams) (query.List[SecretArea], error) {
	if errs := s.validator.Struct(queryParams); errs != nil {
		return query.List[SecretArea]{}, app_error.ValidationError(errs)
	}

	return s.repo.List(queryParams)
}

func (s *usecase) ExistsByUnique(id int, title string) (bool, error) {
	return s.repo.ExistsByUnique(id, title)
}

func (s *usecase) DeleteById(id int) (bool, error) {
	return s.repo.DeleteById(id)
}
