package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	role, exists := c.Get("role")
	var assignments []models.Assignment

	if exists && role == "ADMIN" {
		if err := h.DB.Find(&assignments).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
			return
		}
	} else {
		userID, userExists := c.Get("user_id")
		if !userExists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var allocated []models.CandidateAssignment
		if err := h.DB.Where("candidate_id = ?", userID).Find(&allocated).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch candidate allocations"})
			return
		}

		var allowedIDs []string
		for _, alloc := range allocated {
			allowedIDs = append(allowedIDs, alloc.AssignmentID)
		}

		if len(allowedIDs) == 0 {
			c.JSON(http.StatusOK, []models.Assignment{})
			return
		}

		if err := h.DB.Where("id IN ?", allowedIDs).Find(&assignments).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
			return
		}
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

// POST /assignments/:id/assign
func (h *AssignmentHandler) AssignCandidate(c *gin.Context) {

	assignmentID := c.Param("id")

	var req struct {
		CandidateID string `json:"candidate_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Lookup custom magic ID (Email)
	var user models.User
	if err := h.DB.Where("email = ?", req.CandidateID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Magic ID not found. Ensure Candidate has simulated a sign up."})
		return
	}

	// Create map allocation
	allocation := models.CandidateAssignment{
		ID:           uuid.New(),
		CandidateID:  user.ID,
		AssignmentID: assignmentID,
		Status:       "ASSIGNED",
	}

	if err := h.DB.Create(&allocation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to explicitly allocate assignment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Candidate securely mapped!"})
}

// DELETE /assignments/:id
func (h *AssignmentHandler) DeleteAssignment(c *gin.Context) {

	assignmentID := c.Param("id")

	// 1. Delete associated candidate mappings
	h.DB.Where("assignment_id = ?", assignmentID).Delete(&models.CandidateAssignment{})

	// 2. Fetch associated questions to cascade delete testcases
	var questions []models.CodingQuestion
	h.DB.Where("assignment_id = ?", assignmentID).Find(&questions)

	for _, q := range questions {
		h.DB.Where("question_id = ?", q.ID).Delete(&models.TestCase{})
		h.DB.Delete(&q)
	}

	// 3. Delete parent Assignment
	if err := h.DB.Where("id = ?", assignmentID).Delete(&models.Assignment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute cascade assignment wipe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully shredded assignment hierarchy"})
}

// GET /assignments/:id/candidates
func (h *AssignmentHandler) GetAssignedCandidates(c *gin.Context) {

	assignmentID := c.Param("id")

	var allocations []models.CandidateAssignment
	if err := h.DB.Where("assignment_id = ?", assignmentID).Find(&allocations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch candidate allocations"})
		return
	}

	var candidateIDs []string
	for _, alloc := range allocations {
		candidateIDs = append(candidateIDs, alloc.CandidateID.String())
	}

	if len(candidateIDs) == 0 {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	var users []models.User
	if err := h.DB.Where("id IN ?", candidateIDs).Select("id, name, email").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch candidate details"})
		return
	}

	var questions []models.CodingQuestion
	h.DB.Where("assignment_id = ?", assignmentID).Find(&questions)

	var questionIDs []uuid.UUID
	for _, q := range questions {
		questionIDs = append(questionIDs, q.ID)
	}

	var allSubmissions []models.Submission
	if len(questionIDs) > 0 {
		h.DB.Where("candidate_id IN ? AND question_id IN ?", candidateIDs, questionIDs).
			Order("created_at desc").Find(&allSubmissions)
	}

	var result []gin.H
	for _, u := range users {
		passed := 0
		attempted := 0

		for _, qID := range questionIDs {
			for _, s := range allSubmissions {
				if s.CandidateID.String() == u.ID.String() && s.QuestionID == qID {
					attempted++
					if s.Status == "COMPLETED" && s.PassedTests == s.TotalTests && s.TotalTests > 0 {
						passed++
					}
					break
				}
			}
		}

		result = append(result, gin.H{
			"id":        u.ID,
			"name":      u.Name,
			"email":     u.Email,
			"passed":    passed,
			"attempted": attempted,
			"total":     len(questions),
		})
	}
	c.JSON(http.StatusOK, result)
}

// GET /assignments/:id/candidates/:candidate_id/report
func (h *AssignmentHandler) GetCandidateReport(c *gin.Context) {

	assignmentID := c.Param("id")
	candidateID := c.Param("candidate_id")

	// Verify the Candidate is actually assigned to this Assignment
	var alloc models.CandidateAssignment
	if err := h.DB.Where("assignment_id = ? AND candidate_id = ?", assignmentID, candidateID).First(&alloc).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Candidate is not assigned to this assignment"})
		return
	}

	// Fetch all questions for this assignment
	var questions []models.CodingQuestion
	if err := h.DB.Where("assignment_id = ?", assignmentID).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignment questions"})
		return
	}

	var questionIDs []uuid.UUID
	for _, q := range questions {
		questionIDs = append(questionIDs, q.ID)
	}

	if len(questionIDs) == 0 {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}

	// Fetch the absolute latest submissions for this specific candidate against these specific mapped questions
	var submissions []models.Submission
	if err := h.DB.Where("candidate_id = ? AND question_id IN ?", candidateID, questionIDs).Order("created_at desc").Find(&submissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch candidate submissions"})
		return
	}

	var report []gin.H

	for _, q := range questions {
		var bestSub *models.Submission
		for i, s := range submissions {
			if s.QuestionID == q.ID {
				// We ordered strictly by desc in SQL, so the first match hitting this block is undeniably the absolute latest attempt
				bestSub = &submissions[i]
				break
			}
		}

		if bestSub != nil {
			report = append(report, gin.H{
				"question_title": q.Title,
				"difficulty":     q.Difficulty,
				"status":         bestSub.Status,
				"language":       bestSub.Language,
				"source_code":    bestSub.SourceCode,
				"passed_tests":   bestSub.PassedTests,
				"total_tests":    bestSub.TotalTests,
				"execution_time": bestSub.ExecutionTime,
				"memory_used":    bestSub.MemoryUsed,
				"submitted_at":   bestSub.CreatedAt,
			})
		} else {
			report = append(report, gin.H{
				"question_title": q.Title,
				"difficulty":     q.Difficulty,
				"status":         "NOT_ATTEMPTED",
			})
		}
	}

	c.JSON(http.StatusOK, report)
}
