package app

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func (a *App) startSprintWeeklyReset(ctx context.Context) {
	c := cron.New(cron.WithLocation(time.UTC))

	// Запускаем обнуление каждый понедельник в 00:00 UTC
	_, err := c.AddFunc("0 0 * * 1", func() {
		a.resetSprintProgress(ctx)
	})

	if err != nil {
		log.Printf("error: failed to schedule sprint weekly reset: %v", err)
		return
	}

	c.Start()
	log.Println("sprint weekly reset scheduler started (every Monday at 00:00 UTC)")

	go func() {
		<-ctx.Done()
		log.Println("stopping sprint weekly reset scheduler...")
		c.Stop()
	}()
}

func (a *App) resetSprintProgress(ctx context.Context) {
	sprintService := a.serviceProvider.SprintService(ctx)
	if err := sprintService.ResetWeeklyProgress(ctx); err != nil {
		log.Printf("warning: failed to reset sprint progress: %v", err)
		return
	}
	log.Println("sprint progress reset completed")
}

