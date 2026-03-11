package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vaibhavsh2/ai-interview/internal/config"
	"github.com/vaibhavsh2/ai-interview/internal/database"
	"github.com/vaibhavsh2/ai-interview/internal/handlers"
	"github.com/vaibhavsh2/ai-interview/internal/models"
	"github.com/vaibhavsh2/ai-interview/internal/service"
)

func main() {

	cfg := config.LoadConfig()
	db := database.Connect(cfg)
	db.AutoMigrate(
		&models.Assignment{},
		&models.CodingQuestion{},
		&models.TestCase{},
		&models.Submission{},
		&models.User{},
	)
	questionHandler := handlers.NewQuestionHandler(db)
	testCaseHandler := handlers.NewTestCaseHandler(db)
	submissionHandler := handlers.NewSubmissionHandler(db)

	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/questions", questionHandler.CreateQuestion)
	r.GET("/questions/:id", questionHandler.GetQuestionByID)
	r.GET("/questions", questionHandler.GetAllQuestions)
	r.POST("/questions/:id/testcases", testCaseHandler.CreateTestCase)
	r.GET("/questions/:id/testcases", testCaseHandler.GetTestCasesByQuestion)
	r.GET("/submissions/:id", submissionHandler.GetSubmissionByID)

	assignmentHandler := handlers.NewAssignmentHandler(db)
	r.POST("/assignments", assignmentHandler.CreateAssignment)
	r.GET("/assignments", assignmentHandler.GetAssignments)
	authHandler := handlers.NewAuthHandler(db)

	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)
	r.GET("/assignments/:id/questions", assignmentHandler.GetQuestionsByAssignment)
	r.POST("/questions/:id/submit", submissionHandler.CreateSubmission)
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
	worker := service.NewExecutionWorker(db)
	go worker.Start()
	r.Run(":8080")

	_ = db
}
