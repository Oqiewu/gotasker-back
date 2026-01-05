# GoTasker - Обучающий путь изучения Go

## Цель проекта

**GoTasker** - это комплексный учебный проект для изучения языка Go и современных практик backend-разработки. Вы создадите полноценный REST API для управления задачами (TODO list) с нуля, постепенно изучая все ключевые концепции Go и архитектурные паттерны.

## Что вы изучите

- Основы синтаксиса Go и идиомы языка
- Clean Architecture и разделение ответственности
- REST API разработка с Gin framework
- Работа с разными типами хранилищ данных (JSON → PostgreSQL)
- Dependency Injection и тестирование
- Middleware и обработка ошибок
- Контейнеризация с Docker
- Работа с модулями Go

---

## Этапы реализации

### Этап 0: Подготовка окружения ✅
**Цель**: Настроить базовую структуру проекта

**Задачи**:
- [x] Инициализировать Go модуль: `go mod init github.com/yourusername/gotasker`
- [x] Создать структуру директорий
- [x] Настроить Docker и docker-compose
- [x] Создать .gitignore и .env.example

**Что изучаете**:
- Система модулей Go (go.mod, go.sum)
- Организация проектов в Go
- Базовая настройка окружения

**Время**: 30 минут

---

### Этап 1: Модели данных (Domain Layer)
**Цель**: Создать типы данных для задач

**Файлы**: `models/task.go`

**Задачи**:
1. Создать структуру `Task` с полями:
   - ID (int)
   - Title (string)
   - Description (string)
   - Done (bool)
   - CreatedAt (time.Time)
   - UpdatedAt (time.Time)

2. Создать DTO структуры:
   - `CreateTaskRequest` (для создания задачи)
   - `UpdateTaskRequest` (для обновления задачи)

3. Добавить JSON tags для сериализации
4. Добавить валидационные tags (для будущего использования)

**Что изучаете**:
- Структуры (structs) в Go
- Теги (tags) для метаданных
- Работа с time.Time
- Концепция DTO (Data Transfer Objects)
- Указатели (*string, *bool) для optional полей

**Пример кода**:
```go
type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title" binding:"required"`
    Description string    `json:"description"`
    Done        bool      `json:"done"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
    Title       string `json:"title" binding:"required,min=1,max=200"`
    Description string `json:"description" binding:"max=1000"`
}
```

**Время**: 1 час

---

### Этап 2: Repository Interface (Repository Pattern)
**Цель**: Определить контракт для работы с данными

**Файлы**: `repository/task_repository.go`

**Задачи**:
1. Создать интерфейс `TaskRepository` с методами:
   - `Create(task *models.Task) error`
   - `GetAll() ([]models.Task, error)`
   - `GetByID(id int) (*models.Task, error)`
   - `Update(id int, task *models.Task) error`
   - `Delete(id int) error`

2. Определить кастомные ошибки:
   - `ErrTaskNotFound`
   - `ErrInvalidID`

**Что изучаете**:
- Интерфейсы в Go
- Repository Pattern
- Обработка ошибок (errors.New, fmt.Errorf)
- Указатели vs значения
- Соглашения об именовании

**Пример кода**:
```go
type TaskRepository interface {
    Create(task *models.Task) error
    GetAll() ([]models.Task, error)
    GetByID(id int) (*models.Task, error)
    Update(id int, task *models.Task) error
    Delete(id int) error
}

var (
    ErrTaskNotFound = errors.New("task not found")
    ErrInvalidID    = errors.New("invalid task ID")
)
```

**Время**: 30 минут

---

### Этап 3: JSON Repository Implementation
**Цель**: Реализовать хранилище в JSON файле

**Файлы**: `repository/json_task_repository.go`

**Задачи**:
1. Создать структуру `jsonTaskRepository` с полями:
   - `tasks []models.Task`
   - `filePath string`
   - `mu sync.RWMutex` (для thread-safety)
   - `nextID int`

2. Реализовать все методы интерфейса `TaskRepository`
3. Добавить методы `loadFromFile()` и `saveToFile()`
4. Реализовать конструктор `NewJSONTaskRepository(filePath string)`

**Что изучаете**:
- Слайсы (slices) и работа с ними
- Чтение/запись файлов (io/ioutil, os)
- JSON encoding/decoding
- Мьютексы и конкурентность (sync.RWMutex)
- Методы на структурах
- Defer для закрытия ресурсов
- Обработка ошибок

**Ключевые концепции**:
```go
type jsonTaskRepository struct {
    tasks    []models.Task
    filePath string
    mu       sync.RWMutex // RWMutex для параллельного чтения
    nextID   int
}

func (r *jsonTaskRepository) Create(task *models.Task) error {
    r.mu.Lock()         // Блокируем для записи
    defer r.mu.Unlock() // Автоматически разблокируем

    task.ID = r.nextID
    task.CreatedAt = time.Now()
    task.UpdatedAt = time.Now()
    r.nextID++

    r.tasks = append(r.tasks, *task)
    return r.saveToFile()
}
```

**Время**: 2-3 часа

---

### Этап 4: Service Layer (Business Logic)
**Цель**: Создать слой бизнес-логики

**Файлы**: `service/task_service.go`

**Задачи**:
1. Создать интерфейс `TaskService` с методами:
   - `CreateTask(req *models.CreateTaskRequest) (*models.Task, error)`
   - `GetAllTasks() ([]models.Task, error)`
   - `GetTaskByID(id int) (*models.Task, error)`
   - `UpdateTask(id int, req *models.UpdateTaskRequest) (*models.Task, error)`
   - `DeleteTask(id int) error`
   - `ToggleTaskStatus(id int) (*models.Task, error)`

2. Создать структуру `taskService` с зависимостью от `TaskRepository`
3. Реализовать все методы с валидацией и бизнес-логикой
4. Добавить конструктор с Dependency Injection

**Что изучаете**:
- Разделение concerns (handlers → service → repository)
- Dependency Injection
- Валидация данных
- Преобразование DTO в domain модели
- Обработка и проброс ошибок

**Пример кода**:
```go
type TaskService interface {
    CreateTask(req *models.CreateTaskRequest) (*models.Task, error)
    GetAllTasks() ([]models.Task, error)
    // ...
}

type taskService struct {
    repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
    return &taskService{repo: repo}
}

func (s *taskService) CreateTask(req *models.CreateTaskRequest) (*models.Task, error) {
    // Валидация
    if req.Title == "" {
        return nil, errors.New("title is required")
    }

    // Создание модели
    task := &models.Task{
        Title:       req.Title,
        Description: req.Description,
        Done:        false,
    }

    // Сохранение через repository
    if err := s.repo.Create(task); err != nil {
        return nil, fmt.Errorf("failed to create task: %w", err)
    }

    return task, nil
}
```

**Время**: 2 часа

---

### Этап 5: HTTP Handlers (Presentation Layer)
**Цель**: Создать HTTP endpoints с Gin framework

**Файлы**: `handlers/task_handler.go`

**Задачи**:
1. Добавить Gin: `go get -u github.com/gin-gonic/gin`
2. Создать структуру `TaskHandler` с зависимостью от `TaskService`
3. Реализовать handlers:
   - `CreateTask(c *gin.Context)`
   - `GetAllTasks(c *gin.Context)`
   - `GetTaskByID(c *gin.Context)`
   - `UpdateTask(c *gin.Context)`
   - `DeleteTask(c *gin.Context)`
   - `ToggleTaskStatus(c *gin.Context)`

4. Обработать HTTP коды ответов (200, 201, 400, 404, 500)
5. Добавить парсинг и валидацию запросов

**Что изучаете**:
- Gin web framework
- HTTP методы и коды состояния
- Парсинг JSON из запросов
- Валидация с gin binding
- Обработка path параметров
- Формирование JSON ответов

**Пример кода**:
```go
type TaskHandler struct {
    service service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
    return &TaskHandler{service: service}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
    var req models.CreateTaskRequest

    // Парсинг и валидация
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Вызов service
    task, err := h.service.CreateTask(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Успешный ответ
    c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
    // Получение параметра из URL
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
        return
    }

    task, err := h.service.GetTaskByID(id)
    if err != nil {
        if errors.Is(err, repository.ErrTaskNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, task)
}
```

**Время**: 2-3 часа

---

### Этап 6: Main Application Setup
**Цель**: Собрать все компоненты вместе

**Файлы**: `src/cmd/main.go`

**Задачи**:
1. Создать функцию `main()`
2. Инициализировать все слои приложения (DI цепочка):
   - Repository
   - Service
   - Handlers
3. Настроить Gin router с endpoints
4. Добавить базовый middleware (Logger, Recovery)
5. Запустить HTTP сервер

**Что изучаете**:
- Точка входа Go приложения
- Инициализация зависимостей
- Настройка роутинга
- Middleware в Gin
- Graceful shutdown (опционально)

**Пример кода**:
```go
package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "your-module/handlers"
    "your-module/repository"
    "your-module/service"
)

func main() {
    // Инициализация слоев (Dependency Injection)
    taskRepo := repository.NewJSONTaskRepository("data/tasks.json")
    taskService := service.NewTaskService(taskRepo)
    taskHandler := handlers.NewTaskHandler(taskService)

    // Настройка Gin router
    router := gin.Default() // Default включает Logger и Recovery middleware

    // Регистрация routes
    api := router.Group("/api/v1")
    {
        tasks := api.Group("/tasks")
        {
            tasks.POST("", taskHandler.CreateTask)
            tasks.GET("", taskHandler.GetAllTasks)
            tasks.GET("/:id", taskHandler.GetTaskByID)
            tasks.PUT("/:id", taskHandler.UpdateTask)
            tasks.DELETE("/:id", taskHandler.DeleteTask)
            tasks.PATCH("/:id/toggle", taskHandler.ToggleTaskStatus)
        }
    }

    // Запуск сервера
    log.Println("Server starting on :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

**Время**: 1 час

---

### Этап 7: Middleware и улучшения
**Цель**: Добавить полезный middleware

**Файлы**: `middleware/cors.go`, `middleware/logger.go`

**Задачи**:
1. Создать CORS middleware для разрешения cross-origin запросов
2. Настроить кастомный Logger middleware
3. Добавить Request ID middleware
4. Создать Error Handler middleware

**Что изучаете**:
- Middleware pattern
- HTTP headers
- CORS (Cross-Origin Resource Sharing)
- Логирование запросов
- gin.HandlerFunc

**Пример CORS middleware**:
```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

**Время**: 1-2 часа

---

### Этап 8: Конфигурация и Environment Variables
**Цель**: Сделать приложение конфигурируемым

**Файлы**: `config/config.go`

**Задачи**:
1. Установить viper: `go get github.com/spf13/viper`
2. Создать структуру `Config` для настроек
3. Реализовать загрузку из .env файла
4. Параметризовать порт, путь к данным, режим работы

**Что изучаете**:
- Environment variables
- Конфигурация приложений
- Библиотека viper
- Best practices для sensitive данных

**Пример**:
```go
type Config struct {
    ServerPort   string
    DataFilePath string
    Environment  string
}

func LoadConfig() (*Config, error) {
    viper.SetConfigFile(".env")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    return &Config{
        ServerPort:   viper.GetString("SERVER_PORT"),
        DataFilePath: viper.GetString("DATA_FILE_PATH"),
        Environment:  viper.GetString("ENVIRONMENT"),
    }, nil
}
```

**Время**: 1 час

---

### Этап 9: Unit Testing
**Цель**: Написать тесты для сервисов

**Файлы**: `service/task_service_test.go`

**Задачи**:
1. Создать mock repository
2. Написать тесты для `TaskService`:
   - TestCreateTask
   - TestGetAllTasks
   - TestGetTaskByID_NotFound
   - TestUpdateTask
   - TestDeleteTask
3. Использовать table-driven tests
4. Изучить testify/mock

**Что изучаете**:
- Тестирование в Go (testing package)
- Mocking зависимостей
- Table-driven tests
- testify/assert и testify/mock
- Команда `go test`

**Пример**:
```go
type mockTaskRepository struct {
    mock.Mock
}

func (m *mockTaskRepository) Create(task *models.Task) error {
    args := m.Called(task)
    return args.Error(0)
}

func TestCreateTask_Success(t *testing.T) {
    // Arrange
    mockRepo := new(mockTaskRepository)
    service := NewTaskService(mockRepo)

    req := &models.CreateTaskRequest{
        Title:       "Test Task",
        Description: "Test Description",
    }

    mockRepo.On("Create", mock.AnythingOfType("*models.Task")).Return(nil)

    // Act
    task, err := service.CreateTask(req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, task)
    assert.Equal(t, "Test Task", task.Title)
    mockRepo.AssertExpectations(t)
}
```

**Время**: 2-3 часа

---

### Этап 10: PostgreSQL Integration
**Цель**: Заменить JSON хранилище на PostgreSQL

**Файлы**: `repository/postgres_task_repository.go`

**Задачи**:
1. Добавить PostgreSQL драйвер: `go get github.com/lib/pq`
2. Добавить GORM ORM (опционально): `go get -u gorm.io/gorm`
3. Создать новую реализацию `TaskRepository` для PostgreSQL
4. Добавить миграции для создания таблиц
5. Настроить connection pool
6. Обновить docker-compose.yml с PostgreSQL сервисом

**Что изучаете**:
- Работа с SQL базами данных
- database/sql package
- Connection pooling
- Prepared statements
- Транзакции
- GORM ORM (опционально)
- Docker compose с несколькими сервисами

**Пример без ORM**:
```go
type postgresTaskRepository struct {
    db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) TaskRepository {
    return &postgresTaskRepository{db: db}
}

func (r *postgresTaskRepository) Create(task *models.Task) error {
    query := `
        INSERT INTO tasks (title, description, done, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

    err := r.db.QueryRow(
        query,
        task.Title,
        task.Description,
        task.Done,
        time.Now(),
        time.Now(),
    ).Scan(&task.ID)

    return err
}
```

**docker-compose.yml**:
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: gotasker
      DB_PASSWORD: password
      DB_NAME: gotasker_db

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: gotasker
      POSTGRES_PASSWORD: password
      POSTGRES_DB: gotasker_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

**Время**: 3-4 часа

---

### Этап 11: Миграции базы данных
**Цель**: Автоматизировать управление схемой БД

**Файлы**: `migrations/`, использование golang-migrate

**Задачи**:
1. Установить golang-migrate: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
2. Создать миграции для таблицы tasks:
   - `000001_create_tasks_table.up.sql`
   - `000001_create_tasks_table.down.sql`
3. Интегрировать миграции в приложение
4. Добавить команды в Makefile

**Что изучаете**:
- Database migrations
- Версионирование схемы БД
- golang-migrate
- SQL DDL команды

**Пример миграции**:
```sql
-- 000001_create_tasks_table.up.sql
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    done BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_done ON tasks(done);
```

**Время**: 1-2 часа

---

### Этап 12: Integration Tests
**Цель**: Тестировать весь API end-to-end

**Файлы**: `tests/integration/api_test.go`

**Задачи**:
1. Настроить тестовую БД (PostgreSQL или SQLite)
2. Написать integration тесты для всех endpoints
3. Использовать httptest для тестирования handlers
4. Добавить setup/teardown для тестовых данных

**Что изучаете**:
- Integration testing
- httptest package
- Тестирование с реальной БД
- Test fixtures

**Пример**:
```go
func TestCreateTaskAPI(t *testing.T) {
    // Setup
    router := setupTestRouter()

    // Prepare request
    body := `{"title":"Test Task","description":"Test Desc"}`
    req, _ := http.NewRequest("POST", "/api/v1/tasks", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    // Execute
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)

    var task models.Task
    err := json.Unmarshal(w.Body.Bytes(), &task)
    assert.NoError(t, err)
    assert.Equal(t, "Test Task", task.Title)
}
```

**Время**: 2-3 часа

---

### Этап 13: Документация API (Swagger)
**Цель**: Создать интерактивную документацию API

**Файлы**: Аннотации в handlers

**Задачи**:
1. Установить swaggo: `go get -u github.com/swaggo/swag/cmd/swag`
2. Установить gin-swagger: `go get -u github.com/swaggo/gin-swagger`
3. Добавить Swagger аннотации к handlers
4. Сгенерировать документацию: `swag init`
5. Подключить Swagger UI к приложению

**Что изучаете**:
- API документация
- OpenAPI/Swagger спецификация
- Аннотации в Go
- Автогенерация документации

**Пример**:
```go
// CreateTask godoc
// @Summary      Create a new task
// @Description  Create a new task with title and description
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task body models.CreateTaskRequest true "Task to create"
// @Success      201 {object} models.Task
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
    // ...
}
```

**Время**: 1-2 часа

---

### Этап 14: Логирование (Structured Logging)
**Цель**: Добавить качественное логирование

**Файлы**: `pkg/logger/logger.go`

**Задачи**:
1. Установить zap или logrus: `go get -u go.uber.org/zap`
2. Создать wrapper для logger
3. Добавить структурированное логирование во всех слоях
4. Настроить уровни логирования (Debug, Info, Warn, Error)
5. Логировать в файл и stdout

**Что изучаете**:
- Structured logging
- zap/logrus библиотеки
- Уровни логирования
- Контекстное логирование

**Пример**:
```go
import "go.uber.org/zap"

func (s *taskService) CreateTask(req *models.CreateTaskRequest) (*models.Task, error) {
    logger.Info("Creating new task",
        zap.String("title", req.Title),
        zap.String("description", req.Description),
    )

    task, err := s.repo.Create(req)
    if err != nil {
        logger.Error("Failed to create task",
            zap.Error(err),
            zap.String("title", req.Title),
        )
        return nil, err
    }

    logger.Info("Task created successfully", zap.Int("task_id", task.ID))
    return task, nil
}
```

**Время**: 1-2 часа

---

### Этап 15: Graceful Shutdown
**Цель**: Корректное завершение приложения

**Файлы**: `src/cmd/main.go`

**Задачи**:
1. Обработать системные сигналы (SIGINT, SIGTERM)
2. Реализовать graceful shutdown для HTTP сервера
3. Закрыть соединения с БД перед выходом
4. Использовать context для управления временем ожидания

**Что изучаете**:
- Signals в Unix/Linux
- context.Context
- Graceful shutdown pattern
- Управление ресурсами

**Пример**:
```go
func main() {
    // ... инициализация ...

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    // Запуск сервера в горутине
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // Ожидание сигнала завершения
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // Graceful shutdown с таймаутом
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exited")
}
```

**Время**: 1 час

---

### Этап 16 (Опционально): Аутентификация JWT
**Цель**: Добавить защиту API с JWT токенами

**Файлы**: `middleware/auth.go`, `service/auth_service.go`

**Задачи**:
1. Установить JWT библиотеку: `go get -u github.com/golang-jwt/jwt/v5`
2. Создать модели User (Login, Register)
3. Реализовать генерацию и валидацию JWT токенов
4. Создать Auth middleware
5. Защитить endpoints (кроме login/register)

**Что изучаете**:
- JSON Web Tokens (JWT)
- Аутентификация и авторизация
- Хеширование паролей (bcrypt)
- Protected routes

**Время**: 3-4 часа

---

### Этап 17 (Опционально): Rate Limiting
**Цель**: Защита от abuse

**Файлы**: `middleware/rate_limit.go`

**Задачи**:
1. Использовать библиотеку rate limiter
2. Создать middleware для ограничения запросов
3. Настроить лимиты (например, 100 req/min)

**Что изучаете**:
- Rate limiting
- Token bucket algorithm
- Защита API

**Время**: 1-2 часа

---

### Этап 18 (Опционально): CI/CD
**Цель**: Автоматизация тестирования и деплоя

**Файлы**: `.github/workflows/ci.yml`

**Задачи**:
1. Создать GitHub Actions workflow
2. Настроить автоматический запуск тестов
3. Добавить линтеры (golangci-lint)
4. Настроить сборку Docker образа

**Что изучаете**:
- CI/CD концепции
- GitHub Actions
- Автоматизация тестирования
- Линтинг кода

**Время**: 2-3 часа

---

## Дополнительные улучшения (для продвинутых)

### Этап 19: Pагинация и фильтрация
- Query параметры для пагинации (?page=1&limit=10)
- Фильтрация задач (?done=true)
- Сортировка (?sort=created_at&order=desc)

### Этап 20: Кеширование с Redis
- Установка Redis
- Кеширование часто запрашиваемых данных
- Cache invalidation стратегии

### Этап 21: WebSockets для real-time обновлений
- Добавить WebSocket endpoint
- Push уведомления о новых задачах
- Синхронизация между клиентами

### Этап 22: Микросервисная архитектура
- Разделить на отдельные сервисы (Task Service, User Service)
- gRPC для межсервисного взаимодействия
- Service mesh (опционально)

---

## Итоговая структура проекта

```
gotasker/
├── .github/
│   └── workflows/
│       └── ci.yml
├── src/
│   ├── cmd/
│   │   └── main.go                  # Точка входа
│   └── internal/
│       ├── config/
│       │   └── config.go            # Конфигурация
│       ├── models/
│       │   ├── task.go              # Domain модели
│       │   └── user.go              # (опционально)
│       ├── repository/
│       │   ├── task_repository.go   # Интерфейс
│       │   ├── json_task_repo.go    # JSON реализация
│       │   └── postgres_task_repo.go # PostgreSQL реализация
│       ├── service/
│       │   ├── task_service.go      # Бизнес-логика
│       │   └── task_service_test.go # Unit тесты
│       ├── handlers/
│       │   └── task_handler.go      # HTTP handlers
│       ├── middleware/
│       │   ├── cors.go              # CORS
│       │   ├── logger.go            # Логирование
│       │   ├── auth.go              # (опционально)
│       │   └── rate_limit.go        # (опционально)
│       └── pkg/
│           └── logger/
│               └── logger.go        # Logger wrapper
├── migrations/
│   ├── 000001_create_tasks_table.up.sql
│   └── 000001_create_tasks_table.down.sql
├── tests/
│   └── integration/
│       └── api_test.go
├── data/
│   └── tasks.json                   # JSON хранилище
├── docs/                            # Swagger документация (автогенерация)
├── docker-compose.yml
├── Dockerfile
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile                         # Удобные команды
├── README.md
├── CLAUDE.md                        # Инструкции для Claude
└── LEARNING_PATH.md                 # Этот файл

```

---

## Полезные команды

```bash
# Инициализация
go mod init github.com/yourusername/gotasker
go mod tidy

# Разработка
go run src/cmd/main.go
go build -o bin/gotasker src/cmd/main.go

# Тестирование
go test ./...                        # Все тесты
go test -v ./service/...             # Тесты конкретного пакета
go test -cover ./...                 # С покрытием
go test -race ./...                  # Race detector

# Форматирование и линтинг
go fmt ./...
go vet ./...
golangci-lint run

# Swagger
swag init -g src/cmd/main.go

# Docker
docker-compose up --build
docker-compose down
docker-compose logs -f

# Миграции
migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" up
migrate -path migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" down
```

---

## Рекомендуемые ресурсы для изучения

### Официальная документация
- [Go Tour](https://go.dev/tour/) - интерактивный тур по языку
- [Effective Go](https://go.dev/doc/effective_go) - best practices
- [Go by Example](https://gobyexample.com/) - примеры кода

### Книги
- "The Go Programming Language" (Donovan & Kernighan)
- "Go in Action" (Kennedy, Ketelsen, St. Martin)
- "Clean Architecture" (Robert Martin)

### Web frameworks и библиотеки
- [Gin Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Testify](https://github.com/stretchr/testify)

### Паттерны и архитектура
- [Go Patterns](https://github.com/tmrts/go-patterns)
- [Standard Package Layout](https://github.com/golang-standards/project-layout)

---

## Критерии успешного завершения

После прохождения всех этапов вы должны:

✅ **Понимать основы Go**:
- Синтаксис, типы данных, структуры
- Slices, maps, каналы, горутины
- Интерфейсы и composition
- Error handling

✅ **Знать архитектурные паттерны**:
- Clean Architecture слои
- Repository Pattern
- Dependency Injection
- DTO pattern

✅ **Уметь разрабатывать REST API**:
- Создание endpoints
- Валидация данных
- Обработка ошибок
- Middleware

✅ **Работать с базами данных**:
- SQL и PostgreSQL
- Миграции
- Connection pooling
- Транзакции

✅ **Тестировать код**:
- Unit tests
- Integration tests
- Mocking
- Table-driven tests

✅ **Использовать инструменты**:
- Docker и docker-compose
- Git
- Go модули
- Swagger/OpenAPI

---

## Примерное время прохождения

- **Минимальный MVP** (Этапы 0-6): **10-15 часов**
- **С базой данных** (Этапы 0-10): **20-25 часов**
- **Production-ready** (Этапы 0-15): **30-40 часов**
- **С опциональными этапами**: **50+ часов**

---

## Что дальше?

После завершения GoTasker вы будете готовы к:
- Разработке собственных Go API
- Работе на реальных проектах
- Изучению продвинутых тем (gRPC, микросервисы)
- Прохождению собеседований на Go разработчика

Удачи в изучении Go! 🚀
