# Product Tracker API

A RESTful API for tracking products and their energy consumption, built with Go, Gin, and PostgreSQL.

## Features

- RESTful API endpoints for product management
- JWT-based authentication
- PostgreSQL database integration
- Swagger API documentation
- Configuration management with YAML support
- Health check endpoint

## Prerequisites

- Go 1.16 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/volsafe/product-tracker.git
    cd product-tracker
    ```

2. Install dependencies:

    ```sh
    go mod download
    ```

3. Set up the database:

    Ensure you have PostgreSQL installed and running. Create a database and a table with the following structure:

    ```sql
    CREATE DATABASE consumers;
    CREATE USER pgsql WITH PASSWORD 'pgsql';
    GRANT ALL PRIVILEGES ON DATABASE consumers TO pgsql;
    ```

4. Configure the application:
   - Copy `config/config.yaml.example` to `config/config.yaml`
   - Update the configuration values in `config/config.yaml`

## Configuration

The application can be configured through:
1. `config/config.yaml` file
2. Environment variables (override YAML settings)

### Configuration File Structure

```yaml
Server:
  Port: 8080
  ReadTimeout: 10s
  WriteTimeout: 10s
  IdleTimeout: 120s
  MaxHeaderBytes: 1048576

Database:
  Host: "localhost"
  Port: 5432
  User: "pgsql"
  Password: "pgsql"
  DbName: "consumers"
  SSLMode: "disable"

JWT:
  Secret: "your-secret-key"
  ExpirationTime: 24h
```

### Environment Variables

- `SERVER_PORT`: Server port (default: 8080)
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: pgsql)
- `DB_PASSWORD`: Database password (default: pgsql)
- `DB_NAME`: Database name (default: consumers)
- `DB_SSL_MODE`: Database SSL mode (default: disable)
- `JWT_SECRET`: JWT secret key

## Running the Application

1. Start the server:

    ```sh
    go run cmd/main.go
    ```

2. Access the API:
   - API Base URL: `http://localhost:8080/api/v1`
   - Swagger Documentation: `http://localhost:8080/swagger/index.html`
   - Health Check: `http://localhost:8080/health`

## API Endpoints

### Products

- `POST /api/v1/product/insert`: Import a new product
- `GET /api/v1/product/list`: List all products
- `GET /api/v1/product/list/{name}`: Get products by name

### Health Check

- `GET /health`: Check API health status

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-token>
```

## Development

### Project Structure

```
product-tracker/
├── cmd/
│   └── main.go           # Application entry point
├── config/
│   ├── config.go         # Configuration management
│   └── config.yaml       # Configuration file
├── controllers/
│   └── health.go         # Health check controller
├── db/
│   └── db.go            # Database connection management
├── handlers/
│   ├── health.go        # Health check handler
│   └── products.go      # Product handlers
├── models/
│   └── product.go       # Product model
├── routes/
│   └── routes.go        # Route definitions
├── storage/
│   └── storage.go       # Database operations
├── utils/
│   └── jwt.go          # JWT utilities
└── docs/               # Swagger documentation
```

### Adding New Features

1. Create new models in the `models` package
2. Add database operations in the `storage` package
3. Create handlers in the `handlers` package
4. Define routes in `routes/routes.go`
5. Update Swagger documentation

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
