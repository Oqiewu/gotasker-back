package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Получаем порт из переменной окружения или используем 8080 по умолчанию
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Создаем Gin router с дефолтными middleware (Logger, Recovery)
	router := gin.Default()

	// Временный health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "GoTasker API is running",
			"version": "0.1.0",
		})
	})

	// Временный endpoint для проверки
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to GoTasker API! Visit /health for health check",
			"docs":    "API documentation will be available at /docs (coming soon)",
		})
	})

	// API группа для будущих endpoints
	api := router.Group("/api/v1")
	{
		// Временная заглушка для tasks
		api.GET("/tasks", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Tasks endpoint - coming soon",
				"tasks":   []interface{}{},
			})
		})
	}

	// Запускаем сервер
	log.Printf("🚀 Starting GoTasker server on port %s", port)
	log.Printf("📝 Health check: http://localhost:%s/health", port)
	log.Printf("📚 API base: http://localhost:%s/api/v1", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
