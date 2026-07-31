# ExpenseSync — Smart Expense Tracker API

ExpenseSync is a RESTful API for managing personal expenses, built with Go and the Echo web framework. It uses an in-memory storage layer with thread-safe data access, supports category filtering and title search, calculates overall and monthly totals, and provides interactive OpenAPI/Swagger documentation.

---

## Features

- **CRUD Operations**: Create, read, list, and delete expenses.
- **Financial Precision**: All monetary values are represented as `int64` cents (e.g. `2550` for `$25.50`) to eliminate floating-point rounding errors.
- **Filtering & Search**: Filter expenses by exact category match and case-insensitive title search.
- **Calculations & Analytics**: Calculate overall totals, category totals, and monthly summaries grouped by month and category.
- **OpenAPI / Swagger Documentation**: Interactive API documentation generated with `swag`.
- **Live-Reloading & Local Development**: `make run` integrates with Go Air for automatic live-reloading.
- **Docker Support**: Multi-stage `Dockerfile` and `docker-compose.yml` for containerized execution.

---

## Structure

```
ExpenseSync/
├── AGENTS.md
├── README.md            
├── AI_NOTES.md           
├── Makefile              
├── Dockerfile            
├── docker-compose.yml   
├── .air.toml             
├── go.mod
├── go.sum
├── main.go               
├── src/                  
│   ├── models.go         
│   ├── storage.go        
│   ├── handlers.go       
│   └── router.go         
├── tests/                
│   └── expense_test.go
├── api-docs/             
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
└── docs/                 
    └── implementation.md
```

---

## Prerequisites

- **Go**: Version 1.25 or higher
- **Make**: (Optional, for running `make` targets)
- **Docker & Docker Compose**: (Optional, for containerized execution)

---

## Installation & Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Amazing-Stardom/ExpenseSync.git
   cd ExpenseSync
   ```

2. **Download dependencies**:
   ```bash
   go mod download
   ```

---

## Running the Server

### Option A: Using Makefile
```bash
make run
```
*Note: `make run` automatically uses `air` for live-reloading if installed; otherwise, it executes `go run main.go`.*

### Option B: Using Go CLI Directly
```bash
go run main.go
```

The server starts on port `8080` (or the port specified by the `PORT` environment variable).

---

## Running Tests

### Option A: Using Makefile
```bash
make test
```

### Option B: Using Go Testing Tool Directly
```bash
go test ./... -v
```

To run tests with code coverage:
```bash
make cover
```

---

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Server health check endpoint |
| `POST` | `/api/v1/expenses` | Add a new expense record |
| `GET` | `/api/v1/expenses` | List all expenses (supports `?category=` and `?search=`) |
| `GET` | `/api/v1/expenses/totals` | Calculate overall total expense amount (in cents) |
| `GET` | `/api/v1/expenses/category/:category/total` | Calculate total expense amount for a category (in cents) |
| `GET` | `/api/v1/summary/monthly` | Get monthly expense summaries grouped by month & category |
| `GET` | `/api/v1/expenses/:id` | Retrieve a specific expense by ID |
| `DELETE` | `/api/v1/expenses/:id` | Delete an expense by ID |

---

## Interactive API Documentation (Swagger)

Once the server is running, open your web browser to:

- **Swagger UI**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **Convenience Redirects**: [http://localhost:8080/api](http://localhost:8080/api) or [http://localhost:8080/docs](http://localhost:8080/docs)

To regenerate Swagger documentation after code changes:
```bash
make swagger
```

---

## Running with Docker

### Using Docker Compose
```bash
make docker-up
# Or directly:
docker-compose up --build
```

To stop the container:
```bash
make docker-down
```

### Using Docker CLI directly
```bash
make docker-build
make docker-run
```

---

## License

This project is open-source and available under the [MIT License](LICENSE).