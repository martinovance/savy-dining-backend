# Savy Dining Backend (Go)

High-performance reservation engine for the Savy Dining application.

## Tech Stack
- **Language**: Go 1.21+
- **Framework**: Gin Gonic
- **Features**: RESTful API, UUID generation, Validation

## Getting Started

1. Clone the repository
2. Run `go mod tidy`
3. Start the server: `go run main.go`

## API Endpoints
- `GET /health` - Service health status
- `POST /api/v1/reservations` - Create a new booking
- `GET /api/v1/reservations` - List all bookings
