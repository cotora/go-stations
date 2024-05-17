package main

import (
	"log"
	"os"
	"time"
	"net/http"
	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
	"os/signal"
	"context"
	"sync"
	"errors"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	// NOTE: シグナルを受け取るためのコンテキストを作成
	sigCtx,stop:=signal.NotifyContext(context.Background(),os.Interrupt,os.Kill)
	defer stop()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB)

	// NOTE: サーバーを作成する
	srv:=&http.Server{
		Addr: port,
		Handler: mux,
	}

	wg:=&sync.WaitGroup{}
	wg.Add(1)

	// NOTE: シグナルを受け取ったらサーバーをシャットダウンする
	go func(){
		defer wg.Done()
		<-sigCtx.Done()
		log.Println("received signal, shutting down gracefully")
		ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
		defer cancel()
		err:=srv.Shutdown(ctx)
		if err!=nil{
			log.Println("failed to shutdown gracefully, err =", err)
		}
	}()

	// TODO: サーバーをlistenする
	err = srv.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("main: server closed")
		} else {
			return err
		}
	}

	wg.Wait()

	return nil
}
