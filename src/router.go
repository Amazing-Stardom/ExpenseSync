package expense

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// RegisterRoutes configures all middleware, health check, Swagger UI, and API v1 endpoints on the provided Echo instance.
func RegisterRoutes(e *echo.Echo, h *Handler) {
	// Configure restricted CORS origins to ensure secure cross-origin resource sharing.
	allowedOrigins := []string{"http://localhost:8080", "http://localhost:3000", "http://127.0.0.1:8080"}
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		allowedOrigins = strings.Split(origins, ",")
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	// Swagger UI endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Convenience redirect routes for API documentation
	e.GET("/api", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	e.GET("/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API v1 routes
	// Note: Specific static routes (/expenses/totals, /expenses/category/...) must be registered
	// BEFORE the parameterized route (/expenses/:id) to prevent "totals" or "category" from matching :id.
	v1 := e.Group("/api/v1")
	v1.POST("/expenses", h.CreateExpense)
	v1.GET("/expenses", h.GetAllExpenses)
	v1.GET("/expenses/totals", h.GetTotalExpenses)
	v1.GET("/expenses/category/:category/total", h.GetCategoryTotal)
	v1.GET("/summary/monthly", h.GetMonthlySummary)

	v1.GET("/expenses/:id", h.GetExpenseByID)
	v1.DELETE("/expenses/:id", h.DeleteExpense)
}
