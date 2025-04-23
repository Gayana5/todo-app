package handler

import (
	"github.com/Gayana5/todo-app/pkg/repository"
	"github.com/Gayana5/todo-app/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	repo     repository.Authorization
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/verify-code", h.verifyRegistrationCode)

		forgotPassword := router.Group("/forgot-password")
		{
			forgotPassword.POST("/send-code", h.forgotPassword)
			forgotPassword.POST("/verify-code", h.verifyResetCode)
			forgotPassword.PUT("/reset-password", h.resetPassword)
		}
	}

	api := router.Group("/api", h.userIdentity)
	{
		api.GET("/user", h.getInfo)
		api.PUT("/user", h.updateUserInfo)

		goals := api.Group("/goal")
		{
			goals.POST("/", h.createGoal)
			goals.GET("/", h.getAllGoals)
			goals.GET("/:goal_id", h.getGoalById)
			goals.PUT("/:goal_id", h.updateGoal)
			goals.DELETE("/:goal_id", h.deleteGoal)
			goals.GET("/:goal_id/askAI", h.askAI)

			tasks := goals.Group("/:goal_id/tasks")
			{
				tasks.POST("/", h.createTask)
				tasks.GET("/", h.getAllTasks)
				tasks.GET("/:task_id", h.getTaskById)
				tasks.PUT("/:task_id", h.updateTask)
				tasks.DELETE("/:task_id", h.deleteTask)
			}
		}
	}
	return router
}
