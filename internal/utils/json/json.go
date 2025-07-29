package json

import (
	"encoding/json"
	"net/http"

	"github.com/lakshsetia/learn-RESTAPI/internal/utils/error"
)

func WriteJSON(w http.ResponseWriter, response any) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func ReadJSON(w http.ResponseWriter, r *http.Request, request any) {
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteJSON(w, error.ErrorResponse{
			Level: "backend",
			Message: "invalid json request",
		})
	}
}