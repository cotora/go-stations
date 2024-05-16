package router

import (
	"database/sql"
	"net/http"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	mux.Handle("/healthz", handler.NewHealthzHandler())
	mux.Handle("/todos", handler.NewTODOHandler(service.NewTODOService(todoDB)))
	//mux.Handle("/do-panic", handler.NewDoPanicHandler())
	mux.Handle("/do-panic",middleware.Recovery(handler.NewDoPanicHandler()))
	mux.Handle("/os-info",middleware.SetOSInfo(handler.NewOSInfoHandler()))
	mux.Handle("/print-log",middleware.SetOSInfo(middleware.PrintLog(handler.NewOSInfoHandler())))
	return mux
}
