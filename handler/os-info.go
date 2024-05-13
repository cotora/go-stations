package handler

import (
	"net/http"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"encoding/json"
)

type OSInfoResponse struct {
	OS string `json:"os"`
}

type OSInfoHandler struct{}

func NewOSInfoHandler() *OSInfoHandler {
	return &OSInfoHandler{}
}

func (h *OSInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	os := middleware.GetOSInfo(r.Context())
	resp := OSInfoResponse{OS: os}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

