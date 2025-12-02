package habit

import (
	"github.com/Dokhoyan/daily-routine/internal/repository"
	"github.com/Dokhoyan/daily-routine/internal/service"
)

type serv struct {
	habitRepo repository.HabitRepository
}

func NewService(habitRepo repository.HabitRepository) service.HabitService {
	return &serv{
		habitRepo: habitRepo,
	}
}
