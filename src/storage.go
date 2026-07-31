package expense

import (
	"strings"
	"sync"
)

// ExpenseStore provides thread-safe in-memory storage for Expense items.
type ExpenseStore struct {
	mu       sync.RWMutex
	expenses map[int]Expense
	nextID   int
}

// NewExpenseStore creates a new ExpenseStore initialized with an empty map and nextID set to 1.
func NewExpenseStore() *ExpenseStore {
	return &ExpenseStore{
		expenses: make(map[int]Expense),
		nextID:   1,
	}
}

// Create assigns an auto-incrementing ID to the expense and saves it in the store.
func (s *ExpenseStore) Create(exp Expense) Expense {
	s.mu.Lock()
	defer s.mu.Unlock()

	exp.ID = s.nextID
	s.expenses[exp.ID] = exp
	s.nextID++
	return exp
}

// GetAll returns all expenses as a slice. It returns an empty slice []Expense{}, never nil.
func (s *ExpenseStore) GetAll() []Expense {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]Expense, 0, len(s.expenses))
	for _, exp := range s.expenses {
		results = append(results, exp)
	}
	return results
}

// GetByID retrieves an expense by its unique ID. Returns the expense and true if found, or zero value and false.
func (s *ExpenseStore) GetByID(id int) (Expense, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	exp, ok := s.expenses[id]
	return exp, ok
}

// Delete removes an expense by its unique ID. Returns true if deleted, false if not found.
func (s *ExpenseStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.expenses[id]; ok {
		delete(s.expenses, id)
		return true
	}
	return false
}

// Filter searches expenses by category (exact match if non-empty) and search title (case-insensitive substring match).
// It returns an empty slice []Expense{}, never nil.
func (s *ExpenseStore) Filter(category, search string) []Expense {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]Expense, 0)
	for _, exp := range s.expenses {
		if category != "" && exp.Category != category {
			continue
		}
		if search != "" && !contains(exp.Title, search) {
			continue
		}
		results = append(results, exp)
	}
	return results
}

// GetTotal calculates the sum of amounts for all stored expenses in int64 cents.
func (s *ExpenseStore) GetTotal() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var total int64
	for _, exp := range s.expenses {
		total += exp.Amount
	}
	return total
}

// GetTotalByCategory calculates the sum of amounts for expenses in the specified category in int64 cents.
func (s *ExpenseStore) GetTotalByCategory(category string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var total int64
	for _, exp := range s.expenses {
		if exp.Category == category {
			total += exp.Amount
		}
	}
	return total
}

// GetMonthlySummaries groups expenses by YYYY-MM month and calculates overall total and category totals in int64 cents.
func (s *ExpenseStore) GetMonthlySummaries() []MonthlySummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	summaryMap := make(map[string]*MonthlySummary)
	for _, exp := range s.expenses {
		monthStr := exp.Date.Format("2006-01")
		summary, exists := summaryMap[monthStr]
		if !exists {
			summary = &MonthlySummary{
				Month:      monthStr,
				Total:      0,
				ByCategory: make(map[string]int64),
			}
			summaryMap[monthStr] = summary
		}
		summary.Total += exp.Amount
		summary.ByCategory[exp.Category] += exp.Amount
	}

	results := make([]MonthlySummary, 0, len(summaryMap))
	for _, summary := range summaryMap {
		results = append(results, *summary)
	}
	return results
}

// contains returns true if s contains substr as a case-insensitive substring.
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
