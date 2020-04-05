package errors

import (
	"errors"
	"net/http"
)

// Пакет для определения типичных ошибок, которые потом будут использоваться в кастомных обёртках

var (
	// общие
	ErrInternal     = errors.New("Internal error")
	ErrConflict     = errors.New("Conflict with exists data")
	ErrNoPermission = errors.New("No permission for current operation")

	// ошибки, связанные с юзером
	ErrUserNotFound  = errors.New("User not exist")
	ErrWrongPassword = errors.New("Wrong password")

	// ошибки, связанные с сессией
	ErrNoCookie        = errors.New("Not found cookie header")
	ErrSessionNotFound = errors.New("Cookie invalid, session not exist")

	// ошибки, связанные с досками
	ErrBoardNotFound = errors.New("Boards not found")

	// ошибки, связанные с колонками
	ErrColNotFound = errors.New("Column not found")

	// ошибки, связанные с тасками
	ErrTaskNotFound = errors.New("Task not found")

	// ошибки, связанные с бд
	ErrDbBadOperation  = errors.New("Unsuccessful ORM operation")
	ErrBadAvatarUpload = errors.New("Unsuccessful avatar upload")
)

var errorToCodeMap = map[error]int{
	// общие
	ErrInternal:     http.StatusInternalServerError,
	ErrConflict:     http.StatusConflict,
	ErrNoPermission: http.StatusForbidden,

	// ошибки, связанные с юзером
	ErrUserNotFound:  http.StatusNotFound,
	ErrWrongPassword: http.StatusPreconditionFailed,

	// ошибки, связанные с сессией
	ErrNoCookie:        http.StatusForbidden,
	ErrSessionNotFound: http.StatusNotFound,

	// ошибки, связанные с доской
	ErrBoardNotFound: http.StatusNotFound,

	// ошибки, связанные с колонкой
	ErrColNotFound: http.StatusNotFound,

	// ошибки, связанные с таской
	ErrTaskNotFound: http.StatusNotFound,
}

func ResolveErrorToCode(err error) (code int) {
	code, exist := errorToCodeMap[err]
	if exist != true {
		return http.StatusInternalServerError
	}
	return code
}
