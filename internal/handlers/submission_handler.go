package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vaibhavsh2/ai-interview/internal/models"
	"github.com/vaibhavsh2/ai-interview/internal/service"
	"github.com/vaibhavsh2/ai-interview/internal/state"
	"gorm.io/gorm"
)

type SubmissionHandler struct {
	DB       *gorm.DB
	Executor *service.ExecutionService
}

func NewSubmissionHandler(db *gorm.DB) *SubmissionHandler {
	return &SubmissionHandler{
		DB:       db,
		Executor: service.NewExecutionService(db),
	}
}

// POST /questions/:id/submit
func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {

	questionIDParam := c.Param("id")

	questionID, err := uuid.Parse(questionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var submission models.Submission

	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	submission.QuestionID = questionID
	submission.Status = state.StatusSubmitted

	if err := h.DB.Create(&submission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create submission"})
		return
	}

	// Async execution
	go h.Executor.ExecuteSubmission(submission.ID)

	c.JSON(http.StatusCreated, submission)
}

// GET /submissions/:id
func (h *SubmissionHandler) GetSubmissionByID(c *gin.Context) {

	submissionIDParam := c.Param("id")

	submissionID, err := uuid.Parse(submissionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	var submission models.Submission

	if err := h.DB.First(&submission, "id = ?", submissionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	c.JSON(http.StatusOK, submission)
}
