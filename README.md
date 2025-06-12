# KnowledgeHub API

A knowledge management application with RESTful API built using Go.

## Prerequisites

- Go 1.23+ (for local development)
- PostgreSQL (for local development)
- Git
- Docker
- Docker Compose

## Getting Started

### Running with Docker

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/KnowledgeHub.git
   cd KnowledgeHub
   ```

2. Set up environment variables:
   ```
   cp .env.example .env
   ```
   Edit the `.env` file with your configuration.

3. Run the application using one of the following methods:

   a) Using batch file (Windows):
   ```
   ./start.bat up -d
   ```

   b) Using Docker Compose directly:
   ```
   cd deployments
   docker-compose up -d
   ```

4. The API will be available at `http://localhost:8080`

### Local Development

1. Install dependencies:
   ```
   go mod download
   ```

2. Start the server:
   ```
   go run cmd/main.go
   ```

## API Documentation with Swagger

### Setting up Swagger

1. Install Swagger CLI tool:
   ```
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. Generate Swagger documentation:
   ```
   swag init -g internal/controller/http/router.go
   ```

3. Make sure `SWAGGER_ENABLED=true` is set in your `.env` file.

4. Access Swagger UI at `http://localhost:8080/swagger/index.html`
