package route

import (
	"seal/internal/domain"
	"sync"
)

type valid interface {
	getId() int
	getTitle() string
}

func validate[T valid](s *usecase, model T) (map[string]string, error) {
	if errs := s.validator.Struct(model); errs != nil {
		s.logger.Debug("Ошибки валидации", errs)
		return errs, nil
	}

	var wg sync.WaitGroup
	wg.Add(1)
	errs := map[string]string{}
	resChan := make(chan domain.Res)

	go existsByUnique(s, &wg, resChan, model)

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

func existsByUnique[T valid](s *usecase, wg *sync.WaitGroup, ch chan domain.Res, model T) {
	defer wg.Done()
	if exists, err := s.ExistsByUnique(model.getId(), model.getTitle()); err != nil {
		ch <- domain.Res{Errs: nil, Err: err}
	} else if exists {
		ch <- domain.Res{Errs: map[string]string{"title": "Не уникально"}, Err: nil}
	}
}
