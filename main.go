package main

import (
	"context"
	"log"
	"moma-api/cron"
	"moma-api/handler"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/moma-api/rate", handler.NewRateHandler())

	server := &http.Server{
		Addr:    ":80",
		Handler: mux,
	}

	// 创建系统信号接收器
	done := make(chan os.Signal)
	closeRateRefresher := make(chan bool)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done

		closeRateRefresher <- true
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal("Shutdown server:", err)
		}
	}()

	// create cron job to refresh rates
	ticker := time.NewTicker(time.Hour * 24)
	cron.RefreshRates(ticker, closeRateRefresher)

	log.Println("Starting HTTP server...")
	err := server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Print("Server closed under request")
		} else {
			log.Fatal("Server closed unexpected")
		}
	}
}
