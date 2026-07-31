# AI_NOTES.md — ExpenseSync

**Project**: Smart Expense Tracker API
**Stack**: Go + Echo Framework
**Repository**: https://github.com/Amazing-Stardom/ExpenseSync
**Date**: August 2026

---

## 1. What the AI Generated

The initial blueprint and all code scaffolding were AI-generated using Antigravity (Google DeepMind). The following sections were produced without manual authoring:

| File | Content | Source |
|---|---|---|
| `src/models.go` | All structs (`Expense`, `CreateExpenseRequest`, `MonthlySummary`, `CategoryTotal`, `ErrorResponse`) with JSON and Swagger tags | AI |
| `src/storage.go` | `ExpenseStore` with CRUD, filter, total, and monthly summary methods | AI |
| `src/handlers.go` | All HTTP handler functions with Swagger annotations | AI |
| `main.go` | Echo setup, middleware registration, route declarations, port configuration | AI |
| `tests/expense_test.go` | Unit and HTTP integration test cases | AI |
| `Dockerfile` | Multi-stage build using `golang:1.21-alpine` and `alpine:3.20` | AI |
| `docker-compose.yml` | Single-service compose file with `PORT` environment variable | AI |
| `Makefile` | Targets for run, test, cover, swagger, build, docker-build, docker-run, docker-up, docker-down | AI |
| `docs/implementation.md` | Phase-by-phase implementation plan with rules and test cases | AI |

---

## 2. What Was Validated, Critiqued, and Changed

Every item below represents a critic raised by the reviewer during the review process. Each entry records the problem, the fix, and why.

---

### 2.1 Critical Correctness Bugs

**Bug: `contains` function compared string lengths instead of checking for substring presence**

The AI-generated `contains(s, substr string) bool` function returned `len(s) > len(substr)`. This passes for any string longer than the search term, producing false positives.

Fix: replaced with `strings.Contains(strings.ToLower(s), strings.ToLower(substr))` for correct case-insensitive substring search.

---

**Bug: `GetMonthlySummaries` initialized `ByCategory` as `map[string]interface{}`**

The `MonthlySummary` struct declared `ByCategory map[string]float64`. The storage function created the map as `map[string]interface{}` and accumulated values with a `val.(float64)` type assertion. The assertion panics if the map ever holds a non-float64 value (for example, the first insertion which sets the value from `exp.Amount`, an untyped numeric literal).

Fix: changed the initialization to `make(map[string]float64)` and replaced the conditional accumulation block with `sum.ByCategory[exp.Category] += exp.Amount`. No type assertion is needed when the map type matches the value type.

---

**Bug: `json.Marshal` error was silently discarded in tests**

The HTTP handler test called `body, _ := json.Marshal(payload)`. A silent discard of the error means a failed marshal produces an empty body and a misleading test failure.

Fix: added `body, err := json.Marshal(payload)` followed by `if err != nil { t.Fatalf(...) }`.

---

### 2.2 Architecture Issues

**Issue: All handlers referenced a global variable `globalStore`**

The original handlers called `globalStore.Create(...)` directly. This couples all handler logic to a package-level variable, making unit tests unreliable (tests mutate shared state) and making it impossible to run tests in parallel.

Fix: introduced a `Handler` struct that holds `*ExpenseStore`. All handler functions became methods on `Handler`. The `main()` function creates the store, creates the handler, and passes handler methods to the router.

---

**Issue: `init()` function populated the store with sample data**

The `init()` function is called automatically before every test run. This means tests that relied on an empty store could receive pre-seeded data.

Fix: moved the seed data into `main()`, not `init()`. Tests create their own `ExpenseStore` instances and are never affected by seed data.

---

**Issue: Tests mutated `globalStore` directly**

The HTTP test set `globalStore = NewExpenseStore()` before each run. This is test-order-dependent and breaks parallel execution.

Fix: each test now creates a local `store := NewExpenseStore()` and a local `h := &Handler{store: store}`. No global variable is touched.

---

### 2.3 Validation Configuration

**Issue: `binding:` struct tags were used on `CreateExpenseRequest`**

Echo's default binder does not process `binding:` tags. Required-field validation silently passes regardless of input.

Fix: replaced all `binding:` tags with `validate:` tags. Added `go get github.com/go-playground/validator/v10` as a dependency. Registered a `CustomValidator` wrapping `go-playground/validator` on the Echo instance before the server starts. Added `c.Validate(&req)` after `c.Bind(&req)` in every handler that processes a request body.

---

### 2.4 Test Reliability

**Issue: `time.Now()` was used in test data**

Tests that depend on the current time can behave differently depending on when they run (timezone, daylight saving, test parallelism).

Fix: replaced all `time.Now()` calls in tests with a package-level constant `fixedDate = time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)`. All test fixtures use this value.

---

### 2.5 Infrastructure and Configuration

**Issue: Docker runtime image was `alpine:latest`**

Using `latest` produces non-reproducible builds. The image resolved by `latest` changes over time.

Fix: pinned the runtime stage to `alpine:3.20`.

---

**Issue: Port was hardcoded to `:8080` in `main()`**

Fix: port is read from `os.Getenv("PORT")` with a fallback to `8080`. The `docker-compose.yml` passes `PORT=8080` as an environment variable.

---

**Issue: `go get` was used to install the `swag` CLI**

`go get` modifies `go.mod` and installs a binary as a side effect. The correct command for installing CLI tools is `go install`.

Fix: setup section now uses `go install github.com/swaggo/swag/cmd/swag@latest`.

---

**Issue: Swagger-generated files were mixed into the `docs/` directory**

The `docs/` directory held both the hand-written implementation plan and the machine-generated OpenAPI artifacts (`docs.go`, `swagger.json`, `swagger.yaml`). This makes it harder to exclude generated files from code review and harder for tools (Redoc, Stoplight, GitHub Actions) to find the spec.

Fix: moved generated files to `api-docs/`. The `swag init` command is now `swag init --output ./api-docs --packageName docs`. The import path in `main.go` is `_ "github.com/Amazing-Stardom/ExpenseSync/api-docs"`.

---

**Issue: Module path was a placeholder (`github.com/yourusername/expense-sync`)**

The correct path was already known from the repository URL.

Fix: all module references updated to `github.com/Amazing-Stardom/ExpenseSync`.

---

**Issue: `AI_NOTES.md` was placed in `docs/`**

The assignment submission structure requires `AI_NOTES.md` at the repository root.

Fix: file moved to the repository root.

---

### 2.6 Swagger Annotation

**Issue: `GetMonthlySummary` Swagger annotation was flagged as a type mismatch**

`GetMonthlySummaries()` on the storage layer returns `map[string]MonthlySummary`. The Swagger annotation on the handler says `{array} MonthlySummary`. The reviewer queried whether these were inconsistent.

Clarification: the handler converts the map to a slice before writing the response (`for _, s := range summaries { result = append(result, s) }`). The HTTP response is a JSON array. The `{array} MonthlySummary` annotation is correct. A comment was added inside the handler to make the conversion intent explicit.

---

## 3. AI Suggestions Not Used

| Suggestion | Reason not used |
|---|---|
| Complex query builder for filtering | The assignment requires only category and title search. A query builder adds unnecessary abstraction for two parameters. |
| Custom error types (separate `AppError` struct) | A simple `ErrorResponse{Message, Code}` struct produces cleaner Swagger documentation and is sufficient for the scale of this API. |
| Database persistence (JSON file on disk) | The assignment explicitly permits in-memory storage. Adding file I/O introduces error paths (file not found, corrupt JSON) that are out of scope. |
| `map[string]interface{}` for `ByCategory` | The AI initially used this "flexible" type. It was rejected because it requires type assertions at every read site and conflicts with the struct declaration. `map[string]float64` is the correct type. |
| `init()` for sample data seeding | Rejected because `init()` runs before test code and pollutes the store. Seed data lives in `main()` only. |

---

## 4. Decisions Made by the Reviewer

| Decision | Choice | Reason |
|---|---|---|
| Monetary type | `float64` | Sufficient for an assignment tracker. A production system must use `int64` (cents) or a decimal library. Documented as a known trade-off. |
| Storage | In-memory `sync.RWMutex` map | No database required per the assignment. The mutex ensures concurrent safety. |
| Bonus features | Swagger + Docker + Makefile | Swagger satisfies the OpenAPI bonus. Docker makes reviewer testing easier. Makefile gives one-line commands for every operation. |
| Module path | `github.com/Amazing-Stardom/ExpenseSync` | Matches the actual repository URL. |
| API docs directory | `api-docs/` | Industry convention separates generated OpenAPI artifacts from hand-written docs. |

---

*Last Updated: August 2026*
