# KnowledgeHub API

A knowledge management application with RESTful API built using Go.

## Prerequisites

- Go 1.23+
- PostgreSQL
- Git

## Getting Started

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/KnowledgeHub.git
   cd KnowledgeHub
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Set up environment variables:
   ```
   cp .env.example .env
   ```
   Edit the `.env` file with your configuration.

### Running the Application

1. Start the server:
   ```
   go run cmd/main.go
   ```

2. The API will be available at `http://localhost:8080`

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

### Updating Swagger Documentation
