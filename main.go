package main

import (
	"Elections_Patiala/db"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	db.InitDatabase()

	router := mux.NewRouter()

	corsMiddlerware := cors.Default()
	handler := corsMiddlerware.Handler(router)

	server := &http.Server{Addr: ":8080", Handler: handler}
	go func() {
		fmt.Println("Server is up and running on port : 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Error while starting the server")
		}
	}()
	WaitForTerminationSignal(server)
}

func WaitForTerminationSignal(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}
	log.Println("Server stopped gracefully")
}
