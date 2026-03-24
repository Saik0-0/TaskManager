package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (server *Server) StatsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		response := server.Store.GetStats()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			fmt.Println("Encoding error: ", err)
			return
		}

	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid method"}); err != nil {
			fmt.Println("Encoding error: ", err)
			return
		}
	}
}
