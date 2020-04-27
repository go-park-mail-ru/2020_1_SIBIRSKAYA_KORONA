package errors

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/status"
)

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

	CommentNotFound = "comment not found"

	ChecklistNotFound = "checklist not found"

	ItemNotFound = "item not found"

	FileNotFound = "file not found"

	BadFileUploadS3 = "unsuccessful file upload to s3"
	BadFileDeleteS3 = "unsuccessful file delete on s3"

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

	// ошибки, связанные с комментариями
	ErrCommentNotFound = errors.New(CommentNotFound)

	// ошибки, связанные с чеклистами
	ErrChecklistNotFound = errors.New(ChecklistNotFound)

	// ошибки, связанные с итемами
	ErrItemNotFound = errors.New(ChecklistNotFound)

	// ошибки, связанные с файлами
	ErrFileNotFound    = errors.New(FileNotFound)
	ErrBadFileUploadS3 = errors.New(BadFileUploadS3)
	ErrBadFileDeleteS3 = errors.New(BadFileDeleteS3)

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

	CommentNotFound: ErrCommentNotFound,

	ChecklistNotFound: ErrChecklistNotFound,

	ItemNotFound: ErrItemNotFound,

	BadFileUploadS3: ErrBadFileUploadS3,
	BadFileDeleteS3: ErrBadFileDeleteS3,
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

	// ошибки, связанные с комментариями
	ErrCommentNotFound: http.StatusNotFound,

	// ошибки, связанные с чеклистами
	ErrChecklistNotFound: http.StatusNotFound,

	// ошибки, связанные с итемами
	ErrItemNotFound: http.StatusNotFound,

	// ошибки, связанные с файлами
	ErrFileNotFound:    http.StatusNotFound,
	ErrBadFileUploadS3: http.StatusUnprocessableEntity,
	ErrBadFileUploadS3: http.StatusUnprocessableEntity,
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
		return ErrInternal
	}
	return err
}
