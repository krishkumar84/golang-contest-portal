package contest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/krishkumar84/bdcoe-golang-portal/pkg/storage"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/types"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/utils/response"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateContest(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contestReq types.Contest
		if err := json.NewDecoder(r.Body).Decode(&contestReq); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		contestId, err := storage.CreateContest(contestReq)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"contest_id": contestId})
	}
}

func DeleteContestById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		id := path[strings.LastIndex(path, "/")+1:]
		
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("contest id is required")))
			return
		}

		if err := storage.DeleteContestById(id); err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"status": "success", "message": "contest deleted successfully"})
	}
}

func GetAllContests(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contests, err := storage.GetAllContests()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		
		response.WriteJson(w, http.StatusOK, contests)
	}
}

func GetContestById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		id := path[strings.LastIndex(path, "/")+1:]
		
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("contest id is required")))
			return
		}

		contest, err := storage.GetContestById(id)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("contest not found")))
				return
			}
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, contest)
	}
}

func AddQuestionToContest(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		parts := strings.Split(path, "/")
		if len(parts) < 5 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid URL format")))
			return
		}
		contestId := parts[3]
		
		fmt.Printf("Extracted contest ID: %s\n", contestId)
		
		if contestId == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("contest id is required")))
			return
		}

		var question types.Question
		if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		questionId, err := storage.AddQuestionToContest(contestId, question)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"question_id": questionId})
	}
}

func DeleteQuestionFromContestById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		parts := strings.Split(path, "/")
		if len(parts) < 6 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid URL format")))
			return
		}
		contestId := parts[3]
		questionId := parts[5]
		
		fmt.Printf("Extracted contest ID: %s\n", contestId)
		fmt.Printf("Extracted question ID: %s\n", questionId)
		
		if contestId == "" || questionId == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("contest id and question id are required")))
			return
		}

		if err := storage.DeleteQuestionFromContestById(contestId, questionId); err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"status": "success", "message": "question deleted from contest successfully"})
	}
}