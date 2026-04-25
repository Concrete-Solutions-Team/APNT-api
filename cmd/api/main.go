package main

import (
	"log"

	"github.com/slupx/smartest-backend/internal/auth"
	"github.com/slupx/smartest-backend/internal/config"
	"github.com/slupx/smartest-backend/internal/database"
	"github.com/slupx/smartest-backend/internal/server"
	"github.com/slupx/smartest-backend/internal/test"
)

func main() {
	cfg := config.LoadConfig()

	db := database.InitPostgres(cfg.DatabaseURL)

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	testRepo := test.NewRepository(db)
	testService := test.NewService(testRepo)
	testHandler := test.NewHandler(testService)

	// scan := scanner.NewScanner(ghCachedClient, notifier, subRepo, 10*time.Second)

	svr := server.NewServer(cfg.Port, authHandler, testHandler)
	svr.MountEndpoints()

	if err := svr.Start(); err != nil {
		log.Fatalf("Server start failed: %s", err.Error())
	}
}
