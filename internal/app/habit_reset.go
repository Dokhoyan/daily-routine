package app

import (
	"context"
	"log"
	"sync"
	"time"

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
		log.Printf("error: failed to schedule habit daily reset: %v", err)
		return
	}

	c.Start()
	log.Println("habit daily reset scheduler started (every hour at :00)")

	go func() {
		<-ctx.Done()
		log.Println("stopping habit daily reset scheduler...")
		c.Stop()
	}()
}

func (a *App) processHabitDailyReset(ctx context.Context) {
	userService := a.serviceProvider.UserService(ctx)
	settingsService := a.serviceProvider.SettingsService(ctx)
	habitService := a.serviceProvider.HabitService(ctx)

	users, err := userService.GetAll(ctx)
	if err != nil {
		log.Printf("warning: failed to get users for daily reset: %v", err)
		return
	}

	now := time.Now()

	for _, user := range users {
		settings, err := settingsService.GetByUserID(ctx, user.ID)
		if err != nil {
			log.Printf("warning: failed to get settings for user %d: %v", user.ID, err)
			continue
		}

		loc, err := time.LoadLocation(settings.Timezone)
		if err != nil {
			log.Printf("warning: invalid timezone %s for user %d: %v", settings.Timezone, user.ID, err)
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
				log.Printf("warning: failed to get habits for user %d: %v", user.ID, err)
				continue
			}

			if err := habitService.ProcessDailyReset(ctx, user.ID, habits); err != nil {
				log.Printf("warning: failed to process daily reset for user %d: %v", user.ID, err)
				continue
			}

			// Проверяем выполнение спринтов после обновления привычек
			sprintService := a.serviceProvider.SprintService(ctx)
			if err := sprintService.CheckAndUpdateSprintProgress(ctx, user.ID); err != nil {
				log.Printf("warning: failed to check sprint progress for user %d: %v", user.ID, err)
			}

			lastProcessedDayMu.Lock()
			lastProcessedDay[user.ID] = dateKey
			lastProcessedDayMu.Unlock()

			log.Printf("processed daily reset for user %d (timezone: %s, date: %s)", user.ID, settings.Timezone, dateKey)
		}
	}
}
