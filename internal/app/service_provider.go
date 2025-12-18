package app

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/Dokhoyan/daily-routine/internal/config"
	"github.com/Dokhoyan/daily-routine/internal/logger"
	authHandler "github.com/Dokhoyan/daily-routine/internal/http-server/handlers/auth"
	habitHandler "github.com/Dokhoyan/daily-routine/internal/http-server/handlers/habit"
	sprintHandler "github.com/Dokhoyan/daily-routine/internal/http-server/handlers/sprint"
	settingsHandler "github.com/Dokhoyan/daily-routine/internal/http-server/handlers/settings"
	userHandler "github.com/Dokhoyan/daily-routine/internal/http-server/handlers/user"
	postgresRepo "github.com/Dokhoyan/daily-routine/internal/repository/postgres"
	"github.com/Dokhoyan/daily-routine/internal/service"
	authService "github.com/Dokhoyan/daily-routine/internal/service/auth"
	habitService "github.com/Dokhoyan/daily-routine/internal/service/habit"
	sprintService "github.com/Dokhoyan/daily-routine/internal/service/sprint"
	settingsService "github.com/Dokhoyan/daily-routine/internal/service/settings"
	userService "github.com/Dokhoyan/daily-routine/internal/service/user"
	_ "github.com/lib/pq"
)

type serviceProvider struct {
	pgConfig       config.PGConfig
	httpConfig     config.HTTPConfig
	telegramConfig config.TelegramConfig
	jwtConfig      config.JWTConfig
	corsConfig     config.CORSConfig
	testConfig     config.TestConfig
	authConfig     config.AuthConfig
	adminConfig    config.AdminConfig

	db   *sql.DB
	repo *postgresRepo.Repository

	authService     service.AuthService
	userService     service.UserService
	settingsService service.SettingsService
	habitService    service.HabitService
	sprintService   service.SprintService

	authImpl     *authHandler.Implementation
	userImpl     *userHandler.Implementation
	settingsImpl *settingsHandler.Implementation
	habitImpl    *habitHandler.Implementation
	sprintImpl   *sprintHandler.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			logger.Fatalf("failed to get pg config: %s", err.Error())
		}
		logger.Info("pg config loaded successfully")
		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			logger.Fatalf("failed to get http config: %s", err.Error())
		}
		logger.Infof("http config loaded: %s", cfg.Address())
		s.httpConfig = cfg
	}
	return s.httpConfig
}

func (s *serviceProvider) TelegramConfig() config.TelegramConfig {
	if s.telegramConfig == nil {
		cfg, err := config.NewTelegramConfig()
		if err != nil {
			logger.Fatalf("failed to get telegram config: %s", err.Error())
		}
		logger.Info("telegram config loaded")
		s.telegramConfig = cfg
	}
	return s.telegramConfig
}

func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			logger.Fatalf("failed to get jwt config: %s", err.Error())
		}
		logger.Info("jwt config loaded")
		s.jwtConfig = cfg
	}
	return s.jwtConfig
}

func (s *serviceProvider) CORSConfig() config.CORSConfig {
	if s.corsConfig == nil {
		s.corsConfig = config.NewCORSConfig()
	}
	return s.corsConfig
}

func (s *serviceProvider) TestConfig() config.TestConfig {
	if s.testConfig == nil {
		s.testConfig = config.NewTestConfig()
	}
	return s.testConfig
}

func (s *serviceProvider) AuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		s.authConfig = config.NewAuthConfig()
	}
	return s.authConfig
}

func (s *serviceProvider) AdminConfig() config.AdminConfig {
	if s.adminConfig == nil {
		cfg, err := config.NewAdminConfig()
		if err != nil {
			logger.Fatalf("failed to get admin config: %s", err.Error())
		}
		logger.Info("admin config loaded")
		s.adminConfig = cfg
	}
	return s.adminConfig
}

func (s *serviceProvider) DB(ctx context.Context) *sql.DB {
	if s.db == nil {
		pgCfg := s.PGConfig()
		logger.Info("connecting to database...")
		db, err := sql.Open("postgres", pgCfg.DSN())
		if err != nil {
			logger.Fatalf("failed to open database: %v", err)
		}

		if err := db.Ping(); err != nil {
			logger.Fatalf("failed to ping database: %v", err)
		}

		logger.Info("database connection established")
		s.db = db
	}
	return s.db
}

func (s *serviceProvider) Repository(ctx context.Context) *postgresRepo.Repository {
	if s.repo == nil {
		s.repo = postgresRepo.New(s.DB(ctx))
	}
	return s.repo
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		telegramCfg := s.TelegramConfig()
		jwtCfg := s.JWTConfig()
		authCfg := s.AuthConfig()
		repo := s.Repository(ctx)
		tokenCache := authService.NewMemoryTokenCache()
		s.authService = authService.NewService(telegramCfg.GetBotToken(), jwtCfg.GetSecret(), repo, repo, authCfg, tokenCache)
	}
	return s.authService
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		repo := s.Repository(ctx)
		s.userService = userService.NewService(repo)
	}
	return s.userService
}

func (s *serviceProvider) SettingsService(ctx context.Context) service.SettingsService {
	if s.settingsService == nil {
		repo := s.Repository(ctx)
		s.settingsService = settingsService.NewService(repo)
	}
	return s.settingsService
}

func (s *serviceProvider) HabitService(ctx context.Context) service.HabitService {
	if s.habitService == nil {
		repo := s.Repository(ctx)
		s.habitService = habitService.NewService(repo, repo, repo)
	}
	return s.habitService
}

func (s *serviceProvider) SprintService(ctx context.Context) service.SprintService {
	if s.sprintService == nil {
		repo := s.Repository(ctx)
		s.sprintService = sprintService.NewService(repo, repo, repo)
	}
	return s.sprintService
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *authHandler.Implementation {
	if s.authImpl == nil {
		s.authImpl = authHandler.NewImplementation(s.AuthService(ctx))
	}
	return s.authImpl
}

func (s *serviceProvider) UserImpl(ctx context.Context) *userHandler.Implementation {
	if s.userImpl == nil {
		s.userImpl = userHandler.NewImplementation(s.UserService(ctx))
	}
	return s.userImpl
}

func (s *serviceProvider) SettingsImpl(ctx context.Context) *settingsHandler.Implementation {
	if s.settingsImpl == nil {
		s.settingsImpl = settingsHandler.NewImplementation(s.SettingsService(ctx))
	}
	return s.settingsImpl
}

func (s *serviceProvider) HabitImpl(ctx context.Context) *habitHandler.Implementation {
	if s.habitImpl == nil {
		s.habitImpl = habitHandler.NewImplementation(s.HabitService(ctx))
	}
	return s.habitImpl
}

func (s *serviceProvider) SprintImpl(ctx context.Context) *sprintHandler.Implementation {
	if s.sprintImpl == nil {
		s.sprintImpl = sprintHandler.NewImplementation(s.SprintService(ctx), s.UserService(ctx))
	}
	return s.sprintImpl
}

func (s *serviceProvider) CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowedOrigin := s.CORSConfig().GetAllowedOrigin()
			testModeEnabled := s.TestConfig().IsTestModeEnabled()

			// Функция для установки CORS заголовков
			setCORSHeaders := func(allowedOriginValue string) {
				w.Header().Set("Access-Control-Allow-Origin", allowedOriginValue)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Vary", "Origin")
			}

			// Проверяем, разрешен ли origin
			shouldAllow := false
			originToAllow := ""

			// В тестовом режиме разрешаем два origin
			if testModeEnabled && origin != "" {
				if origin == "https://daily-routine.ru" || origin == "http://localhost:3000" {
					shouldAllow = true
					originToAllow = origin
				}
			}

			// Если не разрешен в тестовом режиме, проверяем обычную конфигурацию
			if !shouldAllow {
				if allowedOrigin == "*" {
					shouldAllow = true
					if origin != "" {
						originToAllow = origin
					} else {
						originToAllow = "*"
					}
				} else if allowedOrigin != "" && origin == allowedOrigin {
					shouldAllow = true
					originToAllow = origin
				} else if allowedOrigin == "" && origin != "" {
					// Для разработки: если ALLOWED_ORIGIN не установлен, разрешаем любой origin
					shouldAllow = true
					originToAllow = origin
				}
			}

			// Устанавливаем заголовки, если origin разрешен
			if shouldAllow && originToAllow != "" {
				setCORSHeaders(originToAllow)
			}

			// Обработка preflight запроса
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (s *serviceProvider) CloseDB() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
