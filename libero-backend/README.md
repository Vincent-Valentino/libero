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