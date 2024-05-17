package handler

import (
	"net/http"
	"time"
	"log"
)

func SleepFive() http.Handler {

	fn:=func(w http.ResponseWriter,r *http.Request){
		time.Sleep(5*time.Second)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, World!"))
		log.Println("sleep-five response sent")
	}

	return http.HandlerFunc(fn)
}