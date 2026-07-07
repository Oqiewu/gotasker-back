package router

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gotasker/gotasker-back/src/internal/handlers"
	"github.com/gotasker/gotasker-back/src/internal/repository"
)

// New собирает Gin-роутер со всеми маршрутами приложения
func New(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "GoTasker API is running",
			"version": "0.1.0",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to GoTasker API! Visit /health for health check",
		})
	})

	taskHandler := handlers.NewTaskHandler(repository.NewTaskRepository(db))

	api := r.Group("/api/v1")
	{
		tasks := api.Group("/tasks")
		{
			tasks.POST("", taskHandler.Create)
			tasks.GET("", taskHandler.List)
			tasks.GET("/:id", taskHandler.GetByID)
			tasks.PATCH("/:id", taskHandler.Update)
			tasks.DELETE("/:id", taskHandler.Delete)
		}
	}

	return r
}
