package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vaibhav/ai-interview/internal/config"
	"github.com/vaibhav/ai-interview/internal/database"
	"github.com/vaibhav/ai-interview/internal/handlers"
	"github.com/vaibhav/ai-interview/internal/models"
)

func main() {

	// Load configuration
	cfg := config.LoadConfig()

	// Connect database
	db := database.Connect(cfg)
	db.AutoMigrate(
		&models.CodingQuestion{},
		&models.TestCase{},
	)
	questionHandler := handlers.NewQuestionHandler(db)
	// Create Gin router
	r := gin.Default()
	r.POST("/questions", questionHandler.CreateQuestion)
	r.GET("/questions", questionHandler.GetAllQuestions)
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	log.Println("Server running on port 8080")

	r.Run(":8080")

	_ = db
}
