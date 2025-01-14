package question

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/krishkumar84/bdcoe-golang-portal/pkg/storage"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/types"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/utils/response"
)

func CreateQuestion(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var questionReq types.Question
		if err := json.NewDecoder(r.Body).Decode(&questionReq); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		questionId, err := storage.CreateQuestion(questionReq)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"question_id": questionId})
	}
}

func GetQuestionById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		id := path[strings.LastIndex(path, "/")+1:]
		
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("contest id is required")))
			return
		}
		fmt.Println("Question ID: ", id)
		question, err := storage.GetQuestionById(id)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, question)
	}
}

func AddTestCaseToQuestion(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		parts := strings.Split(path, "/")
		if len(parts) < 5 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid URL format")))
			return
		}
		questionId := parts[3]
		
		fmt.Printf("Extracted question ID: %s\n", questionId)
		
		if questionId == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("question id is required")))
			return
		}

		var testCase types.TestCase
		if err := json.NewDecoder(r.Body).Decode(&testCase); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		testCaseId, err := storage.AddTestCaseToQuestion(questionId, testCase)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"test_case_id": testCaseId})
	}
}

func DeleteTestCaseFromQuestionById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		parts := strings.Split(path, "/")
		if len(parts) < 6 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid URL format")))
			return
		}
		questionId := parts[3]
		testCaseId := parts[5]
		
		fmt.Printf("Extracted question ID: %s\n", questionId)
		fmt.Printf("Extracted test case ID: %s\n", testCaseId)
		
		if questionId == "" || testCaseId == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("question id and test case id are required")))
			return
		}

		err := storage.DeleteTestCaseFromQuestionById(questionId, testCaseId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "test case deleted successfully"})
	}
}

func EditQuestionById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var questionReq types.Question
		if err := json.NewDecoder(r.Body).Decode(&questionReq); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		path := r.URL.Path
		id := path[strings.LastIndex(path, "/")+1:]
		
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("question id is required")))
			return
		}

		err := storage.EditQuestionById(id, questionReq)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "question updated successfully"})
	}
}