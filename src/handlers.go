package expense

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Handler holds the ExpenseStore dependency for HTTP handlers.
type Handler struct {
	store *ExpenseStore
}

// NewHandler creates a new Handler instance with the given ExpenseStore.
func NewHandler(store *ExpenseStore) *Handler {
	return &Handler{store: store}
}

// CreateExpense handles POST /api/v1/expenses
// @Summary Create new expense
// @Description Add a new expense to the tracker
// @Tags expenses
// @Accept json
// @Produce json
// @Param expense body CreateExpenseRequest true "Expense details"
// @Success 201 {object} Expense
// @Failure 400 {object} ErrorResponse
// @Router /expenses [post]
func (h *Handler) CreateExpense(c echo.Context) error {
	var req CreateExpenseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid request body",
			Code:    http.StatusBadRequest,
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Validation failed: invalid input parameters",
			Code:    http.StatusBadRequest,
		})
	}

	exp := Expense{
		Title:    req.Title,
		Amount:   req.Amount,
		Category: req.Category,
		Date:     req.Date,
	}

	created := h.store.Create(exp)
	return c.JSON(http.StatusCreated, created)
}

// GetAllExpenses handles GET /api/v1/expenses
// @Summary Get all expenses
// @Description Retrieve all expenses with optional filtering by category or title search
// @Tags expenses
// @Produce json
// @Param category query string false "Filter by category"
// @Param search query string false "Search by title"
// @Success 200 {array} Expense
// @Router /expenses [get]
func (h *Handler) GetAllExpenses(c echo.Context) error {
	category := c.QueryParam("category")
	search := c.QueryParam("search")

	var expenses []Expense
	if category != "" || search != "" {
		expenses = h.store.Filter(category, search)
	} else {
		expenses = h.store.GetAll()
	}

	if expenses == nil {
		expenses = make([]Expense, 0)
	}
	return c.JSON(http.StatusOK, expenses)
}

// GetExpenseByID handles GET /api/v1/expenses/:id
// @Summary Get expense by ID
// @Description Retrieve a specific expense by its unique ID
// @Tags expenses
// @Produce json
// @Param id path int true "Expense ID"
// @Success 200 {object} Expense
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /expenses/{id} [get]
func (h *Handler) GetExpenseByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid ID parameter",
			Code:    http.StatusBadRequest,
		})
	}

	exp, found := h.store.GetByID(id)
	if !found {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Message: "Expense not found",
			Code:    http.StatusNotFound,
		})
	}
	return c.JSON(http.StatusOK, exp)
}

// DeleteExpense handles DELETE /api/v1/expenses/:id
// @Summary Delete expense
// @Description Remove an expense by its unique ID
// @Tags expenses
// @Param id path int true "Expense ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /expenses/{id} [delete]
func (h *Handler) DeleteExpense(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid ID parameter",
			Code:    http.StatusBadRequest,
		})
	}

	deleted := h.store.Delete(id)
	if !deleted {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Message: "Expense not found",
			Code:    http.StatusNotFound,
		})
	}
	return c.NoContent(http.StatusNoContent)
}

// GetTotalExpenses handles GET /api/v1/expenses/totals
// @Summary Get total expenses
// @Description Calculate overall total of all expenses
// @Tags calculations
// @Produce json
// @Success 200 {object} map[string]int64
// @Router /expenses/totals [get]
func (h *Handler) GetTotalExpenses(c echo.Context) error {
	total := h.store.GetTotal()
	return c.JSON(http.StatusOK, map[string]int64{"total": total})
}

// GetCategoryTotal handles GET /api/v1/expenses/category/:category/total
// @Summary Get category total
// @Description Calculate total expenses for a specific category
// @Tags calculations
// @Produce json
// @Param category path string true "Category name"
// @Success 200 {object} CategoryTotal
// @Router /expenses/category/{category}/total [get]
func (h *Handler) GetCategoryTotal(c echo.Context) error {
	category := c.Param("category")
	total := h.store.GetTotalByCategory(category)
	return c.JSON(http.StatusOK, CategoryTotal{
		Category: category,
		Total:    total,
	})
}

// GetMonthlySummary handles GET /api/v1/summary/monthly
// @Summary Get monthly summary
// @Description Get total expenses grouped by month and category
// @Tags analytics
// @Produce json
// @Success 200 {array} MonthlySummary
// @Router /summary/monthly [get]
func (h *Handler) GetMonthlySummary(c echo.Context) error {
	// GetMonthlySummaries returns []MonthlySummary.
	// The Swagger annotation {array} MonthlySummary is correct: the HTTP response is a JSON array.
	summaries := h.store.GetMonthlySummaries()
	if summaries == nil {
		summaries = make([]MonthlySummary, 0)
	}
	return c.JSON(http.StatusOK, summaries)
}
