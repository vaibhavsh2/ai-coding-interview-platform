package service

import (
	"log"
	"time"

	"github.com/vaibhavsh2/ai-interview/internal/models"
	"gorm.io/gorm"
)

type ExecutionWorker struct {
	DB *gorm.DB
}

func NewExecutionWorker(db *gorm.DB) *ExecutionWorker {
	return &ExecutionWorker{DB: db}
}

func (w *ExecutionWorker) Start() {

	for {

		var submissions []models.Submission

		err := w.DB.Where("status = ?", "SUBMITTED").
			Limit(5).
			Find(&submissions).Error

		if err != nil {
			log.Println("Worker DB error:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, sub := range submissions {

			log.Println("Processing submission:", sub.ID)

			w.processSubmission(&sub)

		}

		time.Sleep(2 * time.Second)
	}

}
func (w *ExecutionWorker) processSubmission(sub *models.Submission) {

	sub.Status = "RUNNING"
	w.DB.Save(sub)

	// temporary fake execution
	time.Sleep(3 * time.Second)

	sub.Status = "COMPLETED"
	sub.PassedTests = 1
	sub.TotalTests = 1

	w.DB.Save(sub)
}
