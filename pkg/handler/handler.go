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
		auth.POST("/verify", h.verifyCode)
	}

	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/goal")
		{
			lists.POST("/", h.createGoal)
			lists.GET("/", h.getAllGoals)
			lists.GET("/:id", h.getGoalById)
			lists.PUT("/:id", h.updateGoal)
			lists.DELETE("/:id", h.deleteGoal)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
				items.GET("/:item_id", h.getItemById)
				items.PUT("/:item_id", h.updateItem)
				items.DELETE("/:item_id", h.deleteItem)
			}
		}
	}
	return router
}
