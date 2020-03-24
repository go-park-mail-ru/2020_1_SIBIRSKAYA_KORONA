package custom_errors

import (
	"fmt"
)

// type CustomDeliveryError {
// 	URL string
// 	Err error
// 	Code int
// }

// Может ли юзкейс решать какой статус отправлять ?
type CustomUsecaseError struct {
	Err  error
	Code int
}

// type CustomRepositoryError {}

func (customError *CustomUsecaseError) Error() string {
	return fmt.Sprintf("Usecase error: %v", customError.Err)
}
