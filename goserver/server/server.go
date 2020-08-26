package server

import (
	"net/http"
)

const (
	launchPath = "/launch"
)

// Server handles incoming HTTP requests.
type Server struct {
	AuthorizedFHIREndPoint string
}

// Run starts HTTP server on port 8080
func (s *Server) Run() {
	http.HandleFunc(launchPath, s.handleLaunch)

	http.ListenAndServe(":8080", nil)
}

func (s *Server) handleLaunch(w http.ResponseWriter, r *http.Request) {
	iss := r.URL.Query()["iss"]
	if len(iss) == 0 || len(iss[0]) < 1 || iss[0] != s.AuthorizedFHIREndPoint {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
