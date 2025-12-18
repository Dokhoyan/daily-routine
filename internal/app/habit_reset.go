package app

import (
	"context"
	"sync"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/logger"
	"github.com/robfig/cron/v3"
)

var lastProcessedDay = make(map[int64]string)
var lastProcessedDayMu sync.RWMutex

func (a *App) startHabitDailyReset(ctx context.Context) {
	c := cron.New(cron.WithLocation(time.UTC))

	_, err := c.AddFunc("0 * * * *", func() {
		a.processHabitDailyReset(ctx)
	})

	if err != nil {
		logger.Errorf("failed to schedule habit daily reset: %v", err)
		return
	}

	c.Start()
	logger.Info("habit daily reset scheduler started (every hour at :00)")

	go func() {
		<-ctx.Done()
		logger.Info("stopping habit daily reset scheduler...")
		c.Stop()
	}()
}

func (a *App) processHabitDailyReset(ctx context.Context) {
	userService := a.serviceProvider.UserService(ctx)
	settingsService := a.serviceProvider.SettingsService(ctx)
	habitService := a.serviceProvider.HabitService(ctx)

	users, err := userService.GetAll(ctx)
	if err != nil {
		logger.Warnf("failed to get users for daily reset: %v", err)
		return
	}

	now := time.Now()

	for _, user := range users {
		settings, err := settingsService.GetByUserID(ctx, user.ID)
		if err != nil {
			logger.Warnf("failed to get settings for user %d: %v", user.ID, err)
			continue
		}

		loc, err := time.LoadLocation(settings.Timezone)
		if err != nil {
			logger.Warnf("invalid timezone %s for user %d: %v", settings.Timezone, user.ID, err)
			continue
		}

		userTime := now.In(loc)

		if userTime.Hour() == 0 && userTime.Minute() == 0 {
			dateKey := userTime.Format("2006-01-02")

			lastProcessedDayMu.RLock()
			lastDate, alreadyProcessed := lastProcessedDay[user.ID]
			lastProcessedDayMu.RUnlock()

			if alreadyProcessed && lastDate == dateKey {
				continue
			}

			habits, err := habitService.GetByUserID(ctx, user.ID, nil, nil)
			if err != nil {
				logger.Warnf("failed to get habits for user %d: %v", user.ID, err)
				continue
			}

			if err := habitService.ProcessDailyReset(ctx, user.ID, habits); err != nil {
				logger.Warnf("failed to process daily reset for user %d: %v", user.ID, err)
				continue
			}

			// Проверяем выполнение спринтов после обновления привычек
			sprintService := a.serviceProvider.SprintService(ctx)
			if err := sprintService.CheckAndUpdateSprintProgress(ctx, user.ID); err != nil {
				logger.Warnf("failed to check sprint progress for user %d: %v", user.ID, err)
			}

			lastProcessedDayMu.Lock()
			lastProcessedDay[user.ID] = dateKey
			lastProcessedDayMu.Unlock()

			logger.Infof("processed daily reset for user %d (timezone: %s, date: %s)", user.ID, settings.Timezone, dateKey)
		}
	}
}
