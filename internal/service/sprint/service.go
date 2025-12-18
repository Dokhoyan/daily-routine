package sprint

import (
	"github.com/Dokhoyan/daily-routine/internal/repository"
)

type serv struct {
	sprintRepo repository.SprintRepository
	userRepo   repository.UserRepository
	habitRepo  repository.HabitRepository
}

func NewService(sprintRepo repository.SprintRepository, userRepo repository.UserRepository, habitRepo repository.HabitRepository) *serv {
	return &serv{
		sprintRepo: sprintRepo,
		userRepo:   userRepo,
		habitRepo:  habitRepo,
	}
}


