package test

import (
	"net/http"

	"github.com/krishkumar84/bdcoe-golang-portal/pkg/utils/response"
)

func TestUserRoute(w http.ResponseWriter, r *http.Request) {
    response.WriteJson(w, http.StatusOK, map[string]string{"message": "User route accessed successfully"})
}

func TestAdminRoute(w http.ResponseWriter, r *http.Request) {
    response.WriteJson(w, http.StatusOK, map[string]string{"message": "Admin route accessed successfully"})
} 