package cstmerr

import (
	"fmt"
)

// Главным образом используется для логов
type CustomDeliveryError struct {
	URL  string
	Err  error
	Code int
}

type CustomUsecaseError struct {
	Err  error
	Code int
}

// уровень репозитория не заботится о статусах ответа
type CustomRepositoryError struct {
	Err error
}

// Будем использовать эти методы  для логов, а внутренние ошибки - для отдачи пользователям
func (customError *CustomUsecaseError) Error() string {
	return fmt.Sprintf("Usecase error: %v", customError.Err)
}

func (customError *CustomRepositoryError) Error() string {
	return fmt.Sprintf("Repository error: %v", customError.Err)
}
