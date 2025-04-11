# Event App with Gin and MongoDB

[![Test](https://github.com/tanvirrb/event-app-gin/actions/workflows/test.yml/badge.svg)](https://github.com/tanvirrb/event-app-gin/actions/workflows/test.yml)

A RESTful API for managing events built with Go, Gin framework, and MongoDB.

## Features

- CRUD operations for events
- MongoDB database integration
- Docker support for development and testing
- Hot reloading for development
- End-to-end testing
- Linting and formatting with golangci-lint

## Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- Task (optional, for running predefined tasks)
- Air (optional, for hot reloading)
- golangci-lint (optional, for linting and formatting)

## Getting Started

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/tanvirrb/event-app-gin.git
   cd event-app-gin
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

### Running the Application

#### Using Docker (Recommended for Development)

1. Start the development environment:
   ```bash
   task docker-dev
   ```

2. Stop the development environment:
   ```bash
   task docker-dev-down
   ```

#### Running Locally

1. Set up environment variables:
   ```bash
   export PORT=3001
   export MONGODB_URI=mongodb://localhost:27017
   export GIN_MODE=debug
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

### Events API

#### Create Event
- **URL**: `/events`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "Tech Conference 2023",
    "genre": "Technology"
  }
  ```
- **Response**:
  ```json
  {
    "data": {
      "_id": "67f8efef9095da7e5d4c74f8",
      "name": "Tech Conference 2023",
      "genre": "Technology"
    }
  }
  ```
- **Curl Request**:
  ```bash
  curl -X POST http://localhost:3001/events \
    -H "Content-Type: application/json" \
    -d '{"name": "Tech Conference 2023", "genre": "Technology"}'
  ```

#### Get Event
- **URL**: `/events/:id`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "data": {
      "_id": "67f8ecfaea7f1d2363aed711",
      "name": "Tech Conference 2023",
      "genre": "Technology"
    }
  }
  ```
- **Curl Request**:
  ```bash
  curl -X GET http://localhost:3001/events/67f8ecfaea7f1d2363aed711
  ```

#### Get All Events
- **URL**: `/events`
- **Method**: `GET`
- **Response**:
  ```json
  {
    "data": [
      {
        "_id": "67f8ecfaea7f1d2363aed711",
        "name": "Tech Conference 2023",
        "genre": "Technology"
      },
      {
        "_id": "67f8ff8c04a6f3f0b569635b",
        "name": "Summer Music Festival",
        "genre": "Music"
      },
      {
        "_id": "67f8ff8c04a6f3f0b569635c",
        "name": "Food & Wine Expo",
        "genre": "Culinary"
      }
    ]
  }
  ```
- **Curl Request**:
  ```bash
  curl -X GET http://localhost:3001/events
  ```

#### Update Event
- **URL**: `/events/:id`
- **Method**: `PUT`
- **Request Body**:
  ```json
  {
    "name": "Tech Conference 2024",
    "genre": "Technology"
  }
  ```
- **Response**:
  ```json
  {
    "data": {
      "_id": "67f8ecfaea7f1d2363aed711",
      "name": "Tech Conference 2024",
      "genre": "Technology"
    }
  }
  ```
- **Curl Request**:
  ```bash
  curl -X PUT http://localhost:3001/events/67f8ecfaea7f1d2363aed711 \
    -H "Content-Type: application/json" \
    -d '{"name": "Tech Conference 2024", "genre": "Technology"}'
  ```

#### Delete Event
- **URL**: `/events/:id`
- **Method**: `DELETE`
- **Response**:
  ```json
  {
    "data": "Event deleted successfully with id: 67f8ecfaea7f1d2363aed711"
  }
  ```
- **Curl Request**:
  ```bash
  curl -X DELETE http://localhost:3001/events/67f8ecfaea7f1d2363aed711
  ```

## Development

### Available Tasks

- `task lint`: Run linter
- `task format`: Format code
- `task docker-dev`: Start development environment
- `task docker-dev-down`: Stop development environment
- `task docker-test`: Run end-to-end tests
- `task tidy`: Run `go mod tidy`

### Project Structure

```
.
├── src/
│   ├── bootstrap/     # Application bootstrap
│   ├── configs/       # Configuration
│   ├── events/        # Event domain
│   ├── router/        # Route definitions
│   └── tests/         # Tests
├── .air.toml          # Air hot reload configuration
├── .golangci.yml      # Linter configuration
├── docker-compose.dev.yml  # Development Docker Compose
├── docker-compose.test.yml # Test Docker Compose
├── Dockerfile.dev     # Development Dockerfile
├── Taskfile.yml       # Task definitions
└── main.go            # Application entry point
```

## Testing

### Running Tests

```bash
task docker-test
```

This will:
1. Start a test environment with MongoDB
2. Run end-to-end tests
3. Clean up the test environment

## Troubleshooting

### Common Issues

#### Hot Reloading Not Working

If changes to your code aren't being picked up by the hot reloader:

1. Make sure you're using the development environment:
   ```bash
   task docker-dev
   ```

2. Check that the volume mount is working correctly:
   ```bash
   docker-compose -f docker-compose.dev.yml exec app ls -la /app
   ```

3. Restart the development environment:
   ```bash
   task docker-dev-down
   task docker-dev
   ```

#### MongoDB Connection Issues

If the application can't connect to MongoDB:

1. Check if MongoDB is running:
   ```bash
   docker-compose -f docker-compose.dev.yml ps
   ```

2. Verify the MongoDB URI in your environment:
   ```bash
   docker-compose -f docker-compose.dev.yml exec app env | grep MONGODB_URI
   ```

3. Restart the MongoDB container:
   ```bash
   docker-compose -f docker-compose.dev.yml restart mongodb
   ```

## License

[MIT License](LICENSE)

