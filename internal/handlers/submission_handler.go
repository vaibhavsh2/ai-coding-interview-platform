package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vaibhavsh2/ai-interview/internal/models"
	"gorm.io/gorm"
)

type SubmissionHandler struct {
	DB *gorm.DB
}

func NewSubmissionHandler(db *gorm.DB) *SubmissionHandler {
	return &SubmissionHandler{DB: db}
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
	submission.Status = "SUBMITTED"

	if err := h.DB.Create(&submission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create submission"})
		return
	}

	c.JSON(http.StatusCreated, submission)
}
