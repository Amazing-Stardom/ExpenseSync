// Package expense contains the core domain types, storage layer, and HTTP handlers
// for the ExpenseSync REST API.
package expense

import "time"

// Expense represents a single personal expense record.
// Amount is stored as int64 cents (e.g., 2550 for $25.50) to prevent floating-point precision risks.
type Expense struct {
	ID       int       `json:"id"       example:"1"`
	Title    string    `json:"title"    example:"Lunch at restaurant"`
	Amount   int64     `json:"amount"   example:"2550"`
	Category string    `json:"category" example:"food"`
	Date     time.Time `json:"date"     example:"2026-08-01T00:00:00Z"`
}

// CreateExpenseRequest is the validated request body for POST /api/v1/expenses.
// Amount is stored as int64 cents (e.g., 2550 for $25.50).
type CreateExpenseRequest struct {
	Title    string    `json:"title"    validate:"required"`
	Amount   int64     `json:"amount"   validate:"required,gt=0"`
	Category string    `json:"category" validate:"required"`
	Date     time.Time `json:"date"     validate:"required"`
}

// CategoryTotal is the response body for GET /api/v1/expenses/category/:cat/total.
// Total is stored as int64 cents.
type CategoryTotal struct {
	Category string `json:"category" example:"food"`
	Total    int64  `json:"total"    example:"15075"`
}

// MonthlySummary groups expenses by month.
// Total and ByCategory values are stored as int64 cents to ensure precision.
type MonthlySummary struct {
	Month      string           `json:"month"       example:"2026-08"`
	Total      int64            `json:"total"       example:"50000"`
	ByCategory map[string]int64 `json:"by_category"`
}

// ErrorResponse is the standard error body returned by all handlers.
type ErrorResponse struct {
	Message string `json:"message" example:"Invalid input"`
	Code    int    `json:"code"    example:"400"`
}
