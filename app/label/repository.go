package user

//go:generate mockgen -source=repository.go -package=mocks -destination=./mocks/user_repo_mock.go
type Repository interface {
}
