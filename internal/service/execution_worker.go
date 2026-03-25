package service

import (
	"log"
	"strings"
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

		for i := range submissions {
			log.Println("Processing submission:", submissions[i].ID)
			w.processSubmission(&submissions[i])
		}

		time.Sleep(2 * time.Second)
	}
}

// 🔥 REAL EXECUTION USING JUDGE0
func (w *ExecutionWorker) processSubmission(sub *models.Submission) {

	// mark running
	sub.Status = "RUNNING"
	w.DB.Save(sub)

	// get testcases
	var testCases []models.TestCase
	err := w.DB.Where("question_id = ?", sub.QuestionID).Find(&testCases).Error
	if err != nil {
		log.Println("Error fetching testcases:", err)
		sub.Status = "FAILED"
		w.DB.Save(sub)
		return
	}

	passed := 0
	total := len(testCases)

	// Map Submission Language string to Judge0 Language ID
	langID := 62 // Default to Java
	switch strings.ToLower(sub.Language) {
	case "python":
		langID = 71
	case "cpp", "c++":
		langID = 54
	case "c":
		langID = 50
	case "javascript", "js":
		langID = 93
	case "java":
		langID = 62
	}

	for _, test := range testCases {

		// 🔥 send to Judge0
		token, err := SubmitToJudge(sub.SourceCode, test.Input, langID)
		if err != nil {
			log.Println("Judge submit error:", err)
			continue
		}

		var result JudgeResult

		// 🔥 Poll Judge0 for result up to 5 times
		for tries := 0; tries < 5; tries++ {
			time.Sleep(2 * time.Second)
			
			res, err := GetJudgeResult(token)
			if err != nil {
				log.Println("Judge result poll error:", err)
				continue
			}
			
			// If not pending, we can stop polling
			if res.Status.Description != "In Queue" && res.Status.Description != "Processing" {
				result = res
				break
			}
		}

		output := strings.TrimSpace(result.Stdout)
		expected := strings.TrimSpace(test.ExpectedOutput)

		output = strings.ReplaceAll(output, "\r\n", "\n")
		expected = strings.ReplaceAll(expected, "\r\n", "\n")

		log.Println("Output:", output)
		log.Println("Expected:", expected)
		log.Println("Status:", result.Status.Description)

		if result.Status.Description == "Accepted" && output == expected {
			passed++
		}
	}

	// update result
	sub.PassedTests = passed
	sub.TotalTests = total

	if passed == total {
		sub.Status = "COMPLETED"
	} else {
		sub.Status = "FAILED"
	}

	w.DB.Save(sub)
}
