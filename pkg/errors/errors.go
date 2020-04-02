package errors

import (
	"errors"
	"net/http"
)

// Пакет для определения типичных ошибок, которые потом будут использоваться в кастомных обёртках

var (
	// общие
	ErrInternal = errors.New("Internal error")

	// ошибки, связанные с пользователем
	ErrUserNotExist    = errors.New("User not exist")
	ErrWrongPassword   = errors.New("Wrong password")
	ErrUserBadMarshall = errors.New("Invalid data for user update")

	// ошибки, связанные с сессией
	ErrSessionNotExist = errors.New("Session not exist")

	// ошибки, связанные с досками
	ErrBoardsNotFound = errors.New("Boards not found for this user")

	// ошибки, связанные с бд
	ErrDbBadOperation  = errors.New("Unsuccessful ORM operation")
	ErrBadAvatarUpload = errors.New("Unsuccessful avatar upload")
)

var errorToCodeMap = map[error]int{
	ErrInternal: http.StatusInternalServerError,

	ErrUserNotExist:    http.StatusNotFound,
	ErrWrongPassword:   http.StatusUnauthorized,
	ErrUserBadMarshall: http.StatusBadRequest,

	ErrSessionNotExist: http.StatusUnauthorized,
}

func ResolveErrorToCode(err error) (code int) {
	code, exist := errorToCodeMap[err]
	if exist != true {
		return http.StatusInternalServerError
	}

	return code
}
