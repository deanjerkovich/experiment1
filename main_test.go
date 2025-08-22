package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	server := NewServer()

	if server.users == nil {
		t.Error("Expected users map to be initialized")
	}

	if server.sessions == nil {
		t.Error("Expected sessions to be initialized")
	}

	if len(server.users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(server.users))
	}
}

func TestGenerateID(t *testing.T) {
	id1 := generateID()
	id2 := generateID()

	if id1 == "" {
		t.Error("Expected non-empty ID")
	}

	if id2 == "" {
		t.Error("Expected non-empty ID")
	}

	if id1 == id2 {
		t.Error("Expected unique IDs")
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := hashPassword(password)
	if err != nil {
		t.Errorf("Expected no error hashing password: %v", err)
	}

	if hash == "" {
		t.Error("Expected non-empty hash")
	}

	if hash == password {
		t.Error("Expected hash to be different from original password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := hashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password for test: %v", err)
	}

	// Test correct password
	if !checkPassword(password, hash) {
		t.Error("Expected correct password to match hash")
	}

	// Test incorrect password
	if checkPassword("wrongpassword", hash) {
		t.Error("Expected incorrect password to not match hash")
	}

	// Test empty password
	if checkPassword("", hash) {
		t.Error("Expected empty password to not match hash")
	}
}

func TestRegisterHandler(t *testing.T) {
	server := NewServer()

	tests := []struct {
		name            string
		request         RegisterRequest
		expectedStatus  int
		expectedSuccess bool
	}{
		{
			name: "Valid registration",
			request: RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus:  http.StatusCreated,
			expectedSuccess: true,
		},
		{
			name: "Missing username",
			request: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
		},
		{
			name: "Missing email",
			request: RegisterRequest{
				Username: "testuser",
				Password: "password123",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
		},
		{
			name: "Missing password",
			request: RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
		},
		{
			name: "Password too short",
			request: RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "123",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			server.registerHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedSuccess {
				var response Response
				json.Unmarshal(w.Body.Bytes(), &response)

				if !response.Success {
					t.Errorf("Expected success response, got: %v", response)
				}

				if len(server.users) != 1 {
					t.Errorf("Expected 1 user, got %d", len(server.users))
				}
			}
		})
	}
}

func TestRegisterHandlerDuplicateUser(t *testing.T) {
	server := NewServer()

	// Register first user
	user1 := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	body1, _ := json.Marshal(user1)
	req1 := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	server.registerHandler(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Errorf("Expected first registration to succeed, got status %d", w1.Code)
	}

	// Try to register duplicate username
	user2 := RegisterRequest{
		Username: "testuser",
		Email:    "different@example.com",
		Password: "password456",
	}

	body2, _ := json.Marshal(user2)
	req2 := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	server.registerHandler(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("Expected duplicate username to fail, got status %d", w2.Code)
	}

	// Try to register duplicate email
	user3 := RegisterRequest{
		Username: "differentuser",
		Email:    "test@example.com",
		Password: "password789",
	}

	body3, _ := json.Marshal(user3)
	req3 := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body3))
	req3.Header.Set("Content-Type", "application/json")

	w3 := httptest.NewRecorder()
	server.registerHandler(w3, req3)

	if w3.Code != http.StatusConflict {
		t.Errorf("Expected duplicate email to fail, got status %d", w3.Code)
	}

	// Should still only have 1 user
	if len(server.users) != 1 {
		t.Errorf("Expected 1 user after duplicate attempts, got %d", len(server.users))
	}
}

func TestLoginHandler(t *testing.T) {
	server := NewServer()

	// First register a user
	registerReq := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	registerBody, _ := json.Marshal(registerReq)
	registerHTTPReq := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(registerBody))
	registerHTTPReq.Header.Set("Content-Type", "application/json")

	registerW := httptest.NewRecorder()
	server.registerHandler(registerW, registerHTTPReq)

	tests := []struct {
		name            string
		request         LoginRequest
		expectedStatus  int
		expectedSuccess bool
	}{
		{
			name: "Valid login",
			request: LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
		},
		{
			name: "Invalid username",
			request: LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			expectedStatus:  http.StatusUnauthorized,
			expectedSuccess: false,
		},
		{
			name: "Invalid password",
			request: LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			expectedStatus:  http.StatusUnauthorized,
			expectedSuccess: false,
		},
		{
			name: "Missing username",
			request: LoginRequest{
				Password: "password123",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
		},
		{
			name: "Missing password",
			request: LoginRequest{
				Username: "testuser",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			server.loginHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedSuccess {
				var response Response
				json.Unmarshal(w.Body.Bytes(), &response)

				if !response.Success {
					t.Errorf("Expected success response, got: %v", response)
				}

				// Check that user data is returned
				if response.Data == nil {
					t.Error("Expected user data in response")
				}
			}
		})
	}
}

func TestProfileHandler(t *testing.T) {
	server := NewServer()

	// First register and login a user to get a session
	registerReq := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	registerBody, _ := json.Marshal(registerReq)
	registerHTTPReq := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(registerBody))
	registerHTTPReq.Header.Set("Content-Type", "application/json")

	registerW := httptest.NewRecorder()
	server.registerHandler(registerW, registerHTTPReq)

	// Login to create a session
	loginReq := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	loginBody, _ := json.Marshal(loginReq)
	loginHTTPReq := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(loginBody))
	loginHTTPReq.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	server.loginHandler(loginW, loginHTTPReq)

	// Extract cookies from login response
	cookies := loginW.Result().Cookies()

	tests := []struct {
		name            string
		cookies         []*http.Cookie
		expectedStatus  int
		expectedSuccess bool
	}{
		{
			name:            "Valid session",
			cookies:         cookies,
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
		},
		{
			name:            "No cookies",
			cookies:         []*http.Cookie{},
			expectedStatus:  http.StatusUnauthorized,
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/profile", nil)

			// Add cookies if they exist
			for _, cookie := range tt.cookies {
				req.AddCookie(cookie)
			}

			w := httptest.NewRecorder()
			server.profileHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedSuccess {
				var response Response
				json.Unmarshal(w.Body.Bytes(), &response)

				if !response.Success {
					t.Errorf("Expected success response, got: %v", response)
				}

				if response.Data == nil {
					t.Error("Expected user data in response")
				}
			}
		})
	}
}

func TestLogoutHandler(t *testing.T) {
	server := NewServer()

	// First register and login a user to get a session
	registerReq := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	registerBody, _ := json.Marshal(registerReq)
	registerHTTPReq := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(registerBody))
	registerHTTPReq.Header.Set("Content-Type", "application/json")

	registerW := httptest.NewRecorder()
	server.registerHandler(registerW, registerHTTPReq)

	// Login to create a session
	loginReq := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	loginBody, _ := json.Marshal(loginReq)
	loginHTTPReq := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(loginBody))
	loginHTTPReq.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	server.loginHandler(loginW, loginHTTPReq)

	// Extract cookies from login response
	cookies := loginW.Result().Cookies()

	// Test logout
	req := httptest.NewRequest("POST", "/api/logout", nil)

	// Add cookies
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	w := httptest.NewRecorder()
	server.logoutHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response Response
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response.Success {
		t.Errorf("Expected success response, got: %v", response)
	}
}

func TestChangePasswordHandler(t *testing.T) {
	server := NewServer()

	// First register and login a user to get a session
	registerReq := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	registerBody, _ := json.Marshal(registerReq)
	registerHTTPReq := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(registerBody))
	registerHTTPReq.Header.Set("Content-Type", "application/json")

	registerW := httptest.NewRecorder()
	server.registerHandler(registerW, registerHTTPReq)

	// Login to create a session
	loginReq := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	loginBody, _ := json.Marshal(loginReq)
	loginHTTPReq := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(loginBody))
	loginHTTPReq.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	server.loginHandler(loginW, loginHTTPReq)

	// Extract cookies from login response
	cookies := loginW.Result().Cookies()

	tests := []struct {
		name            string
		request         ChangePasswordRequest
		cookies         []*http.Cookie
		expectedStatus  int
		expectedSuccess bool
	}{
		{
			name: "Valid password change",
			request: ChangePasswordRequest{
				CurrentPassword: "password123",
				NewPassword:     "newpassword456",
			},
			cookies:         cookies,
			expectedStatus:  http.StatusOK,
			expectedSuccess: true,
		},
		{
			name: "Wrong current password",
			request: ChangePasswordRequest{
				CurrentPassword: "wrongpassword",
				NewPassword:     "newpassword456",
			},
			cookies:         cookies,
			expectedStatus:  http.StatusUnauthorized,
			expectedSuccess: false,
		},
		{
			name: "New password too short",
			request: ChangePasswordRequest{
				CurrentPassword: "password123",
				NewPassword:     "123",
			},
			cookies:         cookies,
			expectedStatus:  http.StatusBadRequest,
			expectedSuccess: false,
		},
		{
			name: "No session",
			request: ChangePasswordRequest{
				CurrentPassword: "password123",
				NewPassword:     "newpassword456",
			},
			cookies:         []*http.Cookie{},
			expectedStatus:  http.StatusUnauthorized,
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/change-password", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Add cookies if they exist
			for _, cookie := range tt.cookies {
				req.AddCookie(cookie)
			}

			w := httptest.NewRecorder()
			server.changePasswordHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedSuccess {
				var response Response
				json.Unmarshal(w.Body.Bytes(), &response)

				if !response.Success {
					t.Errorf("Expected success response, got: %v", response)
				}

				// Verify the password was actually changed
				if tt.name == "Valid password change" {
					// Try to login with new password
					newLoginReq := LoginRequest{
						Username: "testuser",
						Password: "newpassword456",
					}

					newLoginBody, _ := json.Marshal(newLoginReq)
					newLoginHTTPReq := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(newLoginBody))
					newLoginHTTPReq.Header.Set("Content-Type", "application/json")

					newLoginW := httptest.NewRecorder()
					server.loginHandler(newLoginW, newLoginHTTPReq)

					if newLoginW.Code != http.StatusOK {
						t.Error("Expected to be able to login with new password after change")
					}
				}
			}
		})
	}
}

func TestHealthHandler(t *testing.T) {
	server := NewServer()

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()

	server.healthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response Response
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response.Success {
		t.Errorf("Expected success response, got: %v", response)
	}

	if response.Data == nil {
		t.Error("Expected health data in response")
	}

	// Check that health data contains expected fields
	healthData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Error("Expected health data to be a map")
	}

	if _, exists := healthData["timestamp"]; !exists {
		t.Error("Expected timestamp in health data")
	}

	if _, exists := healthData["users"]; !exists {
		t.Error("Expected users count in health data")
	}
}

func TestUserStruct(t *testing.T) {
	user := User{
		ID:       "test-id",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Created:  time.Now(),
	}

	if user.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", user.ID)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected Email 'test@example.com', got '%s'", user.Email)
	}

	if user.Password != "hashedpassword" {
		t.Errorf("Expected Password 'hashedpassword', got '%s'", user.Password)
	}

	if user.Created.IsZero() {
		t.Error("Expected Created to be set")
	}
}

func TestResponseStruct(t *testing.T) {
	response := Response{
		Success: true,
		Message: "Test message",
		Data:    "test data",
	}

	if !response.Success {
		t.Error("Expected Success to be true")
	}

	if response.Message != "Test message" {
		t.Errorf("Expected Message 'Test message', got '%s'", response.Message)
	}

	if response.Data != "test data" {
		t.Errorf("Expected Data 'test data', got '%v'", response.Data)
	}
}
