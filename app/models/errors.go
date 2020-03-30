package models

import (
	"errors"
)

// Пакет для определения типичных ошибок, которые потом будут использоваться в кастомных обёртках

var (
	// общие
	ErrInternal = errors.New("Internal error")

	// ошибки, связанные с пользователем
	ErrUserNotExist    = errors.New("User not exist")
	ErrWrongPassword   = errors.New("Wrong password")
	ErrUserBadMarshall = errors.New("Invalid data for user update")

	// ошибки, связанные с доской
	ErrBoardNotExist = errors.New("board does not exist")

	// ошибки, связанные с колонками
	ErrColumnNotExist = errors.New("column does not exist")

	// ошибки, связанные с тасками
	ErrTaskNotExist = errors.New("task does not exist")


	// ошибки, связанные с бд
	ErrDbBadOperation  = errors.New("Unsuccessful ORM operation")
	ErrBadAvatarUpload = errors.New("Unsuccessful avatar upload")
)
