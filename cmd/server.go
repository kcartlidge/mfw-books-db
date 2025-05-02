package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// Server is a simple web server that serves the database
type Server struct {
	Port          int
	Router        *mux.Router
	Filename      string
	CookieHandler *CookieHandler
}

// NewServer creates a new server
func NewServer(port int, filename string) (*Server, error) {
	// Initialize templates
	_, err := NewTemplates()
	if err != nil {
		return nil, fmt.Errorf("error initializing templates: %w", err)
	}

	// Initialize cookie handler
	cookieHandler, err := NewCookieHandler()
	if err != nil {
		return nil, fmt.Errorf("error initializing cookie handler: %w", err)
	}

	s := &Server{
		Port:          port,
		Router:        mux.NewRouter(),
		Filename:      filename,
		CookieHandler: cookieHandler,
	}

	// Add handlers
	s.Router.HandleFunc("/", s.HomeHandler).Methods("GET")
	s.Router.HandleFunc("/sort/{field}", s.SortHandler).Methods("GET")
	s.Router.HandleFunc("/books/edit/{isbn}", s.EditHandler).Methods("GET")
	s.Router.HandleFunc("/books/save/{isbn}", s.SaveHandler).Methods("POST")

	// Add root-level static file handler (must come after specific routes)
	s.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	s.Router.NotFoundHandler = http.HandlerFunc(s.NotFoundHandler)

	return s, nil
}

// Start starts the server
func (s *Server) Start() {
	// Create server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: s.Router,
	}

	// Channel to listen for errors coming from the server
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		fmt.Println()
		log.Printf("Server starting on %s", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking select waiting for either a server error or a signal
	select {
	case err := <-serverErrors:
		fmt.Println()
		log.Printf("Error starting server: %v", err)

	case sig := <-shutdown:
		fmt.Println()
		log.Printf("Shutdown signal received: %v", sig)

		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Could not stop server gracefully: %v", err)
			if err := srv.Close(); err != nil {
				log.Printf("Could not stop server: %v", err)
			}
		}
	}

	// Done
	fmt.Println()
	fmt.Println("Server shutdown complete")
	fmt.Println()
}
