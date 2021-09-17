package app

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/lotostudio/financial-api/internal/config"
	"github.com/lotostudio/financial-api/internal/handler"
	"github.com/lotostudio/financial-api/internal/repo"
	"github.com/lotostudio/financial-api/internal/server"
	"github.com/lotostudio/financial-api/internal/service"
	"github.com/lotostudio/financial-api/pkg/auth"
	"github.com/lotostudio/financial-api/pkg/database"
	"github.com/lotostudio/financial-api/pkg/hash"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Financial API
// @version 1.0.0
// @description API for personal financing

// @host localhost:8080
// @BasePath /api/v1/

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

// Run initializes application
func Run(configPath string) {
	// Load configs
	cfg := config.LoadConfig(configPath)

	// Database
	db, err := database.NewPostgres(cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.User, cfg.DB.Password, cfg.DB.SSLMode)

	if err != nil {
		log.Error(err)
		return
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})

	if err != nil {
		log.Error(err)
		return
	}

	// Migrations
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", cfg.DB.Name, driver)

	if err != nil {
		log.Error(err)
		return
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error(err)
		return
	}

	// Utils
	passwordHasher := hash.NewSHA1PasswordHasher(cfg.Auth.PasswordSalt)
	tokenManager, err := auth.NewJWTManager(cfg.Auth.JWT.Key, cfg.Auth.AccessTokenTTL, cfg.Auth.RefreshTokenLength)

	if err != nil {
		log.Error(err)
		return
	}

	// Init handlers
	repos := repo.NewRepos(db)
	services := service.NewServices(repos, passwordHasher, tokenManager, cfg.Auth.AccessTokenTTL, cfg.Auth.RefreshTokenTTL)
	handlers := handler.NewHandler(services, tokenManager)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	log.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err = srv.Stop(ctx); err != nil {
		log.Errorf("failed to stop server: %v", err)
	}

	if err = db.Close(); err != nil {
		log.Errorf("error occured on db connection close: %v", err)
	}
}
