package service

import (
	"github.com/google/uuid"
	"github.com/vaibhavsh2/ai-interview/internal/models"
	"github.com/vaibhavsh2/ai-interview/internal/state"
	"gorm.io/gorm"
)

type ExecutionService struct {
	DB *gorm.DB
}

func NewExecutionService(db *gorm.DB) *ExecutionService {
	return &ExecutionService{DB: db}
}

func (e *ExecutionService) ExecuteSubmission(submissionID uuid.UUID) error {

	var submission models.Submission

	if err := e.DB.First(&submission, "id = ?", submissionID).Error; err != nil {
		return err
	}

	// Update status → RUNNING
	submission.Status = state.StatusRunning
	e.DB.Save(&submission)

	// Fetch test cases
	var testCases []models.TestCase
	e.DB.Where("question_id = ?", submission.QuestionID).Find(&testCases)

	total := len(testCases)
	passed := 0

	// 🔥 For now we simulate success
	for range testCases {
		passed++
	}

	submission.TotalTests = total
	submission.PassedTests = passed
	submission.Status = state.StatusCompleted

	e.DB.Save(&submission)

	return nil
}
