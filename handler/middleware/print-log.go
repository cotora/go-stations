package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Latency  int64 `json:"latency"`
	Path    string	`json:"path"`
	OS 	string	`json:"os"`
}

func PrintLog(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		accessTime := time.Now()
		h.ServeHTTP(w, r)
		latency := time.Since(accessTime).Milliseconds()
		log := Log{
			Timestamp: accessTime,
			Latency: latency,
			Path: r.URL.Path,
			OS: GetOSInfo(r.Context()),
		}

		log_json,err:=json.Marshal(log)
		if err!=nil{
			fmt.Println(err)
			return
		}
		fmt.Println(string(log_json))
	}
	return http.HandlerFunc(fn)
}