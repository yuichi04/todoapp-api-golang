# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **Go-based Todo API educational project** designed as a comprehensive learning resource for backend API development. It implements Clean Architecture with Domain-Driven Design principles using **Go standard packages only** (no external frameworks) to provide deep understanding of Go fundamentals and API development concepts.

**Critical Architecture Decision**: This project deliberately avoids frameworks like Gin, GORM, etc., using only Go standard packages (`net/http`, `database/sql`, `encoding/json`) to teach fundamental concepts without abstraction layers.

## Development Commands

### Hot Reload Development (Recommended)
```bash
# Start development server with hot reload
make dev-hot

# Alternative: Direct Air usage (if installed)
air -c .air.toml
```

### Standard Development
```bash
# Quick setup and run
make setup run

# Build and run
go build -o todoapp ./cmd/api
./todoapp

# Run directly without build
go run ./cmd/api/main.go
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
go test ./internal/infrastructure/database/
```

### Code Quality
```bash
# Format and lint
make lint          # Runs both fmt and vet
make fmt          # go fmt ./...
make vet          # go vet ./...
```

### Docker Development
```bash
# Complete Docker setup and start
make docker-setup docker-start

# Individual Docker commands
make docker-logs       # View logs
make docker-stop       # Stop services
make docker-restart    # Restart services
make docker-clean      # Cleanup containers
make docker-reset      # Complete reset with data deletion
```

## Architecture Overview

### Clean Architecture Implementation
```
cmd/api/main.go                    # Entry point with manual dependency injection
internal/
├── domain/                        # Business logic (core)
│   ├── entity/todo.go            # Business entities
│   ├── repository/               # Data access interfaces
│   └── service/                  # Business logic services
├── application/                  # Application services
│   ├── handler/                  # HTTP handlers (standard net/http)
│   ├── middleware/               # Custom middleware chain
│   └── dto/                      # Data transfer objects
└── infrastructure/               # External concerns
    ├── database/                 # Database implementations
    │   ├── connection.go         # Connection management
    │   └── todo_repository_impl.go
    └── web/                      # HTTP server
        ├── server.go             # Standard HTTP server
        └── routes.go             # Manual routing logic
```

### Key Architectural Patterns
- **Manual Dependency Injection**: Explicit construction in `main.go`, no DI frameworks
- **Repository Pattern**: Interfaces in domain, implementations in infrastructure
- **Standard HTTP Handling**: Custom routing with `net/http.ServeMux` and manual path parsing
- **Custom Middleware Chain**: Hand-built middleware system without frameworks
- **Standard Package Focus**: `database/sql`, `net/http`, `encoding/json` only

### Data Flow
1. HTTP Request → Custom Router (`infrastructure/web/routes.go`)
2. Router → Middleware Chain → Handler (`application/handler/`)
3. Handler → DTO Conversion → Domain Service (`domain/service/`)
4. Service → Repository Interface → Database Implementation (`infrastructure/database/`)
5. Response flows back through layers with appropriate transformations

## Standard Packages Philosophy

This project emphasizes Go standard library usage for educational purposes:
- **`net/http`** instead of Gin/Echo - Learn HTTP fundamentals
- **`database/sql`** instead of GORM - Understand SQL and database operations
- **Custom middleware** - Learn middleware patterns and HTTP request lifecycle
- **Manual routing** - Understand URL pattern matching and path parsing
- **`encoding/json`** - Learn JSON serialization/deserialization

## Database Architecture

- **Supported Drivers**: MySQL (`github.com/go-sql-driver/mysql`), SQLite (`github.com/mattn/go-sqlite3`)
- **Connection Management**: Custom pooling in `infrastructure/database/connection.go`
- **Query Approach**: Raw SQL with prepared statements for security and performance
- **Schema Management**: Automatic table creation in development, manual migration in production
- **Testing Strategy**: SQLite in-memory databases for integration tests

## Configuration System

Environment-based configuration through `pkg/config/`:
- **Database**: `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_DRIVER`
- **Server**: `SERVER_HOST`, `SERVER_PORT`
- **Application**: `APP_ENV` for environment-specific behavior
- **Defaults**: Development-friendly defaults in code

## Development Tools

### Hot Reload Setup
- **Air Configuration**: `.air.toml` for file watching and auto-restart
- **Development Scripts**: Cross-platform scripts in `scripts/` directory
- **Makefile Integration**: `make dev-hot` for easy startup

### Testing Strategy
- **Unit Tests**: Domain logic with mock repositories
- **Integration Tests**: Database operations with test containers/SQLite
- **HTTP Tests**: Handler testing with `net/http/httptest`
- **Table-Driven Tests**: Standard Go testing patterns throughout

## Code Organization Principles

### Dependency Rules
- **Domain Layer**: No dependencies on external layers
- **Application Layer**: Depends only on domain interfaces
- **Infrastructure Layer**: Implements domain interfaces, handles external systems

### Interface Usage
- Repository interfaces defined in `domain/repository/`
- Service interfaces for testing in `domain/service/`
- All external dependencies accessed through interfaces

### Error Handling Patterns
- Domain errors bubble up through layers
- HTTP status codes determined at handler level
- Database errors wrapped with context
- Graceful degradation where appropriate

## Common Development Tasks

### Adding New Entities
1. Create entity in `internal/domain/entity/`
2. Define repository interface in `internal/domain/repository/`
3. Implement repository in `internal/infrastructure/database/`
4. Create service in `internal/domain/service/`
5. Add DTOs in `internal/application/dto/`
6. Implement handlers in `internal/application/handler/`
7. Update routing in `internal/infrastructure/web/routes.go`

### Database Schema Changes
- Development: Modify `CreateTables()` in database manager
- Production: Create migration scripts (not included in this educational project)

### Adding Middleware
- Implement in `internal/application/middleware/`
- Add to chain in `web.NewRouter()`
- Follow the `func(http.Handler) http.Handler` pattern

## Educational Focus Areas

This codebase is designed to teach:
- **Clean Architecture implementation** without frameworks
- **Standard Go HTTP server** development
- **Database/SQL package** usage and patterns  
- **Dependency injection** without external libraries
- **Error handling** and logging best practices
- **Testing strategies** for different architectural layers
- **JSON handling** and API design
- **Middleware patterns** and HTTP request lifecycle

The verbose comments and educational structure make this ideal for learning Go backend development fundamentals.