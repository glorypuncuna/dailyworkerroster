package server

import (
	handler "dailyworkerroster/handlers"
	"dailyworkerroster/middleware"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, shiftHandler *handler.ShiftHandler, userHandler *handler.UserHandler) {
	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.POST("/signup", userHandler.SignUp)
	router.POST("/login", userHandler.Login)
	userGroup := router.Group("/")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.GET("/workers", userHandler.GetAllWorkers)
		userGroup.GET("/worker/:id", userHandler.GetWorkerByID)
		userGroup.GET("/worker/assigned", shiftHandler.GetAssignedShifts)
		userGroup.GET("/worker/available/:workerID", shiftHandler.GetAvailableShifts)
		userGroup.POST("/shift/:shiftID/request/:workerID", shiftHandler.RequestShift)
		userGroup.GET("/worker/requests/:workerID", shiftHandler.GetAllRequestedShifts)
	}

	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		adminGroup.POST("/shift", shiftHandler.CreateShift)
		adminGroup.PUT("/shift/:shiftID/approve/:workerID", shiftHandler.ApproveShiftRequest)
		adminGroup.PUT("/shift/:shiftID/reject/:workerID", shiftHandler.RejectShiftRequest)
		adminGroup.GET("/shifts/day", shiftHandler.GetShiftsByDay)
	}
}
