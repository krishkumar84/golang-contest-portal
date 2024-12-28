package auth

import (
	"encoding/json"
	"net/http"
	"time"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/storage"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/types"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/utils/response"
)

func Login(storage storage.Storage, secretKey string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var loginReq types.LoginRequest

        if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
            response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
            return
        }

        if err := validator.New().Struct(loginReq); err != nil {
            validateErrs := err.(validator.ValidationErrors)
            response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
            return
        }

        user, err := storage.GetUserByEmail(loginReq.Email)
        if err != nil {
            response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err))
            return
        }

        if user.Password != loginReq.Password {
            response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid credentials")))
            return
        }

        // Create JWT token
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "user_id": user.ID,
            "student_id": user.StudentId,
            "exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
        })

        tokenString, err := token.SignedString([]byte(secretKey))
        if err != nil {
            response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
            return
        }

        // Set cookie
        http.SetCookie(w, &http.Cookie{
            Name:     "access_token",
            Value:    tokenString,
            HttpOnly: true,
            Secure:   true,
            SameSite: http.SameSiteStrictMode,
            Path:     "/",
            Expires:  time.Now().Add(time.Hour * 24),
        })

        response.WriteJson(w, http.StatusOK, types.TokenResponse{
            AccessToken: tokenString,
            TokenType:   "Bearer",
        })
    }
}