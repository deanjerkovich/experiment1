package main

import (
	"auth-server/pkg/base64util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// User represents a user account
type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"-"` // Don't include password in JSON responses
	Created  time.Time `json:"created"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// Response represents a generic API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Server represents the HTTP server
type Server struct {
	authHandler *AuthHandler
	mutex       sync.RWMutex
}

// NewServer creates a new server instance
func NewServer() *Server {
	return &Server{
		authHandler: NewAuthHandler([]byte("0mgn3wcryptok3y")),
		mutex:       sync.RWMutex{},
	}
}

// registerHandler delegates to AuthHandler
func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	s.authHandler.RegisterHandler(w, r)
}

// loginHandler delegates to AuthHandler
func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	s.authHandler.LoginHandler(w, r)
}

// logoutHandler delegates to AuthHandler
func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	s.authHandler.LogoutHandler(w, r)
}

// profileHandler delegates to AuthHandler
func (s *Server) profileHandler(w http.ResponseWriter, r *http.Request) {
	s.authHandler.ProfileHandler(w, r)
}

// changePasswordHandler delegates to AuthHandler
func (s *Server) changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	s.authHandler.ChangePasswordHandler(w, r)
}

// base64EncodeHandler handles base64 encoding requests
func (s *Server) base64EncodeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[DEBUG] Base64 encode request received\n")

	if r.Method != http.MethodPost {
		fmt.Fprintf(os.Stderr, "[DEBUG] Invalid method: %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] Empty text provided\n")
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	encoder := base64util.NewEncoder()
	encoded, err := encoder.Encode(req.Text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Encoding failed: %v\n", err)
		http.Error(w, "Encoding failed", http.StatusInternalServerError)
		return
	}

	response := Response{
		Success: true,
		Message: "Text encoded successfully",
		Data: map[string]interface{}{
			"original": req.Text,
			"encoded":  encoded,
		},
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] Base64 encoding successful for text: %s\n", req.Text)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// base64DecodeHandler handles base64 decoding requests
func (s *Server) base64DecodeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[DEBUG] Base64 decode request received\n")

	if r.Method != http.MethodPost {
		fmt.Fprintf(os.Stderr, "[DEBUG] Invalid method: %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] Empty text provided\n")
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	encoder := base64util.NewEncoder()
	decoded, err := encoder.Decode(req.Text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Decoding failed: %v\n", err)
		http.Error(w, "Invalid base64 text", http.StatusBadRequest)
		return
	}

	response := Response{
		Success: true,
		Message: "Text decoded successfully",
		Data: map[string]interface{}{
			"original": req.Text,
			"decoded":  decoded,
		},
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] Base64 decoding successful for text: %s\n", req.Text)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// healthHandler provides a health check endpoint
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Success: true,
		Message: "Server is healthy",
		Data: map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"status":    "running",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	fmt.Fprintf(os.Stderr, "[DEBUG] Starting authentication server...\n")

	server := NewServer()
	fmt.Fprintf(os.Stderr, "[DEBUG] Server instance created\n")

	// Create router
	router := mux.NewRouter()
	fmt.Fprintf(os.Stderr, "[DEBUG] Router created\n")

	// API routes
	router.HandleFunc("/api/register", server.registerHandler).Methods("POST")
	router.HandleFunc("/api/login", server.loginHandler).Methods("POST")
	router.HandleFunc("/api/logout", server.logoutHandler).Methods("POST")
	router.HandleFunc("/api/profile", server.profileHandler).Methods("GET")
	router.HandleFunc("/api/change-password", server.changePasswordHandler).Methods("POST")
	router.HandleFunc("/api/base64/encode", server.base64EncodeHandler).Methods("POST")
	router.HandleFunc("/api/base64/decode", server.base64DecodeHandler).Methods("POST")
	router.HandleFunc("/api/health", server.healthHandler).Methods("GET")
	fmt.Fprintf(os.Stderr, "[DEBUG] All API routes registered\n")

	// Serve static files (optional - for a simple frontend)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))
	fmt.Fprintf(os.Stderr, "[DEBUG] Static file handler registered\n")

	// Start server
	port := ":8080"
	fmt.Fprintf(os.Stderr, "[DEBUG] Server starting on port %s\n", port)
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Printf("Available endpoints:\n")
	fmt.Printf("  POST /api/register        - Create a new account\n")
	fmt.Printf("  POST /api/login           - Login to existing account\n")
	fmt.Printf("  POST /api/logout          - Logout from account\n")
	fmt.Printf("  GET  /api/profile         - Get current user profile\n")
	fmt.Printf("  POST /api/change-password - Change user password\n")
	fmt.Printf("  POST /api/base64/encode   - Encode text to base64\n")
	fmt.Printf("  POST /api/base64/decode   - Decode base64 to text\n")
	fmt.Printf("  GET  /api/health          - Health check\n")
	fmt.Printf("\nServer running at http://localhost%s\n", port)

	fmt.Fprintf(os.Stderr, "[DEBUG] Server ready to accept connections\n")
	log.Fatal(http.ListenAndServe(port, router))
}
