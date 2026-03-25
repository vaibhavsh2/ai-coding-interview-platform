package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vaibhavsh2/ai-interview/internal/config"
	"github.com/vaibhavsh2/ai-interview/internal/database"
	"github.com/vaibhavsh2/ai-interview/internal/handlers"
	"github.com/vaibhavsh2/ai-interview/internal/middleware"
	"github.com/vaibhavsh2/ai-interview/internal/models"
	"github.com/vaibhavsh2/ai-interview/internal/service"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
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
		&models.CandidateAssignment{},
	)

	// Ensure Admin Account Exists and has ADMIN role
	var adminUser models.User
	if err := db.Where("email = ?", "admin@platform.com").First(&adminUser).Error; err != nil {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		db.Create(&models.User{
			ID:       uuid.New(),
			Email:    "admin@platform.com",
			Password: string(hash),
			Name:     "System Administrator",
			Role:     "ADMIN",
		})
		log.Println("Created default admin account: admin@platform.com / admin123")
	} else if adminUser.Role != "ADMIN" {
        hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
        adminUser.Role = "ADMIN"
        adminUser.Password = string(hash)
        db.Save(&adminUser)
        log.Println("Upgraded existing admin@platform.com to ADMIN role and reset password to admin123")
    }
	questionHandler := handlers.NewQuestionHandler(db)
	testCaseHandler := handlers.NewTestCaseHandler(db)
	submissionHandler := handlers.NewSubmissionHandler(db)

	r := gin.Default()
	
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	r.POST("/questions", middleware.AuthMiddleware(), middleware.RequireAdmin(), questionHandler.CreateQuestion)
	r.GET("/questions/:id", middleware.AuthMiddleware(), questionHandler.GetQuestionByID)
	r.GET("/questions", middleware.AuthMiddleware(), questionHandler.GetAllQuestions)
	r.POST("/questions/:id/testcases", middleware.AuthMiddleware(), middleware.RequireAdmin(), testCaseHandler.CreateTestCase)
	r.GET("/questions/:id/testcases", middleware.AuthMiddleware(), testCaseHandler.GetTestCasesByQuestion)
	r.GET("/submissions/:id", middleware.AuthMiddleware(), submissionHandler.GetSubmissionByID)

	assignmentHandler := handlers.NewAssignmentHandler(db)
	r.POST("/assignments", middleware.AuthMiddleware(), middleware.RequireAdmin(), assignmentHandler.CreateAssignment)
	r.GET("/assignments", middleware.AuthMiddleware(), assignmentHandler.GetAssignments)
	r.POST("/assignments/:id/assign", middleware.AuthMiddleware(), middleware.RequireAdmin(), assignmentHandler.AssignCandidate)
	r.GET("/assignments/:id/candidates", middleware.AuthMiddleware(), middleware.RequireAdmin(), assignmentHandler.GetAssignedCandidates)
	r.GET("/assignments/:id/candidates/:candidate_id/report", middleware.AuthMiddleware(), middleware.RequireAdmin(), assignmentHandler.GetCandidateReport)
	r.DELETE("/assignments/:id", middleware.AuthMiddleware(), middleware.RequireAdmin(), assignmentHandler.DeleteAssignment)
	authHandler := handlers.NewAuthHandler(db)

	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/candidate-login", authHandler.CandidateLogin)
	r.GET("/assignments/:id/questions", middleware.AuthMiddleware(), assignmentHandler.GetQuestionsByAssignment)
	r.POST("/questions/:id/submit", middleware.AuthMiddleware(), submissionHandler.CreateSubmission)
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
