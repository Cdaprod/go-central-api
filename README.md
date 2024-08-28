# Go Central API

This is a central API designed to work as an API Gateway for various services including Repocate and MinIO. It's built with Go, implements a registry pattern for dynamic API management, and can be easily deployed and accessed through Tailscale.

## Features

- API Gateway pattern for unified access to multiple services
- Registry pattern for dynamic API management
- Modular design with loose coupling and high cohesion
- Easy configuration through JSON files
- Middleware for logging and authentication
- Extensible routing system
- Docker support for easy deployment

## Getting Started

```plaintext
go-central-api/
│
├── main.go
├── go.mod
├── go.sum
├── config.json
├── Dockerfile
├── README.md
│
├── config/
│   └── config.go
│
├── handlers/
│   └── handlers.go
│
├── middleware/
│   └── middleware.go
│
├── registry/
│   └── registry.go
│
├── services/
│   ├── repocate/
│   │   └── repocate.go
│   ├── minio/
│   │   └── minio.go
│   └── ... (other services)
│
├── utils/
│   └── utils.go
│
└── tests/
    ├── handlers_test.go
    ├── middleware_test.go
    └── registry_test.go
```

1. Clone the repository:

   ```
   git clone https://github.com/Cdaprod/go-central-api.git
   ```

2. Navigate to the project directory:

   ```
   cd go-central-api
   ```

3. Create a `config.json` file in the root directory:
 
   ```json
   {
     "server_address": ":8080",
     "database_url": "your_database_url",
     "jwt_secret": "your_jwt_secret"
   }
   ```

4. Build and run the application:
 
   ```
   go build
   ./go-central-api
   ```

## Docker Deployment

To build and run the Docker container:

```
docker build -t go-central-api .
docker run -p 8080:8080 go-central-api
```

## API Endpoints

- `GET /api/health`: Health check endpoint
- `GET/POST/PUT/DELETE /api/{service}/{path}`: Proxy requests to registered services

## Adding New Services

To add a new service to the API Gateway:

1. Create a new API type in `handlers/handlers.go`
2. Implement the `API` interface for your new type
3. Register the new API in `main.go`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.