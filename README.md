# experiment1: LLM slop for fun and learning

(ignore this entire repo and please don't use anything)

use `air` to monitor for changes and re-run - `go install github.com/air-verse/air@latest`

# Go Authentication Server

A simple HTTP server built in Go that provides user authentication functionality including account creation, login, logout, and profile management.

## Features

- **User Registration**: Create new accounts with username, email, and password
- **User Login**: Authenticate existing users
- **User Logout**: Securely end user sessions
- **Profile Management**: View user profile information
- **Session Management**: Secure cookie-based sessions
- **Password Security**: Bcrypt password hashing
- **RESTful API**: Clean HTTP endpoints
- **Web Interface**: Simple HTML frontend for testing

## Project Structure

```
.
├── main.go           # Main Go server code
├── go.mod            # Go module dependencies
├── static/           # Static web files
│   └── index.html    # Web interface
└── README.md         # This file
```

## Prerequisites

- Go 1.21 or later
- Git (for cloning)

## Installation

1. **Clone or download the project files**

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the server**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication Endpoints

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| `POST` | `/api/register` | Create a new user account | `{"username": "string", "email": "string", "password": "string"}` |
| `POST` | `/api/login` | Login to existing account | `{"username": "string", "password": "string"}` |
| `POST` | `/api/logout` | Logout from account | None |
| `GET` | `/api/profile` | Get current user profile | None |
| `GET` | `/api/health` | Health check endpoint | None |

### Request/Response Format

All API endpoints return JSON responses in the following format:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // Response data (varies by endpoint)
  }
}
```

## Usage Examples

### Using cURL

1. **Register a new user**
   ```bash
   curl -X POST http://localhost:8080/api/register \
     -H "Content-Type: application/json" \
     -d '{"username":"john_doe","email":"john@example.com","password":"secret123"}'
   ```

2. **Login**
   ```bash
   curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"username":"john_doe","password":"secret123"}' \
     -c cookies.txt
   ```

3. **Get profile (using cookies from login)**
   ```bash
   curl -X GET http://localhost:8080/api/profile \
     -b cookies.txt
   ```

4. **Logout**
   ```bash
   curl -X POST http://localhost:8080/api/logout \
     -b cookies.txt
   ```

### Using the Web Interface

1. Open your browser and navigate to `http://localhost:8080`
2. Use the forms to register, login, and manage your account
3. The interface will automatically handle sessions and show your profile

## Security Features

- **Password Hashing**: All passwords are hashed using bcrypt with cost factor 14
- **Session Management**: Secure cookie-based sessions with configurable secret keys
- **Input Validation**: Server-side validation of all user inputs
- **Conflict Prevention**: Username and email uniqueness enforcement

## Configuration

### Session Secret

The session secret key is currently hardcoded in the `main.go` file. **For production use, change this to a secure random string:**

```go
sessions: sessions.NewCookieStore([]byte("your-secure-random-secret-key")),
```

### Server Port

The server runs on port 8080 by default. To change this, modify the `port` variable in the `main()` function:

```go
port := ":3000" // Change to your desired port
```

## Dependencies

- **github.com/gorilla/mux**: HTTP router and URL matcher
- **github.com/gorilla/sessions**: Session management
- **golang.org/x/crypto/bcrypt**: Password hashing

## Development

### Adding New Features

The code is structured to make it easy to add new features:

1. **New Handlers**: Add new handler methods to the `Server` struct
2. **New Routes**: Register new routes in the `main()` function
3. **New Data Types**: Define new structs for requests/responses

### Database Integration

Currently, user data is stored in memory. To add persistence:

1. Replace the `map[string]*User` with a database connection
2. Implement database operations in the handler methods
3. Add proper error handling for database operations

## Testing

### Manual Testing

1. Start the server: `go run main.go`
2. Use the web interface at `http://localhost:8080`
3. Test all authentication flows

### API Testing

Use tools like:
- **cURL**: Command-line HTTP client
- **Postman**: GUI-based API testing
- **Insomnia**: Modern API client

## Production Considerations

Before deploying to production:

1. **Change the session secret key**
2. **Use HTTPS** with proper SSL certificates
3. **Implement rate limiting** to prevent abuse
4. **Add logging** for security monitoring
5. **Use environment variables** for configuration
6. **Implement proper error handling** and logging
7. **Add database persistence** instead of in-memory storage
8. **Set up monitoring** and health checks

## Troubleshooting

### Common Issues

1. **Port already in use**: Change the port number in `main.go`
2. **Dependencies not found**: Run `go mod tidy`
3. **Session not working**: Check browser cookie settings

### Debug Mode

To add debug logging, you can modify the handlers to include `fmt.Printf` statements or use Go's `log` package.

## License

This project is open source and available under the MIT License.

## Contributing

Feel free to submit issues, feature requests, or pull requests to improve this project.
