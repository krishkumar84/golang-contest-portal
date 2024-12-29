package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	// "strconv"

	"github.com/go-playground/validator/v10"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/storage"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/types"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("New User Handler")
        var user types.User

		err := json.NewDecoder(r.Body).Decode(&user)
		fmt.Println(user)
		if errors.Is(err, io.EOF){ 
			response.WriteJson(w , http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(user); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateUser(user.Name, user.Email, user.Password, user.StudentId,types.RoleUser)

		slog.Info("User created sucessfully",slog.String("userId",fmt.Sprint(lastId)))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"id":lastId})
	}
}