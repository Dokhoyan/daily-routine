package habit

import (
	"github.com/Dokhoyan/daily-routine/internal/repository"
	"github.com/Dokhoyan/daily-routine/internal/service"
)

type serv struct {
	habitRepo  repository.HabitRepository
	sprintRepo repository.SprintRepository
	userRepo   repository.UserRepository
}

func NewService(habitRepo repository.HabitRepository, sprintRepo repository.SprintRepository, userRepo repository.UserRepository) service.HabitService {
	return &serv{
		habitRepo:  habitRepo,
		sprintRepo: sprintRepo,
		userRepo:   userRepo,
	}
}
