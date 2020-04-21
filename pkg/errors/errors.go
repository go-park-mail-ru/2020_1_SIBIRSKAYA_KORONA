package errors

import (
	"errors"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

// Пакет для определения типичных ошибок, которые потом будут использоваться в кастомных обёртках

const (
	Internal     = "internal error"
	Conflict     = "conflict with exists data"
	NoPermission = "no permission for current operation"

	UserNotFound  = "user not exist"
	WrongPassword = "wrong password"

	NoCookie        = "not found cookie header"
	SessionNotFound = "cookie invalid, session not exist"
	DetectedCSRF    = "CSRF is confirmed"

	BoardNotFound = "boards not found"

	ColNotFound = "column not found"

	LabelNotFound = "label not found"

	TaskNotFound = "task not found"

	DbBadOperation  = "unsuccessful ORM operation"
	BadAvatarUpload = "unsuccessful avatar upload"
)

var (
	// общие
	ErrInternal     = errors.New(Internal)
	ErrConflict     = errors.New(Conflict)
	ErrNoPermission = errors.New(NoPermission)

	// ошибки, связанные с юзером
	ErrUserNotFound  = errors.New(UserNotFound)
	ErrWrongPassword = errors.New(WrongPassword)

	// ошибки, связанные с сессией
	ErrNoCookie        = errors.New(NoCookie)
	ErrSessionNotFound = errors.New(SessionNotFound)
	ErrDetectedCSRF    = errors.New(DetectedCSRF) // В целях дебага

	// ошибки, связанные с досками
	ErrBoardNotFound = errors.New(BoardNotFound)

	// ошибки, связанные с колонками
	ErrColNotFound = errors.New(ColNotFound)

	// ошибки, связанные с лейблом
	ErrLabelNotFound = errors.New(LabelNotFound)

	// ошибки, связанные с тасками
	ErrTaskNotFound = errors.New(TaskNotFound)

	// ошибки, связанные с бд
	ErrDbBadOperation  = errors.New(DbBadOperation)
	ErrBadAvatarUpload = errors.New(BadAvatarUpload)
)

var messToError = map[string]error{
	Internal:     ErrInternal,
	Conflict:     ErrConflict,
	NoPermission: ErrNoPermission,

	UserNotFound:  ErrUserNotFound,
	WrongPassword: ErrWrongPassword,

	NoCookie:        ErrNoCookie,
	SessionNotFound: ErrSessionNotFound,
	DetectedCSRF:    ErrDetectedCSRF,

	BoardNotFound: ErrBoardNotFound,

	ColNotFound: ErrColNotFound,

	LabelNotFound: ErrLabelNotFound,

	TaskNotFound: ErrTaskNotFound,
}

var errorToCodeMap = map[error]int{
	// общие
	ErrInternal:     http.StatusInternalServerError,
	ErrConflict:     http.StatusConflict,
	ErrNoPermission: http.StatusForbidden,

	// ошибки, связанные с юзером
	ErrUserNotFound:  http.StatusNotFound,
	ErrWrongPassword: http.StatusPreconditionFailed,
	//TODO код csfr?

	// ошибки, связанные с сессией
	ErrNoCookie:        http.StatusForbidden,
	ErrSessionNotFound: http.StatusNotFound,

	// ошибки, связанные с доской
	ErrBoardNotFound: http.StatusNotFound,

	// ошибки, связанные с колонкой
	ErrColNotFound: http.StatusNotFound,

	// ошибки, связанные с лейблом
	ErrLabelNotFound: http.StatusNotFound,

	// ошибки, связанные с таской
	ErrTaskNotFound: http.StatusNotFound,
}

func ResolveErrorToCode(err error) (code int) {
	code, has := errorToCodeMap[err]
	if !has {
		return http.StatusInternalServerError
	}
	return code
}

func ResolveFromRPC(err error) error {
	err, has := messToError[status.Convert(err).Message()]
	if !has {
		log.Println(err)
		return ErrInternal
	}
	return err
}
