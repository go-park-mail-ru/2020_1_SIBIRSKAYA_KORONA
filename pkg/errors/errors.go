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
	ErrUserBadNickname = errors.New("Nickname is already in use")

	// ошибки, связанные с сессией
	ErrNoCookie        = errors.New("Not found cookie header")
	ErrSessionNotExist = errors.New("Cookie invalid, session not exist")

	// ошибки, связанные с досками
	ErrBoardsNotFound = errors.New("Boards not found for this user")
	ErrNoPermission   = errors.New("No permission for current operation")

	// ошибки, связанные с бд
	ErrDbBadOperation  = errors.New("Unsuccessful ORM operation")
	ErrBadAvatarUpload = errors.New("Unsuccessful avatar upload")
)

var errorToCodeMap = map[error]int{
	ErrInternal: http.StatusInternalServerError,

	ErrUserNotExist:    http.StatusNotFound,
	ErrUserBadNickname: http.StatusConflict,
	ErrWrongPassword:   http.StatusUnauthorized,
	ErrUserBadMarshall: http.StatusBadRequest,

	ErrNoCookie:        http.StatusForbidden,
	ErrSessionNotExist: http.StatusNotFound,

	ErrBoardsNotFound: http.StatusNotFound,
	ErrNoPermission:   http.StatusUnauthorized,
}

func ResolveErrorToCode(err error) (code int) {
	code, exist := errorToCodeMap[err]
	if exist != true {
		return http.StatusInternalServerError
	}

	return code
}
