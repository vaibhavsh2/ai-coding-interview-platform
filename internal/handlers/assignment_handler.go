package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaibhavsh2/ai-interview/internal/models"
	"gorm.io/gorm"
)

type AssignmentHandler struct {
	DB *gorm.DB
}

func NewAssignmentHandler(db *gorm.DB) *AssignmentHandler {
	return &AssignmentHandler{DB: db}
}

func (h *AssignmentHandler) CreateAssignment(c *gin.Context) {

	var assignment models.Assignment

	if err := c.ShouldBindJSON(&assignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assignment"})
		return
	}

	c.JSON(http.StatusCreated, assignment)
}

func (h *AssignmentHandler) GetAssignments(c *gin.Context) {

	var assignments []models.Assignment

	if err := h.DB.Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
		return
	}

	c.JSON(http.StatusOK, assignments)
}
func (h *AssignmentHandler) GetQuestionsByAssignment(c *gin.Context) {

	id := c.Param("id")

	var questions []models.CodingQuestion

	if err := h.DB.Where("assignment_id = ?", id).
		Find(&questions).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch questions",
		})
		return
	}

	c.JSON(http.StatusOK, questions)
}
