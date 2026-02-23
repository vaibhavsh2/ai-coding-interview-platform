package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaibhavsh2/ai-interview/internal/models"
	"gorm.io/gorm"
)

type QuestionHandler struct {
	DB *gorm.DB
}

func NewQuestionHandler(db *gorm.DB) *QuestionHandler {
	return &QuestionHandler{DB: db}
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {

	var question models.CodingQuestion

	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Create(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}

	c.JSON(http.StatusCreated, question)
}

func (h *QuestionHandler) GetAllQuestions(c *gin.Context) {

	var questions []models.CodingQuestion

	if err := h.DB.Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}

	c.JSON(http.StatusOK, questions)
}
