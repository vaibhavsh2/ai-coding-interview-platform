package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vaibhavsh2/ai-interview/internal/models"
	"gorm.io/gorm"
)

type TestCaseHandler struct {
	DB *gorm.DB
}

func NewTestCaseHandler(db *gorm.DB) *TestCaseHandler {
	return &TestCaseHandler{DB: db}
}

func (h *TestCaseHandler) CreateTestCase(c *gin.Context) {

	questionIDParam := c.Param("id")

	questionID, err := uuid.Parse(questionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var testCase models.TestCase

	if err := c.ShouldBindJSON(&testCase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	testCase.QuestionID = questionID

	if err := h.DB.Create(&testCase).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test case"})
		return
	}

	c.JSON(http.StatusCreated, testCase)
}

func (h *TestCaseHandler) GetTestCasesByQuestion(c *gin.Context) {

	questionIDParam := c.Param("id")

	questionID, err := uuid.Parse(questionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var testCases []models.TestCase

	if err := h.DB.Where("question_id = ?", questionID).Find(&testCases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test cases"})
		return
	}

	c.JSON(http.StatusOK, testCases)
}
