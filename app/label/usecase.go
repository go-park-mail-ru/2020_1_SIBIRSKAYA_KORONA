package user

import (
)

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/user_usecase_mock.go
type UseCase interface {
}
