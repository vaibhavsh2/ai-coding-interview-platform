package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaibhavsh2/ai-interview/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionHandler struct {
	DB *gorm.DB
}

func NewQuestionHandler(db *gorm.DB) *QuestionHandler {
	return &QuestionHandler{DB: db}
}

// POST /questions
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {

	var req struct {
		AssignmentID string `json:"assignment_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Difficulty  string `json:"difficulty"`

		TestCases []struct {
			Input          string `json:"input"`
			ExpectedOutput string `json:"expectedOutput"`
		} `json:"testcases"`
	}

	// Bind request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Create question
	var assignUUID uuid.UUID
	if req.AssignmentID != "" {
		if parsed, err := uuid.Parse(req.AssignmentID); err == nil {
			assignUUID = parsed
		}
	}

	question := models.CodingQuestion{
		Title:        req.Title,
		Description:  req.Description,
		Difficulty:   req.Difficulty,
		AssignmentID: assignUUID,
	}

	if err := h.DB.Create(&question).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create question"})
		return
	}

	// 🔥 Create testcases automatically
	for _, tc := range req.TestCases {

		testcase := models.TestCase{
			QuestionID:     question.ID,
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
		}

		h.DB.Create(&testcase)
	}

	c.JSON(201, gin.H{
		"message": "Question created with testcases",
		"id":      question.ID,
	})
}

// GET /questions
func (h *QuestionHandler) GetAllQuestions(c *gin.Context) {

	var questions []models.CodingQuestion

	if err := h.DB.Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}

	c.JSON(http.StatusOK, questions)
}

// GET /questions/:id
func (h *QuestionHandler) GetQuestionByID(c *gin.Context) {

	id := c.Param("id")

	var question models.CodingQuestion

	if err := h.DB.First(&question, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, question)
}
