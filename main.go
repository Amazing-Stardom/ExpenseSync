package main

import (
	"log"
	"os"
	"time"

	expense "github.com/Amazing-Stardom/ExpenseSync/src"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	_ "github.com/Amazing-Stardom/ExpenseSync/api-docs"
)

// CustomValidator wraps go-playground/validator/v10 for Echo request validation.
type CustomValidator struct {
	validator *validator.Validate
}

// Validate executes struct tag validations.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// @title           ExpenseSync API
// @version         1.0.0
// @description     A smart REST API for managing personal expenses
// @termsOfService  http://swagger.io/terms/

// @contact.name   Support
// @contact.email  support@expensesync.io

// @license.name  MIT

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	// ExpenseStore uses sync.RWMutex internally and is completely safe for concurrent access across HTTP requests.
	store := expense.NewExpenseStore()

	// Optional sample data seeding for manual testing environment
	if os.Getenv("SEED_SAMPLE_DATA") == "true" {
		seedDate := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
		store.Create(expense.Expense{Title: "Groceries", Amount: 4550, Category: "food", Date: seedDate.AddDate(0, 0, -5)})
		store.Create(expense.Expense{Title: "Taxi fare", Amount: 1200, Category: "transport", Date: seedDate.AddDate(0, 0, -3)})
	}

	h := expense.NewHandler(store)

	e := echo.New()

	// Register custom validator so c.Validate() processes validate: struct tags.
	e.Validator = &CustomValidator{validator: validator.New()}

	// Register middleware and routes using expense.RegisterRoutes
	expense.RegisterRoutes(e, h)

	// Read port from environment variable with fallback to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting ExpenseSync server on port %s", port)
	log.Fatal(e.Start(":" + port))
}
