package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vaibhavsh2/ai-interview/internal/config"
	"github.com/vaibhavsh2/ai-interview/internal/database"
	"github.com/vaibhavsh2/ai-interview/internal/handlers"
	"github.com/vaibhavsh2/ai-interview/internal/models"
	"log"
)

func main() {

	// Load configuration
	cfg := config.LoadConfig()

	// Connect database
	db := database.Connect(cfg)
	db.AutoMigrate(
		&models.CodingQuestion{},
		&models.TestCase{},
		&models.Submission{},
	)
	questionHandler := handlers.NewQuestionHandler(db)
	testCaseHandler := handlers.NewTestCaseHandler(db)
	submissionHandler := handlers.NewSubmissionHandler(db)
	// Create Gin router
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/questions", questionHandler.CreateQuestion)
	r.GET("/questions", questionHandler.GetAllQuestions)
	r.POST("/questions/:id/testcases", testCaseHandler.CreateTestCase)
	r.GET("/questions/:id/testcases", testCaseHandler.GetTestCasesByQuestion)
	r.GET("/submissions/:id", submissionHandler.GetSubmissionByID)

	r.POST("/questions/:id/submit", submissionHandler.CreateSubmission)
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	log.Println("Server running on port 8080")
	r.GET("/submissions", func(c *gin.Context) {
		var subs []models.Submission
		db.Find(&subs)
		c.JSON(200, subs)
	})
	r.Run(":8080")

	_ = db
}
