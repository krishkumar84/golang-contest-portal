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