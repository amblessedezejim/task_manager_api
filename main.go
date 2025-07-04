package main

import (
	"errors"
	"log"

	"github.com/amblessedezejim/task_manager_api/config"
	"github.com/amblessedezejim/task_manager_api/handlers.go"
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
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	api.GET("/tasks", handlers.GetTasks)
	api.GET("/tasks/:id", handlers.GetTaskById)
	api.POST("/tasks", handlers.CreateTask)
	api.DELETE("/tasks/:id", handlers.DeleteTask)
	api.PUT("/tasks/:id", handlers.UpdateTask)
	port := ":8080"
	log.Printf("Server starting on %s\n", port)
	router.Run(port)
}
