package cstmerr

import (
	"fmt"
)

// Главным образом используется для логов
type DeliveryError struct {
	URL  string
	Err  error
	Code int
}

type UseError struct {
	Err  error
	Code int
}

// уровень репозитория не заботится о статусах ответа
type RepoError struct {
	Err error
}

// Будем использовать эти методы  для логов, а внутренние ошибки - для отдачи пользователям
func (customError *UseError) Error() string {
	return fmt.Sprintf("Usecase error: %v", customError.Err)
}

func (customError *RepoError) Error() string {
	return fmt.Sprintf("Repository error: %v", customError.Err)
}
