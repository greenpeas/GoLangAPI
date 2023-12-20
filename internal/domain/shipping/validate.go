package shipping

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
	errs := map[string]string{}
	resChan := make(chan domain.Res)

	wg.Add(4)
	go s.existsUser(&wg, resChan, model)
	go s.existsRoute(&wg, resChan, model)
	go s.existsTransport(&wg, resChan, model)
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

func (s *usecase) existsUser(wg *sync.WaitGroup, ch chan domain.Res, shipping Db) {
	defer wg.Done()
	if exists, err := s.usecase.User.Exists(shipping.Author); err != nil {
		ch <- domain.Res{Errs: nil, Err: err}
	} else if !exists {
		ch <- domain.Res{Errs: map[string]string{"author": fmt.Sprintf("Пользователь %d не существует", shipping.Author)}, Err: nil}
	}
}

func (s *usecase) existsRoute(wg *sync.WaitGroup, ch chan domain.Res, shipping Db) {
	defer wg.Done()
	if exists, err := s.usecase.Route.Exists(shipping.Route); err != nil {
		ch <- domain.Res{Errs: nil, Err: err}
	} else if !exists {
		ch <- domain.Res{Errs: map[string]string{"route": fmt.Sprintf("Маршрут %d не существует", shipping.Route)}, Err: nil}
	}
}

func (s *usecase) existsTransport(wg *sync.WaitGroup, ch chan domain.Res, shipping Db) {
	defer wg.Done()
	if exists, err := s.usecase.Transport.Exists(shipping.Transport); err != nil {
		ch <- domain.Res{Errs: nil, Err: err}
	} else if !exists {
		ch <- domain.Res{Errs: map[string]string{"transport": fmt.Sprintf("Транспорт %d не существует", shipping.Transport)}, Err: nil}
	}
}

func (s *usecase) existsByUnique(wg *sync.WaitGroup, ch chan domain.Res, shipping Db) {
	defer wg.Done()
	if exists, err := s.ExistsByUnique(shipping.CustomNumber, shipping.CreateDate, shipping.Id, shipping.Number); err != nil {
		ch <- domain.Res{Errs: nil, Err: err}
	} else if exists {
		errs := map[string]string{}
		errs["custom_number"] = "Не уникально"
		errs["create_date"] = "Не уникально"
		errs["number"] = "Не уникально"
		ch <- domain.Res{Errs: errs, Err: nil}
	}
}
