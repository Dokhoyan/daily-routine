package app

import (
	"context"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/logger"
	"github.com/robfig/cron/v3"
)

func (a *App) startSprintWeeklyReset(ctx context.Context) {
	c := cron.New(cron.WithLocation(time.UTC))

	// Запускаем обнуление каждый понедельник в 00:00 UTC
	_, err := c.AddFunc("0 0 * * 1", func() {
		a.resetSprintProgress(ctx)
	})

	if err != nil {
		logger.Errorf("failed to schedule sprint weekly reset: %v", err)
		return
	}

	c.Start()
	logger.Info("sprint weekly reset scheduler started (every Monday at 00:00 UTC)")

	go func() {
		<-ctx.Done()
		logger.Info("stopping sprint weekly reset scheduler...")
		c.Stop()
	}()
}

func (a *App) resetSprintProgress(ctx context.Context) {
	sprintService := a.serviceProvider.SprintService(ctx)
	if err := sprintService.ResetWeeklyProgress(ctx); err != nil {
		logger.Warnf("failed to reset sprint progress: %v", err)
		return
	}
	logger.Info("sprint progress reset completed")
}


