package main

import (
	"Elections_Patiala/pkg/db"
	"Elections_Patiala/pkg/handlers"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var store = sessions.NewCookieStore([]byte("2eb7ddef6411bca4205d73dbfcbf9115fcf2ec43"))

func main() {
	db.InitDatabase()

	r := mux.NewRouter()
	r.Use(sessionMiddleware)
	r.HandleFunc("/", handlers.HandleIndex).Methods("GET")
	r.HandleFunc("/admin/login", handlers.HandlerAdminLogin).Methods("GET")
	r.HandleFunc("/admin/login", handlers.HandleAuthenticate).Methods("POST")
	r.HandleFunc("/admin/dashboard", handlers.HanldeAdminDashboard).Methods("GET")

	corsMiddlerware := cors.Default()
	handler := corsMiddlerware.Handler(r)

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

func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			log.Println("Error getting session:", err)
		}
		defer session.Save(r, w)
		next.ServeHTTP(w, r)
	})
}
