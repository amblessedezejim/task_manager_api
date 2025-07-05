package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/amblessedezejim/task_manager_api/config"
	"github.com/amblessedezejim/task_manager_api/handlers.go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %s", err.Error())
	}

	// Connect to database
	config.InitDB()
	if config.DB == nil {
		panic(errors.New("Connection to database failed"))
	}
	log.Println("Connected to Database successfully")
	defer config.CloseDB()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Context-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(customErrorHandler())

	taskHandler := handlers.NewTaskHandler(config.DB)
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Task manager API is running fine",
			"time":    time.Now().UTC(),
		})
	})

	router.SetTrustedProxies([]string{
		"127.0.0.1",
		"192.168.0.0/16",
		"10.0.0.0/8",
	})

	api := router.Group("/api/v1")
	api.GET("/tasks", taskHandler.GetTasks)
	api.GET("/tasks/:id", taskHandler.GetTaskById)
	api.POST("/tasks", taskHandler.CreateTask)
	api.DELETE("/tasks/:id", taskHandler.DeleteTask)
	api.PUT("/tasks/:id", taskHandler.UpdateTask)
	port := config.GetServerPort()
	log.Printf("Server starting on %s\n", port)
	log.Printf("Health check available at: http://localhost%s/health", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start sever:", err)
	}
}

func customErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(ctx *gin.Context, recoverd any) {
		if err, ok := recoverd.(string); ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Internal server error,",
				"error":   err,
			})
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
	})
}
