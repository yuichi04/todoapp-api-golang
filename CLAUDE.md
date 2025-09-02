# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **Go-based Todo API educational project** designed as a learning resource for API development. It implements Clean Architecture with Domain-Driven Design principles using **Go standard packages only** (no frameworks like Gin or GORM) to help beginners understand fundamental Go concepts.

**Critical Architecture Decision**: This project deliberately uses Go standard packages (`net/http`, `database/sql`) instead of frameworks to provide deeper understanding of Go fundamentals.

## Commands

### Development (Hot Reload - Recommended)
```bash
# Start development server with hot reload
make dev-hot

# Alternative: Direct Air usage
air -c .air.toml

# Alternative: Development scripts
./scripts/dev.sh        # Linux/macOS
scripts\dev.bat         # Windows
```

### Standard Development
```bash
# Build and run
go build -o todoapp ./cmd/api
./todoapp

# Run directly
go run ./cmd/api

# Environment setup
make setup
```

### Testing
```bash
# Run all tests
go test ./...
make test

# Test with coverage
go test -cover ./...
make test-coverage

# Test specific package
go test ./internal/domain/service/
```

### Code Quality
```bash
# Format and lint
make lint
# Or individually:
make fmt    # go fmt ./...
make vet    # go vet ./...
```

### Docker Development
```bash
# Complete Docker setup
make docker-dev

# Individual Docker commands
make docker-setup      # Initial setup
make docker-start      # Start services
make docker-logs       # View logs
make docker-stop       # Stop services
make docker-clean      # Cleanup
```

## Architecture

### Clean Architecture Implementation
```
cmd/api/main.go                    # Entry point with manual DI
internal/
├── domain/                        # Business logic layer
│   ├── entity/todo.go            # Core business entities
│   ├── repository/               # Data access interfaces
│   └── service/                  # Business logic services
├── application/                  # Application layer
│   ├── handler/                  # HTTP handlers (standard net/http)
│   ├── middleware/               # Custom middleware implementations
│   └── dto/                      # Data transfer objects
└── infrastructure/               # Infrastructure layer
    ├── database/                 # database/sql implementations
    │   ├── connection.go         # DB connection management
    │   └── todo_repository_impl.go # Repository implementations
    └── web/                      # HTTP server setup
        ├── server.go             # Standard HTTP server
        └── routes.go             # Manual routing logic
```

### Key Architecture Patterns
- **Manual Dependency Injection**: No DI framework, explicit construction in `main.go`
- **Repository Pattern**: Interface in domain, implementation in infrastructure
- **Standard HTTP Routing**: Custom routing logic using `http.ServeMux`
- **Custom Middleware**: Hand-built middleware chain implementation

### Data Flow
1. HTTP Request → Custom Router (`web/routes.go`)
2. Router → Handler (`application/handler/`)
3. Handler → Domain Service (`domain/service/`)
4. Service → Repository Interface → Database Implementation
5. Response flows back through the layers

## Standard Packages Focus

This project emphasizes Go standard library usage:
- `net/http` instead of Gin/Echo for HTTP handling
- `database/sql` instead of GORM for database operations
- Custom middleware instead of framework middleware
- Manual routing instead of framework routers

## Database

- **Driver**: MySQL (`github.com/go-sql-driver/mysql` - only non-standard dependency)
- **Connection Management**: Custom connection pooling in `database/connection.go`
- **Queries**: Raw SQL with prepared statements for security
- **Schema**: Automatic table creation in development mode

## Configuration

Environment variables are loaded through `pkg/config/` with defaults for development:
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD` for database
- `SERVER_HOST`, `SERVER_PORT` for HTTP server
- `APP_ENV` for environment-specific behavior

## Hot Reload Setup (Air)

Development efficiency features:
- Air configuration in `.air.toml`
- Automatic rebuild on file changes
- Cross-platform development scripts
- Makefile integration for easy commands

## Testing Strategy

- Unit tests for domain logic
- Integration tests with database
- HTTP handler testing with `httptest` package
- Repository testing with test database