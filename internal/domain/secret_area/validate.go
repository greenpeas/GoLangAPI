package secret_area

import (
	"fmt"
	"seal/internal/domain"
	"sync"
)

func (s *usecase) validate(model Db) (map[string]string, error) {
	if errs := s.validator.Struct(model); errs != nil {
		s.logger.Debug("Ошибки валидации", errs)
		return errs, nil
	}

	var wg sync.WaitGroup
	wg.Add(2)
	errs := map[string]string{}
	resChan := make(chan domain.Res)

	go s.existsUser(&wg, resChan, model)
	go s.existsByUnique(&wg, resChan, model)

	go domain.CloseChannel(&wg, resChan)

	for res := range resChan {
		if res.Err != nil {
			s.logger.Error("Ошибки валидации", errs)
			return nil, res.Err
		}

		for k, v := range res.Errs {
			errs[k] = v
		}
	}

	if len(errs) > 0 {
		s.logger.Debug("Ошибки валидации", errs)
	}

	return errs, nil
}

func (s *usecase) existsUser(wg *sync.WaitGroup, ch chan domain.Res, secretArea Db) {
	defer wg.Done()
	if exists, err := s.usecase.User.Exists(secretArea.Author); err != nil {
		ch <- domain.Res{Errs: nil, Err: err}
	} else if !exists {
		ch <- domain.Res{Errs: map[string]string{"author": fmt.Sprintf("Пользователь %d не существует", secretArea.Author)}, Err: nil}
	}
}

func (s *usecase) existsByUnique(wg *sync.WaitGroup, ch chan domain.Res, secretArea Db) {
	defer wg.Done()
	if exists, err := s.ExistsByUnique(secretArea.Id, secretArea.Title); err != nil {
		ch <- domain.Res{Errs: nil, Err: err}
	} else if exists {
		errs := map[string]string{}
		errs["title"] = "Не уникально"
		ch <- domain.Res{Errs: errs, Err: nil}
	}
}
