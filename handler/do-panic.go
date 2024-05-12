package handler

import (
	"net/http"
)

type DoPanicHandler struct{}

func NewDoPanicHandler() *DoPanicHandler {
	return &DoPanicHandler{}
}

func (h *DoPanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("panic")
}