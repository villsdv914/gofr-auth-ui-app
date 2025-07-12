# GoFr Auth UI App

A simple authentication application built with GoFr framework, featuring user registration, login, and JWT-based authentication.

## Features

- User registration and login
- JWT-based authentication
- PostgreSQL database
- Docker support
- Static file serving for UI
- Database migrations

## Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- PostgreSQL (if running locally)

## Quick Start

### Using Docker (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd gofr-auth-ui-app
```

2. Create environment file:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Run with Docker Compose:
```bash
make docker-up
```

The application will be available at `http://localhost:8000`

### Local Development

1. Set up environment:
```bash
make setup
```

2. Start PostgreSQL (if not using Docker):
```bash
# Install and start PostgreSQL
# Create database: mydb
```

3. Run the application:
```bash
make run
```

## Configuration

The application uses environment variables for configuration. Key variables:

### Database
- `DATABASE_HOST`: Database host (default: localhost)
- `DATABASE_PORT`: Database port (default: 5432)
- `DATABASE_USER`: Database user (default: postgres)
- `DATABASE_PASSWORD`: Database password (default: postgres)
- `DATABASE_NAME`: Database name (default: mydb)
- `DATABASE_SSLMODE`: SSL mode (default: disable)

### Application
- `APP_PORT`: Application port (default: 8000)
- `APP_HOST`: Application host (default: 0.0.0.0)

### JWT
- `JWT_SECRET`: JWT secret key
- `JWT_EXPIRY`: JWT expiry in seconds (default: 3600)

## API Endpoints

- `POST /signup` - User registration
- `POST /login` - User login
- `GET /me` - Get current user info (requires JWT)
- `GET /ui/*` - Static UI files

## Make Commands

- `make run` - Run the application locally
- `make build` - Build the application
- `make test` - Run tests
- `make docker-up` - Start with Docker Compose
- `make docker-down` - Stop Docker Compose
- `make docker-clean` - Clean Docker volumes
- `make migrate` - Run database migrations
- `make gen-migration name=<name>` - Generate new migration
- `make clean` - Clean build artifacts
- `make deps` - Install dependencies

## Project Structure

```
gofr-auth-ui-app/
├── cmd/
│   └── main.go          # Application entry point
├── configs/
│   ├── config.yaml      # Configuration file
│   └── .env            # Environment variables
├── handler/
│   └── auth.go         # HTTP handlers
├── migration/
│   └── users.go        # Database migrations
├── model/
│   └── user.go         # Data models
├── service/
│   └── auth.go         # Business logic
├── ui/                 # Static UI files
├── docker-compose.yml  # Docker services
├── Dockerfile         # Docker build
├── Makefile          # Build commands
└── README.md         # This file
```

## Database

The application uses PostgreSQL with automatic migrations. The users table is created automatically on startup.

## Troubleshooting

### Common Issues

1. **Database connection failed**: Check database credentials and ensure PostgreSQL is running
2. **Static files not found**: Ensure UI files are copied to the Docker container
3. **JWT authentication failed**: Check JWT_SECRET environment variable

### Logs

Check application logs for detailed error information:
```bash
docker-compose logs app
```

## License

This project is licensed under the MIT License. 