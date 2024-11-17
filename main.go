package main

import (
	"context"
	"fmt"
	"log"
	"moma-api/cron"
	"moma-api/handler"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ()

func main() {
	mux := http.NewServeMux()

	rateHandler, err := handler.NewRateHandler()
	if err != nil {
		log.Fatalf("fail to create rate handler due to: %v\n", err)
	}
	mux.Handle("/moma-api/rate", rateHandler)

	verifyHandler, err := handler.NewVerifyHandler()
	if err != nil {
		log.Fatalf("fail to create verify handler due to: %v", err)
	}
	mux.Handle("/moma-api/verify_receipt", verifyHandler)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:    addr,
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
	cron.RefreshRates(rateHandler.DB, ticker, closeRateRefresher)

	log.Printf("Starting HTTP server with addr: %v... \n", addr)
	err = server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Printf("Server closed under request due to %v\n", err)
		} else {
			log.Fatalf("Server closed unexpected due to %v\n", err)
		}
	}
}
