package app

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Dokhoyan/daily-routine/internal/config"
	"github.com/Dokhoyan/daily-routine/internal/http-server/middleware"
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
				log.Printf("server shutdown error in defer: %v", err)
			}
		}
		if a.serviceProvider != nil {
			if err := a.serviceProvider.CloseDB(); err != nil {
				log.Printf("database close error: %v", err)
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

	// Запускаем периодическую очистку истекших токенов
	go a.startTokenCleanup(ctx)

	// Запускаем ежедневный сброс привычек
	go a.startHabitDailyReset(ctx)

	select {
	case <-ctx.Done():
		log.Println("shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
			return err
		}
		log.Println("server stopped")
		return nil
	case err := <-serverErrChan:
		log.Printf("server error occurred, shutting down: %v", err)
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if shutdownErr := a.httpServer.Shutdown(shutdownCtx); shutdownErr != nil {
			log.Printf("server shutdown error after server error: %v", shutdownErr)
		}
		return err
	}
}

func (a *App) initDeps(ctx context.Context) error {
	flag.Parse()

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

	corsMiddleware := a.serviceProvider.CORSMiddleware()
	authMiddleware := middleware.AuthMiddleware(a.serviceProvider.AuthService(ctx))

	// Public routes
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

	protectedHandler := corsMiddleware(authMiddleware(protectedMux))

	usersHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "/settings") {
			settingsImpl.HandleSettingsRoutes(w, r)
			return
		}
		userImpl.HandleUserRoutes(w, r)
	})
	mux.Handle("/users/", corsMiddleware(authMiddleware(usersHandler)))
	mux.Handle("/habits", protectedHandler)
	mux.Handle("/habits/", protectedHandler)

	testCfg := a.serviceProvider.TestConfig()
	if testCfg.IsTestModeEnabled() {
		testTokenHandler := corsMiddleware(http.HandlerFunc(authImpl.TestToken(testCfg.IsTestModeEnabled(), testCfg.GetTestUserID())))
		mux.Handle("/auth/test/token", testTokenHandler)
		log.Println("⚠️  Test mode enabled: /auth/test/token endpoint is available")
	}

	httpConfig := a.serviceProvider.HTTPConfig()
	a.httpServer = &http.Server{
		Addr:              httpConfig.Address(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return nil
}

func (a *App) runHTTPServer() error {
	addr := a.serviceProvider.HTTPConfig().Address()
	log.Printf("HTTP server is running on http://%s", addr)
	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *App) startTokenCleanup(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour) // Очистка каждый час
	defer ticker.Stop()

	// Выполняем очистку сразу при старте
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
		log.Printf("warning: failed to delete expired tokens: %v", err)
	}
	if err := repo.DeleteExpiredBlacklistEntries(ctx); err != nil {
		log.Printf("warning: failed to delete expired blacklist entries: %v", err)
	}
}
