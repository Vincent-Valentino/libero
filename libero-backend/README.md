# Libero Backend

This is the backend service for the Libero application.

## Requirements

- Go 1.23 or later
- PostgreSQL

## Setup

1. Clone the repository
2. Install dependencies:
   ```
   go mod download
   ```
3. Set up environment variables in `.env` file (see `.env.example` for reference)
4. Run the application:
   ```
   go run main.go
   ```

## Development with Hot Reload

For development, you can use hot reloading to automatically restart the server when code changes are detected.

### Prerequisites

Install the Air tool for hot reloading:

```bash
go install github.com/air-verse/air@latest
```

Make sure the `$GOPATH/bin` directory is in your PATH.

### Running with Hot Reload

#### Windows

Run the dev.bat script:

```
dev.bat
```

Or directly run:

```
air
```

#### Unix/Linux/macOS

Run the dev.sh script:

```
./dev.sh
```

Or directly run:

```
air
```

## Project Structure

- `main.go`: Application entry point
- `app.go`: Application initialization
- `config/`: Configuration settings
- `internal/`:
  - `api/`:
    - `controllers/`: HTTP request handlers
    - `dto/`: Data Transfer Objects
    - `routes/`: API route definitions
  - `middleware/`: HTTP middleware
  - `models/`: Database models
  - `repository/`: Data access layer
  - `service/`: Business logic layer

## Code Architecture

The application follows a clean architecture with clear separation of concerns:

1. **Controllers**: Handle HTTP requests and responses
2. **Services**: Implement business logic
3. **Repositories**: Manage data access
4. **Models**: Define data structures
5. **DTOs**: Define data transfer objects

Each component is defined by interfaces, allowing for loose coupling and easier testing.

## Features
- **User Authentication**: Register, login, password reset and change (`/auth/register`, `/auth/login`, `/auth/password/*`).
- **OAuth2 Integration**: Social login with Google, Facebook, GitHub (`/auth/*/login`, `/auth/*/callback`).
- **Sports Data API**: Fetch upcoming matches, results, player stats, fixtures summary (`/api/matches/*`, `/api/players/{id}/stats`).
- **User Profile & Preferences**: Retrieve and update preferences (`GET /api/users/profile`, `PUT /api/users/preferences`).
- **Background Tasks**:
  - **Cache Cleanup**: Runs every 15 minutes to purge expired entries.
  - **Fixtures Scheduler**: Refreshes fixtures data every 4 hours.

## Data Flow & Request Lifecycle
1. **Router Layer**: `routes.SetupRoutes` registers public and protected routes using Gorilla Mux.
2. **Middleware**: CORS handling and JWT authentication (`AuthMiddleware`).
3. **Controllers**: Parse HTTP requests, validate input DTOs, and invoke service methods.
4. **Services**: Contain business logic and orchestrate calls to repositories and external APIs.
5. **Repositories**: Abstract database interactions and caching logic (e.g., Redis or in-memory cache).
6. **Response**: Controllers format service results into JSON responses with appropriate HTTP status codes.

## Coding Style & Conventions
- **Clean Architecture**: Separation into `controllers`, `service`, `repository`, `models`, `dto` packages under `internal/`.
- **Dependency Injection**: Pass `Repository`, `Service`, and `Config` structs to constructors for testability.
- **Routing & Middleware**: Use Gorilla Mux for flexible routing patterns and middleware chaining.
- **Configuration**: Centralize in `config` package, loading from environment variables via `.env`.
- **Logging**: Use Go's standard `log` package; structured logs recommended for production.
- **Error Handling**: Wrap and return errors with context; controllers translate errors to HTTP responses.
- **Project Layout**: Follow Go conventions (lowercase package names, `internal` for private code).
- **Testing**: Write unit tests for services and controllers; mock repositories for isolation. 