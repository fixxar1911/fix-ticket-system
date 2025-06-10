# Ticket System

A RESTful API-based ticketing system built with Go, Gin, and PostgreSQL.

## Features

- Create, read, update, and delete tickets
- Track ticket status and priority
- Assign tickets to team members
- RESTful API endpoints
- PostgreSQL database integration

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

## Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/fix-ticket-system.git
cd fix-ticket-system
```

2. Install dependencies:
```bash
go mod download
```

3. Create a PostgreSQL database:
```sql
CREATE DATABASE ticket_system;
```

4. Configure environment variables:
Create a `.env` file in the root directory with the following content:
```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ticket_system
DB_PORT=5432
```

5. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`.

## API Endpoints

### Tickets

- `POST /api/v1/tickets` - Create a new ticket
- `GET /api/v1/tickets` - Get all tickets
- `GET /api/v1/tickets/:id` - Get a specific ticket
- `PUT /api/v1/tickets/:id` - Update a ticket
- `DELETE /api/v1/tickets/:id` - Delete a ticket

### Example Request

Create a new ticket:
```bash
curl -X POST http://localhost:8080/api/v1/tickets \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Bug in login page",
    "description": "Users cannot log in after password reset",
    "created_by": "john.doe@example.com"
  }'
```

## Project Structure

```
.
├── config/         # Configuration files
├── models/         # Data models
├── repository/     # Database operations
├── service/        # Business logic
├── main.go         # Application entry point
├── go.mod          # Go module file
└── README.md       # This file
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Unit Tests and Code Coverage

The project now includes comprehensive unit tests for all ticket endpoints, including both success and error cases. The tests use the `testify` package for mocking and assertions.

### Running Tests

To run the tests, execute the following command from the project root:

```bash
go test ./...
```

To generate a coverage report, run:

```bash
go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out
```

The current code coverage is approximately 77.5%, with detailed coverage for models, repository, service, and HTTP handlers.

### Test Structure

- **Models**: Tests for the `NewTicket` function.
- **Repository**: Tests for CRUD operations using an in-memory SQLite database.
- **Service**: Tests for business logic, including error handling.
- **HTTP Handlers**: Tests for all ticket endpoints, including error cases.

### Mocking

The tests use `testify/mock` to mock the `TicketService` interface, allowing for isolated testing of HTTP handlers without a real database connection. 