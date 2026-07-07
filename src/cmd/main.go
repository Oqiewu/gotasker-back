package main

import (
	"log"
	"os"

	"github.com/gotasker/gotasker-back/src/internal/database"
	"github.com/gotasker/gotasker-back/src/internal/router"
)

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Подключение к базе данных
	dbConfig := database.NewConfigFromEnv()
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Запуск миграций
	migrationsPath := "migrations"
	if _, err := os.Stat("/app/migrations"); err == nil {
		migrationsPath = "/app/migrations"
	}
	if err := database.RunMigrations(db, migrationsPath); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	r := router.New(db)

	// Запускаем сервер
	log.Printf("🚀 Starting GoTasker server on port %s", port)
	log.Printf("📝 Health check: http://localhost:%s/health", port)
	log.Printf("📚 API base: http://localhost:%s/api/v1", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
