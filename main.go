package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
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
	users    map[string]*User
	sessions *sessions.CookieStore
	mutex    sync.RWMutex
}

// NewServer creates a new server instance
func NewServer() *Server {
	return &Server{
		users:    make(map[string]*User),
		sessions: sessions.NewCookieStore([]byte("your-secret-key-change-this-in-production")),
		mutex:    sync.RWMutex{},
	}
}

// generateID generates a random ID
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// checkPassword checks if a password matches a hash
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// registerHandler handles user registration
func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[DEBUG] Registration request received\n")

	if r.Method != http.MethodPost {
		fmt.Fprintf(os.Stderr, "[DEBUG] Invalid method: %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] Registration attempt for username: %s, email: %s\n", req.Username, req.Email)

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] Missing required fields\n")
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		fmt.Fprintf(os.Stderr, "[DEBUG] Password too short: %d characters\n", len(req.Password))
		http.Error(w, "Password must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	fmt.Fprintf(os.Stderr, "[DEBUG] Checking for existing users with same username/email\n")
	fmt.Fprintf(os.Stderr, "[DEBUG] Total users in system: %d\n", len(s.users))

	// Check if username already exists
	for _, user := range s.users {
		if user.Username == req.Username {
			fmt.Fprintf(os.Stderr, "[DEBUG] Username already exists: %s\n", req.Username)
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		if user.Email == req.Email {
			fmt.Fprintf(os.Stderr, "[DEBUG] Email already exists: %s\n", req.Email)
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] No conflicts found, proceeding with user creation\n")

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Failed to hash password: %v\n", err)
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Create user
	user := &User{
		ID:       generateID(),
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Created:  time.Now(),
	}

	s.users[user.ID] = user
	fmt.Fprintf(os.Stderr, "[DEBUG] User created successfully: %s (ID: %s)\n", user.Username, user.ID)

	// Return user data (without password)
	response := Response{
		Success: true,
		Message: "User registered successfully. Please login with your credentials.",
		Data:    map[string]string{"username": user.Username},
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] Sending successful registration response for user: %s\n", user.Username)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// loginHandler handles user login
func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[DEBUG] Login request received\n")

	if r.Method != http.MethodPost {
		fmt.Fprintf(os.Stderr, "[DEBUG] Invalid method: %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] Login attempt for username: %s\n", req.Username)

	// Validate input
	if req.Username == "" || req.Password == "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] Missing username or password\n")
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	fmt.Fprintf(os.Stderr, "[DEBUG] Searching for user with username: %s\n", req.Username)
	fmt.Fprintf(os.Stderr, "[DEBUG] Total users in system: %d\n", len(s.users))

	// Find user by username
	var user *User
	for _, u := range s.users {
		if u.Username == req.Username {
			user = u
			fmt.Fprintf(os.Stderr, "[DEBUG] Found user: %s (ID: %s)\n", u.Username, u.ID)
			break
		}
	}

	if user == nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] User not found: %s\n", req.Username)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Check password
	fmt.Fprintf(os.Stderr, "[DEBUG] Verifying password for user: %s\n", user.Username)
	if !checkPassword(req.Password, user.Password) {
		fmt.Fprintf(os.Stderr, "[DEBUG] Password verification failed for user: %s\n", user.Username)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] Password verified successfully for user: %s\n", user.Username)

	// Create session
	session, _ := s.sessions.Get(r, "user-session")
	session.Values["user_id"] = user.ID
	session.Save(r, w)
	fmt.Fprintf(os.Stderr, "[DEBUG] Session created for user: %s (ID: %s)\n", user.Username, user.ID)

	response := Response{
		Success: true,
		Message: "Login successful",
		Data:    user,
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] Sending successful login response for user: %s\n", user.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// logoutHandler handles user logout
func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Clear session
	session, _ := s.sessions.Get(r, "user-session")
	session.Values["user_id"] = ""
	session.Options.MaxAge = -1
	session.Save(r, w)

	response := Response{
		Success: true,
		Message: "Logout successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// profileHandler returns the current user's profile
func (s *Server) profileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[DEBUG] Profile request received\n")

	if r.Method != http.MethodGet {
		fmt.Fprintf(os.Stderr, "[DEBUG] Invalid method: %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session
	session, _ := s.sessions.Get(r, "user-session")
	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] No valid session found\n")
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] Session found for user ID: %s\n", userID)

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[userID]
	if !exists {
		fmt.Fprintf(os.Stderr, "[DEBUG] User not found for ID: %s\n", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] Profile retrieved for user: %s (ID: %s)\n", user.Username, user.ID)

	response := Response{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// changePasswordHandler handles password changes
func (s *Server) changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session to identify user
	session, _ := s.sessions.Get(r, "user-session")
	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.CurrentPassword == "" || req.NewPassword == "" {
		http.Error(w, "Current password and new password are required", http.StatusBadRequest)
		return
	}

	if len(req.NewPassword) < 6 {
		http.Error(w, "New password must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Get user
	user, exists := s.users[userID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify current password
	if !checkPassword(req.CurrentPassword, user.Password) {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	// Hash new password
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "Error processing new password", http.StatusInternalServerError)
		return
	}

	// Update password
	user.Password = hashedPassword

	response := Response{
		Success: true,
		Message: "Password changed successfully",
	}

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
			"users":     len(s.users),
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
	fmt.Printf("  GET  /api/health          - Health check\n")
	fmt.Printf("\nServer running at http://localhost%s\n", port)

	fmt.Fprintf(os.Stderr, "[DEBUG] Server ready to accept connections\n")
	log.Fatal(http.ListenAndServe(port, router))
}
