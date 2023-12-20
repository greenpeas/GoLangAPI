package user

import (
	app_interface "seal/internal/app/interface"
	"seal/internal/repository/pg/query"
	"seal/internal/transport"
	"seal/pkg/app_error"
	"seal/pkg/utils"
)

type usecase struct {
	repo           Repo
	passwordHasher passwordHasher
	logger         app_interface.Logger
	validator      app_interface.Validator
}

type passwordHasher interface {
	Hash(password string) (string, error)
}

func NewUsecase(repo Repo, hasher passwordHasher, logger app_interface.Logger, validator app_interface.Validator) Usecase {
	return &usecase{repo, hasher, logger, validator}
}

func (s *usecase) Create(data CreateRequest) (User, error) {
	var user Db

	if err := utils.BindFromStruct(data, &user); err != nil {
		return User{}, app_error.InternalServerError(err)
	}

	if hashedPassword, err := s.passwordHasher.Hash(user.Password); err != nil {
		return User{}, err
	} else {
		user.Password = hashedPassword
	}

	if errs, err := s.validate(user); err != nil {
		return User{}, err
	} else if len(errs) > 0 {
		return User{}, app_error.ValidationError(errs)
	}

	return s.repo.Create(user)
}

func (s *usecase) Update(id int, data UpdateRequest) (User, error) {
	user, err := s.GetDbById(id)

	if err != nil {
		return User{}, app_error.ErrNotFound
	}

	oldHashedPassword := user.Password

	if err := utils.BindFromStruct(data, &user); err != nil {
		return User{}, app_error.InternalServerError(err)
	}

	if user.Password != oldHashedPassword {
		if hashedPassword, err := s.passwordHasher.Hash(user.Password); err != nil {
			s.logger.Error(err.Error())
			return User{}, err
		} else {
			user.Password = hashedPassword
		}
	}

	if errs, err := s.validate(user); err != nil {
		return User{}, err
	} else if len(errs) > 0 {
		return User{}, app_error.ValidationError(errs)
	}

	return s.repo.Update(user)
}

func (s *usecase) GetByCredentials(login, password string) (User, error) {
	hashedPassword, err := s.passwordHasher.Hash(password)

	if err != nil {
		s.logger.Error(err.Error())
		return User{}, nil
	}

	return s.repo.GetByCredentials(login, hashedPassword)
}

func (s *usecase) GetById(id int) (User, error) {
	return s.repo.GetById(id)
}

func (s *usecase) GetDbById(id int) (Db, error) {
	return s.repo.GetDbById(id)
}

func (s *usecase) List(params transport.QueryParams) (query.List[User], error) {
	if errs := s.validator.Struct(params); errs != nil {
		return query.List[User]{}, app_error.ValidationError(errs)
	}

	return s.repo.List(params)
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
