package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/lakshsetia/learn-RESTAPI/internal/storage"
	"github.com/lakshsetia/learn-RESTAPI/internal/types"
	"github.com/lakshsetia/learn-RESTAPI/internal/utils/error"
	"github.com/lakshsetia/learn-RESTAPI/internal/utils/json"
)

func UserHandler(storage storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			users, err := storage.GetUsers()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.WriteJSON(w, error.ErrorResponse{
					Level: "database",
					Message: err.Error(),
				})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.WriteJSON(w, users)
		case http.MethodPost:
			var user types.User
			json.ReadJSON(w, r, &user)
			if err := user.Validate(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.WriteJSON(w, error.ErrorResponse{
					Level: "server",
					Message: err.Error(),
				})
			}
			w.WriteHeader(http.StatusCreated)
			if err := storage.CreateUser(user.Name, user.Email, user.Age); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.WriteJSON(w, error.ErrorResponse{
					Level: "database",
					Message: err.Error(),
				})
				return
			}	
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.WriteJSON(w, error.ErrorResponse{
				Level: "backend",
				Message: "invalid method",
			})
		}
	})
}

func UserByIdHandler(storage storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		pathSegments := strings.Split(r.URL.Path, "/")
		if len(pathSegments) != 3 {
			w.WriteHeader(http.StatusBadRequest)
			json.WriteJSON(w, error.ErrorResponse{
				Level: "server",
				Message: "invalid url",
			})
			return
		}
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.WriteJSON(w, error.ErrorResponse{
				Level: "server",
				Message: "invalid url",
			})
			return
		}
		switch r.Method {
		case http.MethodGet:
			user, err := storage.GetUserById(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.WriteJSON(w, error.ErrorResponse{
					Level: "server",
					Message: err.Error(),
				})
				return	
			}
			w.WriteHeader(http.StatusOK)
			json.WriteJSON(w, user)
		case http.MethodPut:
			_, err := storage.GetUserById(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.WriteJSON(w, error.ErrorResponse{
					Level: "server",
					Message: err.Error(),
				})
				return	
			}
			var newUser types.User
			json.ReadJSON(w, r, &newUser)
			if err = newUser.Validate(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.WriteJSON(w, error.ErrorResponse{
					Level: "server",
					Message: err.Error(),
				})
			}
			w.WriteHeader(http.StatusOK)
			if err = storage.UpdateUserById(id, newUser.Name, newUser.Email, newUser.Age); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.WriteJSON(w, error.ErrorResponse{
					Level: "database",
					Message: err.Error(),
				})
				return	
			}
		case http.MethodDelete:
			_, err := storage.GetUserById(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.WriteJSON(w, error.ErrorResponse{
					Level: "server",
					Message: err.Error(),
				})
				return	
			}
			w.WriteHeader(http.StatusNoContent)
			if err = storage.DeleteUserById(id); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.WriteJSON(w, error.ErrorResponse{
					Level: "database",
					Message: err.Error(),
				})
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.WriteJSON(w, error.ErrorResponse{
				Level: "backend",
				Message: "invalid method",
			})
		}
	})
}