package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles all authentication-related operations
type AuthHandler struct {
	users    map[string]*User
	sessions *sessions.CookieStore
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(secretKey []byte) *AuthHandler {
	return &AuthHandler{
		users:    make(map[string]*User),
		sessions: sessions.NewCookieStore(secretKey),
	}
}

// RegisterHandler handles user registration
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[DEBUG] Registration request received\n")

	if r.Method != http.MethodPost {
		fmt.Fprintf(os.Stderr, "[DEBUG] Invalid method: %s\n", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Failed to decode request body: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] Missing required fields\n")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Username, email, and password are required",
		})
		return
	}

	// Check if user already exists by username
	for _, existingUser := range h.users {
		if existingUser.Username == req.Username {
			fmt.Fprintf(os.Stderr, "[DEBUG] User already exists: %s\n", req.Username)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Username already exists",
			})
			return
		}
		if existingUser.Email == req.Email {
			fmt.Fprintf(os.Stderr, "[DEBUG] Email already exists: %s\n", req.Email)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Email already exists",
			})
			return
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Failed to hash password: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	// Create user
	user := &User{
		ID:       generateID(),
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Created:  time.Now(),
	}

	h.users[user.ID] = user

	// Return user data (without password)
	response := Response{
		Success: true,
		Message: "User registered successfully. Please login with your credentials.",
		Data:    map[string]string{"username": user.Username},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Fprintf(os.Stderr, "[DEBUG] User registered successfully: %s\n", user.Username)
}

// LoginHandler handles user login
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[DEBUG] Login request received\n")

	if r.Method != http.MethodPost {
		fmt.Fprintf(os.Stderr, "[DEBUG] Invalid method: %s\n", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Failed to decode request body: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Find user by username
	var user *User
	for _, u := range h.users {
		if u.Username == req.Username {
			user = u
			break
		}
	}

	if user == nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] User not found: %s\n", req.Username)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Invalid password for user: %s\n", req.Username)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// Create session
	session, _ := h.sessions.Get(r, "user-session")
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	// Return user data (without password)
	userResponse := User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Created:  user.Created,
	}

	response := Response{
		Success: true,
		Message: "Login successful",
		Data:    userResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Fprintf(os.Stderr, "[DEBUG] User logged in successfully: %s\n", user.Username)
}

// LogoutHandler handles user logout
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[DEBUG] Logout request received\n")

	if r.Method != http.MethodPost {
		fmt.Fprintf(os.Stderr, "[DEBUG] Invalid method: %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Clear session
	session, _ := h.sessions.Get(r, "user-session")
	session.Values["user_id"] = ""
	session.Options.MaxAge = -1
	session.Save(r, w)

	response := Response{
		Success: true,
		Message: "Logged out successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Fprintf(os.Stderr, "[DEBUG] User logged out successfully\n")
}

// ProfileHandler returns user profile information
func (h *AuthHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[DEBUG] Profile request received\n")

	if r.Method != http.MethodGet {
		fmt.Fprintf(os.Stderr, "[DEBUG] Invalid method: %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session
	session, err := h.sessions.Get(r, "user-session")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Failed to get session: %v\n", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] No valid user ID in session\n")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Find user
	user, exists := h.users[userID]
	if !exists {
		fmt.Fprintf(os.Stderr, "[DEBUG] User not found: %s\n", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Return user data (without password)
	userResponse := User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Created:  user.Created,
	}

	response := Response{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    userResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Fprintf(os.Stderr, "[DEBUG] Profile retrieved for user: %s\n", user.Username)
}

// ChangePasswordHandler handles password changes
func (h *AuthHandler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[DEBUG] Password change request received\n")

	if r.Method != http.MethodPost {
		fmt.Fprintf(os.Stderr, "[DEBUG] Invalid method: %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session
	session, err := h.sessions.Get(r, "user-session")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Failed to get session: %v\n", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] No valid user ID in session\n")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.CurrentPassword == "" || req.NewPassword == "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] Missing required fields\n")
		http.Error(w, "Current and new password are required", http.StatusBadRequest)
		return
	}

	if len(req.NewPassword) < 6 {
		fmt.Fprintf(os.Stderr, "[DEBUG] New password too short\n")
		http.Error(w, "New password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	// Find user
	user, exists := h.users[userID]
	if !exists {
		fmt.Fprintf(os.Stderr, "[DEBUG] User not found: %s\n", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Invalid current password for user: %s\n", user.Username)
		http.Error(w, "Invalid current password", http.StatusUnauthorized)
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] Failed to hash new password: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Update user password
	user.Password = string(hashedPassword)
	h.users[userID] = user

	response := Response{
		Success: true,
		Message: "Password changed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Fprintf(os.Stderr, "[DEBUG] Password changed successfully for user: %s\n", user.Username)
}

// Helper function to generate unique IDs
func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
