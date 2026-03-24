package handlers

import (
	"net/http"
)

func (server *Server) StatsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		response := server.Store.GetStats()
		writeJSON(w, http.StatusOK, response)

	default:
		writeError(w, http.StatusMethodNotAllowed, "Invalid method")
	}
}
