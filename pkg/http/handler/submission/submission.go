package submission

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/krishkumar84/bdcoe-golang-portal/pkg/judge0"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/storage"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/types"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/utils/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateSubmission(storage storage.Storage, judgeClient *judge0.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse submission request
		var submissionReq struct {
			QuestionID  string `json:"question_id"`
			ContestID   string `json:"contest_id"`
			Code        string `json:"code"`
			LanguageID  string `json:"language_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&submissionReq); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		question, err := storage.GetQuestionById(submissionReq.QuestionID)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}

		questionID, _ := primitive.ObjectIDFromHex(submissionReq.QuestionID)
		contestID, _ := primitive.ObjectIDFromHex(submissionReq.ContestID)

		userID := primitive.NewObjectID() 

		submission := types.Submission{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			QuestionID:  questionID,
			ContestID:   contestID,
			Code:        submissionReq.Code,
			LanguageID:  submissionReq.LanguageID,
			Status:      types.StatusPending,
			Score:       0,
			SubmittedAt: time.Now(),
		}

		// Create Judge0 submission request
		judgeReq := judge0.SubmissionRequest{
			SourceCode:     submissionReq.Code,
			LanguageID:     submissionReq.LanguageID,
			TimeLimit:      float64(question[0]["cpu_time_limit"].(int)) / 1000.0, 
			MemoryLimit:    question[0]["memory_limit"].(int),
		}

		testCases := question[0]["test_cases"].([]interface{})
		
		totalScore := 0
		passedTests := 0
		
		for _, tc := range testCases {
			testCase := tc.(map[string]interface{})
			
			// Update judge request with test case input/output
			judgeReq.Stdin = testCase["input"].(string)
			judgeReq.ExpectedOutput = testCase["expected_output"].(string)
			
			token, err := judgeClient.SubmitCode(judgeReq)
			if err != nil {
				response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
				return
			}
			
			var status *judge0.SubmissionStatus
			for i := 0; i < 10; i++ { 
				status, err = judgeClient.GetSubmissionStatus(token)
				if err != nil {
					time.Sleep(time.Second)
					continue
				}
				if status.Status.ID != 1 && status.Status.ID != 2 { // Not In Queue or Processing
					break
				}
				time.Sleep(time.Second)
			}
			
			if status.Status.ID == 3 {
				passedTests++
			}
		}
		
		if len(testCases) > 0 {
			totalScore = (passedTests * 100) / len(testCases)
		}
		
		finalStatus := types.StatusAccepted
		if totalScore < 100 {
			finalStatus = types.StatusWrongAnswer
		}
		
		submission.Status = finalStatus
		submission.Score = totalScore
		
		submissionID, err := storage.CreateSubmission(submission)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]interface{}{
			"submission_id": submissionID,
			"status": finalStatus,
			"score": totalScore,
		})
	}
} 