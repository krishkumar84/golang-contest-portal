package testcase

import (
	"encoding/json"
	"net/http"

	"github.com/krishkumar84/bdcoe-golang-portal/pkg/storage"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/types"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/utils/response"
)

func CreateTestCase(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var testCaseReq types.TestCase
		if err := json.NewDecoder(r.Body).Decode(&testCaseReq); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		testCaseId, err := storage.CreateTestCase(testCaseReq)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"test_case_id": testCaseId})
	}
}