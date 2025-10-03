# Hackers Auth Service

A simple authentication service built with Go and Gin, providing JWT-based authentication for demo purposes.

## Features

- JWT-based authentication
- Swagger documentation
- Hardcoded test users with different access levels
- Docker support
- Comprehensive test suite

## Test Users

The service includes two hardcoded test users:

1. Beta User
   ```json
   {
     "username": "betauser",
     "password": "betauser",
     "company": "acme global",
     "beta_access": true
   }
   ```

2. Normal User
   ```json
   {
     "username": "normaluser",
     "password": "normaluser",
     "company": "generic co",
     "beta_access": false
   }
   ```

## API Documentation

The API documentation is available via Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

### Endpoints

#### POST /login
Authenticates a user and returns a JWT token along with user details.

Request body:
```json
{
  "username": "string",
  "password": "string"
}
```

Success Response (200):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "username": "string",
    "company": "string",
    "beta_access": boolean
  }
}
```

#### GET /users
Returns a list of all available users with their credentials. This endpoint is provided for demo purposes to easily populate the login form.

Success Response (200):
```json
[
  {
    "username": "betauser",
    "password": "betauser"
  },
  {
    "username": "normaluser",
    "password": "normaluser"
  }
]
```

## Running Locally

### Prerequisites
- Go 1.23 or higher
- Docker (optional)

### Without Docker

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Generate Swagger documentation:
   ```bash
   swag init
   ```
4. Run the application:
   ```bash
   go run main.go
   ```

### With Docker

1. Build the image:
   ```bash
   docker build -t hackers-auth .
   ```
2. Run the container:
   ```bash
   docker run -p 8080:8080 hackers-auth
   ```

## Helm Deployment

### Prerequisites
- Kubernetes cluster
- Helm 3.x
- cert-manager installed (for TLS)

### Installation

1. Create a custom values file (e.g., `my-values.yaml`):
   ```yaml
   # Image configuration
   image:
     repository: your-registry/hackers-auth  # Change to your image repository
     tag: your-tag                          # Change to your image tag
     pullPolicy: IfNotPresent

   # Domain configuration
   hostname: auth.your-domain.com           # Change to your domain

   # Resource limits
   resources:
     requests:
       cpu: 100m
       memory: 128Mi
     limits:
       cpu: 200m
       memory: 256Mi
   ```

2. Install the chart:
   ```bash
   helm install hackers-auth ./chart -f my-values.yaml
   ```

   Or upgrade an existing release:
   ```bash
   helm upgrade hackers-auth ./chart -f my-values.yaml
   ```

### Configuration Options

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.repository` | Image repository | `cloudbees-days/hackers-auth` |
| `image.tag` | Image tag | `latest` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `hostname` | Service hostname | `hackers-auth.example.com` |
| `resources.requests.cpu` | CPU request | `100m` |
| `resources.requests.memory` | Memory request | `128Mi` |
| `resources.limits.cpu` | CPU limit | `200m` |
| `resources.limits.memory` | Memory limit | `256Mi` |

## Testing

Run the test suite:
```bash
go test -v
```

The test suite includes:
- Login endpoint tests
- User authentication tests
- JWT token validation
- Error handling tests

## Development

### Project Structure
```
.
├── main.go          # Main application file
├── main_test.go     # Test suite
├── Dockerfile       # Docker configuration
├── .dockerignore    # Docker ignore file
├── .gitignore       # Git ignore file
├── chart/           # Helm chart directory
└── docs/           # Generated Swagger documentation (gitignored)
```

### Notes
- The JWT token expires after 24 hours
- Swagger documentation is generated during Docker build
- The service runs in release mode when deployed via Docker
