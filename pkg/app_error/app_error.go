package app_error

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	ErrInternalServer = newAppError("Ошибка сервера", http.StatusInternalServerError, nil)
	ErrNotFound       = newAppError("Запрашиваемый ресурс не найден", http.StatusNotFound, nil)
	ErrUnauthorized   = newAppError("Недостаточно прав", http.StatusUnauthorized, nil)
	ErrValidation     = newAppError("Ошибка в запросе", http.StatusUnprocessableEntity, nil)
	ErrBadRequest     = newAppError("Ошибка в запросе", http.StatusBadRequest, nil)
	ErrLoginPassword  = newAppError("Неверный логин или пароль", http.StatusForbidden, nil)
)

type AppError struct {
	Body     any   `json:"message,omitempty"`
	httpCode int   `json:"-"`
	Err      error `json:"err,omitempty"`
}

func InternalServerError(err error) *AppError {
	return newAppError(ErrInternalServer.Body, ErrInternalServer.httpCode, err)
}

func ValidationError(obj any) *AppError {
	return newAppError(obj, ErrValidation.httpCode, nil)
}

func BadRequestError(err error) *AppError {
	return newAppError(ErrBadRequest.Body, ErrBadRequest.httpCode, err)
}

func LoginPasswordError() *AppError {
	return newAppError(ErrLoginPassword.Body, ErrLoginPassword.httpCode, nil)
}

func (e AppError) Error() string {
	return fmt.Sprintf("%v", e.Body)
}

func (e AppError) Unwrap() error {
	return e.Err
}

func (e AppError) GetBody() any {
	return e.Body
}

func (e AppError) GetHttpCode() int {
	return e.httpCode
}

func (e AppError) Marshal() []byte {
	if marshal, err := json.Marshal(e); err != nil {
		return nil
	} else {
		return marshal
	}
}

func newAppError(body any, httpCode int, err error) *AppError {
	return &AppError{
		Body:     body,
		httpCode: httpCode,
		Err:      err,
	}
}
