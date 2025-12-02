package user

import (
	"github.com/Dokhoyan/daily-routine/internal/repository"
	"github.com/Dokhoyan/daily-routine/internal/service"
)

type serv struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) service.UserService {
	return &serv{
		userRepo: userRepo,
	}
}
