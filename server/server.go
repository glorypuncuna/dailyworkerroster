package server

import (
	"database/sql"
	"log"
	"os"

	handler "dailyworkerroster/handlers"
	"dailyworkerroster/repository"
	"dailyworkerroster/service"

	_ "dailyworkerroster/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func NewServer() {
	// Example: Get DB connection string from environment variable
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(127.0.0.1:3306)/dailyworkerroster?parseTime=true"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	userRepo := &repository.UserRepository{DB: db}
	shiftRepo := &repository.ShiftRepository{DB: db}
	workerShiftRepo := &repository.WorkerShiftRepository{DB: db}

	userService := service.NewUserService(userRepo)
	shiftService := &service.ShiftService{
		ShiftRepo:       shiftRepo,
		WorkerShiftRepo: workerShiftRepo,
	}

	userHandler := handler.NewUserHandler(userService)
	shiftHandler := handler.NewShiftHandler(shiftService)

	router := gin.Default()

	SetupRoutes(router, shiftHandler, userHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running at http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
