package main

import (
	"cashback-serv/config"
	"cashback-serv/internal/handler"
	"cashback-serv/internal/repository"
	"cashback-serv/internal/service"
	"database/sql"
	"fmt"
	"log"

	_ "cashback-serv/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Cashback Service API
// @version 1.0
// @description CASHBACK SERVICE API
// @host localhost:8080
// @BasePath /
func main() {
	gin.SetMode(gin.ReleaseMode)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		log.Fatalf("Connection error with database: %v", err)
	}
	defer db.Close()

	cashbackRepo := repository.NewCashbackRepository(db)
	sourceRepo := repository.NewSourceRepository(db)
	sourceService := service.NewSourceService(sourceRepo)

	cashbackService := service.NewCashbackService(cashbackRepo, sourceService)

	cashbackHandler := handler.NewCashbackHandler(cashbackService)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	cashbackHandler.RegisterRoutes(router)

	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Something went wrogn: %v", err)
	}
}
