package app

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/config"
	"github.com/Dokhoyan/daily-routine/internal/http-server/middleware"
	"github.com/Dokhoyan/daily-routine/internal/logger"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		if a.httpServer != nil {
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer shutdownCancel()
			if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
				logger.Errorf("server shutdown error in defer: %v", err)
			}
		}
		if a.serviceProvider != nil {
			if err := a.serviceProvider.CloseDB(); err != nil {
				logger.Errorf("database close error: %v", err)
			}
		}
	}()

	serverErrChan := make(chan error, 1)

	go func() {
		err := a.runHTTPServer()
		if err != nil && err != http.ErrServerClosed {
			serverErrChan <- err
		}
	}()

	go a.startTokenCleanup(ctx)

	go a.startHabitDailyReset(ctx)

	go a.startSprintWeeklyReset(ctx)

	select {
	case <-ctx.Done():
		logger.Info("shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Errorf("server shutdown error: %v", err)
			return err
		}
		logger.Info("server stopped")
		return nil
	case err := <-serverErrChan:
		logger.Errorf("server error occurred, shutting down: %v", err)
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if shutdownErr := a.httpServer.Shutdown(shutdownCtx); shutdownErr != nil {
			logger.Errorf("server shutdown error after server error: %v", shutdownErr)
		}
		return err
	}
}

func (a *App) initDeps(ctx context.Context) error {
	flag.Parse()

	// Инициализируем логгер первым делом
	logger.InitDefault()
	logger.Info("initializing application dependencies")

	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := http.NewServeMux()

	authImpl := a.serviceProvider.AuthImpl(ctx)
	userImpl := a.serviceProvider.UserImpl(ctx)
	settingsImpl := a.serviceProvider.SettingsImpl(ctx)
	habitImpl := a.serviceProvider.HabitImpl(ctx)
	sprintImpl := a.serviceProvider.SprintImpl(ctx)

	corsMiddleware := a.serviceProvider.CORSMiddleware()
	authMiddleware := middleware.AuthMiddleware(a.serviceProvider.AuthService(ctx))
	adminMiddleware := middleware.AdminMiddleware(a.serviceProvider.AdminConfig())

	publicMux := http.NewServeMux()
	authImpl.RegisterRoutes(publicMux)
	publicMux.HandleFunc("/users", userImpl.GetAll)

	publicHandler := corsMiddleware(publicMux)
	mux.Handle("/login/", publicHandler)
	mux.Handle("/auth/getaccesstoken", publicHandler)
	mux.Handle("/auth/getrefreshtoken", publicHandler)
	mux.Handle("/auth/revoke", publicHandler)
	mux.Handle("/users", publicHandler)

	protectedAuthMux := http.NewServeMux()
	protectedAuthMux.HandleFunc("/auth/revokeall", authImpl.RevokeAllTokens)
	protectedAuthHandler := corsMiddleware(authMiddleware(protectedAuthMux))
	mux.Handle("/auth/revokeall", protectedAuthHandler)

	protectedMux := http.NewServeMux()
	settingsImpl.RegisterRoutes(protectedMux)
	habitImpl.RegisterRoutes(protectedMux)
	userImpl.RegisterRoutes(protectedMux)
	sprintImpl.RegisterRoutes(protectedMux)

	protectedHandler := corsMiddleware(authMiddleware(protectedMux))

	mux.Handle("/user/me", protectedHandler)
	mux.Handle("/user/me/", protectedHandler)
	mux.Handle("/habits", protectedHandler)
	mux.Handle("/habits/", protectedHandler)
	mux.Handle("/sprints", protectedHandler)
	mux.Handle("/sprints/", protectedHandler)

	// Админские роуты для управления спринтами (только Basic Auth, без JWT)
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/sprints", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sprintImpl.GetAll(w, r)
		} else if r.Method == http.MethodPost {
			sprintImpl.Create(w, r)
		}
	})
	adminMux.HandleFunc("/sprints/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut || r.Method == http.MethodPatch {
			sprintImpl.Update(w, r)
		} else if r.Method == http.MethodDelete {
			sprintImpl.Delete(w, r)
		} else if r.Method == http.MethodGet {
			sprintImpl.GetAll(w, r)
		}
	})
	adminHandler := http.StripPrefix("/admin", corsMiddleware(adminMiddleware(adminMux)))
	mux.Handle("/admin/sprints", adminHandler)
	mux.Handle("/admin/sprints/", adminHandler)

	testCfg := a.serviceProvider.TestConfig()
	if testCfg.IsTestModeEnabled() {
		testTokenHandler := corsMiddleware(http.HandlerFunc(authImpl.TestToken(testCfg.IsTestModeEnabled(), testCfg.GetTestUserID())))
		mux.Handle("/auth/test/token", testTokenHandler)
		logger.Warn("test mode enabled: /auth/test/token endpoint is available")
	}

	rootMux := http.NewServeMux()
	rootMux.Handle("/api/", http.StripPrefix("/api", mux))
	rootMux.Handle("/", mux)

	httpConfig := a.serviceProvider.HTTPConfig()
	a.httpServer = &http.Server{
		Addr:              httpConfig.Address(),
		Handler:           rootMux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return nil
}

func (a *App) runHTTPServer() error {
	addr := a.serviceProvider.HTTPConfig().Address()
	logger.Infof("HTTP server is running on http://%s", addr)
	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Errorf("HTTP server error: %v", err)
		return err
	}
	return nil
}

func (a *App) startTokenCleanup(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	a.cleanupExpiredTokens(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.cleanupExpiredTokens(ctx)
		}
	}
}

func (a *App) cleanupExpiredTokens(ctx context.Context) {
	repo := a.serviceProvider.Repository(ctx)
	if err := repo.DeleteExpiredTokens(ctx); err != nil {
		logger.Warnf("failed to delete expired tokens: %v", err)
	}
	if err := repo.DeleteExpiredBlacklistEntries(ctx); err != nil {
		logger.Warnf("failed to delete expired blacklist entries: %v", err)
	}
}
