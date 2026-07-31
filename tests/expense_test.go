package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	expense "github.com/Amazing-Stardom/ExpenseSync/src"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	v *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.v.Struct(i)
}

var fixedDate = time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)

func TestStorageCreateAndGet(t *testing.T) {
	store := expense.NewExpenseStore()
	exp := expense.Expense{
		Title:    "Lunch",
		Amount:   2550,
		Category: "food",
		Date:     fixedDate,
	}

	created := store.Create(exp)
	if created.ID != 1 {
		t.Fatalf("Expected ID 1, got %d", created.ID)
	}

	got, ok := store.GetByID(1)
	if !ok {
		t.Fatalf("Expected to find expense with ID 1")
	}
	if got.Title != "Lunch" {
		t.Errorf("Expected title Lunch, got %s", got.Title)
	}
}

func TestStorageDelete(t *testing.T) {
	store := expense.NewExpenseStore()
	exp := store.Create(expense.Expense{
		Title:    "Movie",
		Amount:   1500,
		Category: "entertainment",
		Date:     fixedDate,
	})

	deleted := store.Delete(exp.ID)
	if !deleted {
		t.Fatalf("Expected Delete to return true for existing item")
	}

	_, ok := store.GetByID(exp.ID)
	if ok {
		t.Fatalf("Expense should not exist after deletion")
	}

	if store.Delete(999) != false {
		t.Fatalf("Expected Delete to return false for non-existent ID")
	}
}

func TestStorageFilterAndTotals(t *testing.T) {
	store := expense.NewExpenseStore()
	store.Create(expense.Expense{Title: "Pizza Lunch", Amount: 2000, Category: "food", Date: fixedDate})
	store.Create(expense.Expense{Title: "Bus Ticket", Amount: 500, Category: "transport", Date: fixedDate})

	foodItems := store.Filter("food", "")
	if len(foodItems) != 1 {
		t.Errorf("Expected 1 food item, got %d", len(foodItems))
	}

	pizzaItems := store.Filter("", "PIZZA")
	if len(pizzaItems) != 1 {
		t.Errorf("Expected case-insensitive search for PIZZA to find 1 item, got %d", len(pizzaItems))
	}

	total := store.GetTotal()
	if total != 2500 {
		t.Errorf("Expected total 2500 cents, got %d", total)
	}

	foodTotal := store.GetTotalByCategory("food")
	if foodTotal != 2000 {
		t.Errorf("Expected food total 2000 cents, got %d", foodTotal)
	}
}

func TestStorageMonthlySummaries(t *testing.T) {
	store := expense.NewExpenseStore()
	store.Create(expense.Expense{Title: "Groceries", Amount: 5000, Category: "food", Date: fixedDate})

	summaries := store.GetMonthlySummaries()
	if len(summaries) != 1 {
		t.Fatalf("Expected 1 monthly summary, got %d", len(summaries))
	}
	if summaries[0].Month != "2026-08" {
		t.Errorf("Expected month 2026-08, got %s", summaries[0].Month)
	}
	if summaries[0].ByCategory["food"] != 5000 {
		t.Errorf("Expected food total 5000 cents in monthly summary, got %d", summaries[0].ByCategory["food"])
	}
}

func TestHTTPCreateExpense(t *testing.T) {
	store := expense.NewExpenseStore()
	h := expense.NewHandler(store)

	e := echo.New()
	e.Validator = &CustomValidator{v: validator.New()}
	e.POST("/api/v1/expenses", h.CreateExpense)

	// Happy path
	reqBody, _ := json.Marshal(expense.CreateExpenseRequest{
		Title:    "Coffee",
		Amount:   450,
		Category: "food",
		Date:     fixedDate,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/expenses", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d", rec.Code)
	}

	var res expense.Expense
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if res.ID != 1 || res.Title != "Coffee" {
		t.Errorf("Unexpected created expense data: %+v", res)
	}

	// Error path: missing title
	badBody := []byte(`{"amount": 450, "category": "food"}`)
	reqBad := httptest.NewRequest(http.MethodPost, "/api/v1/expenses", bytes.NewReader(badBody))
	reqBad.Header.Set("Content-Type", "application/json")
	recBad := httptest.NewRecorder()

	e.ServeHTTP(recBad, reqBad)

	if recBad.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 Bad Request, got %d", recBad.Code)
	}
}

func TestHTTPGetAndDeleteExpense(t *testing.T) {
	store := expense.NewExpenseStore()
	h := expense.NewHandler(store)

	e := echo.New()
	e.GET("/api/v1/expenses", h.GetAllExpenses)
	e.GET("/api/v1/expenses/:id", h.GetExpenseByID)
	e.DELETE("/api/v1/expenses/:id", h.DeleteExpense)

	// Pre-populate store
	store.Create(expense.Expense{Title: "Book", Amount: 1500, Category: "education", Date: fixedDate})

	// Get all
	req := httptest.NewRequest(http.MethodGet, "/api/v1/expenses", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", rec.Code)
	}

	var list []expense.Expense
	json.Unmarshal(rec.Body.Bytes(), &list)
	if len(list) != 1 {
		t.Errorf("Expected 1 item in list, got %d", len(list))
	}

	// Get by ID success
	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/expenses/1", nil)
	recGet := httptest.NewRecorder()
	e.ServeHTTP(recGet, reqGet)
	if recGet.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK for GET /expenses/1, got %d", recGet.Code)
	}

	// Get by ID 404
	req404 := httptest.NewRequest(http.MethodGet, "/api/v1/expenses/999", nil)
	rec404 := httptest.NewRecorder()
	e.ServeHTTP(rec404, req404)
	if rec404.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found, got %d", rec404.Code)
	}

	// Delete success
	reqDel := httptest.NewRequest(http.MethodDelete, "/api/v1/expenses/1", nil)
	recDel := httptest.NewRecorder()
	e.ServeHTTP(recDel, reqDel)
	if recDel.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 No Content for DELETE /expenses/1, got %d", recDel.Code)
	}

	// Delete 404
	reqDel404 := httptest.NewRequest(http.MethodDelete, "/api/v1/expenses/1", nil)
	recDel404 := httptest.NewRecorder()
	e.ServeHTTP(recDel404, reqDel404)
	if recDel404.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found for second DELETE, got %d", recDel404.Code)
	}
}

func TestHTTPTotalsAndSummaries(t *testing.T) {
	store := expense.NewExpenseStore()
	h := expense.NewHandler(store)

	e := echo.New()
	e.GET("/api/v1/expenses/totals", h.GetTotalExpenses)
	e.GET("/api/v1/summary/monthly", h.GetMonthlySummary)

	store.Create(expense.Expense{Title: "Item A", Amount: 1000, Category: "cat1", Date: fixedDate})

	// Totals
	reqTotal := httptest.NewRequest(http.MethodGet, "/api/v1/expenses/totals", nil)
	recTotal := httptest.NewRecorder()
	e.ServeHTTP(recTotal, reqTotal)

	if recTotal.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for totals, got %d", recTotal.Code)
	}

	var totalMap map[string]int64
	json.Unmarshal(recTotal.Body.Bytes(), &totalMap)
	if totalMap["total"] != 1000 {
		t.Errorf("Expected total 1000 cents, got %d", totalMap["total"])
	}

	// Monthly summary
	reqSummary := httptest.NewRequest(http.MethodGet, "/api/v1/summary/monthly", nil)
	recSummary := httptest.NewRecorder()
	e.ServeHTTP(recSummary, reqSummary)

	if recSummary.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for monthly summary, got %d", recSummary.Code)
	}

	var summaries []expense.MonthlySummary
	if err := json.Unmarshal(recSummary.Body.Bytes(), &summaries); err != nil {
		t.Fatalf("Failed to unmarshal monthly summary response: %v", err)
	}
	if len(summaries) != 1 {
		t.Fatalf("Expected 1 monthly summary entry, got %d", len(summaries))
	}
	if summaries[0].Month != "2026-08" {
		t.Errorf("Expected month 2026-08 in summary, got %s", summaries[0].Month)
	}
	if summaries[0].Total != 1000 {
		t.Errorf("Expected monthly total 1000 cents, got %d", summaries[0].Total)
	}
	if summaries[0].ByCategory["cat1"] != 1000 {
		t.Errorf("Expected cat1 category total 1000 cents, got %d", summaries[0].ByCategory["cat1"])
	}
}
