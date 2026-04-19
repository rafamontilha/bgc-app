# BGC Strategic Roadmap - E2E Testing Integrated

**Version:** 3.0
**Date:** 2026-01-16
**Owner:** Product Management Team
**Status:** Active Roadmap with Continuous E2E Validation

---

## Executive Summary

Este roadmap estratégico integra **testes E2E como cidadãos de primeira classe** em cada fase de desenvolvimento, garantindo que cada milestone seja entregue com qualidade production-ready validada de ponta a ponta.

### Current State (v0.4.0)
- **Backend:** 70% production-ready (gaps: graceful shutdown, logging, migrations automation)
- **Frontend:** Stable (Tailwind CSS conflict identified, needs resolution)
- **Data Coverage:** 8% (NCM chapter 17 only) - **BLOCKER for beta**
- **E2E Testing:** Basic suite exists (`api/tests/e2e/simulator_test.go`) but not integrated in CI/CD
- **Technical Debt:** 51 temporary files to clean

### Target State (v0.6.0 - Beta)
- Backend: 100% production-ready with graceful shutdown, structured logging, automated migrations
- Frontend: Production-grade with resolved CSS conflicts, optimized performance
- Data Coverage: 100% (all 97 NCM chapters)
- E2E Testing: Comprehensive suite running on every commit with 95%+ reliability
- CI/CD: Fully automated with E2E gates before deployment

---

## Quality Gates Philosophy

Each phase has **mandatory quality gates** that block progression:

1. **Smoke Tests:** Deploy validation (can the system start?)
2. **Regression Tests:** Don't break existing functionality
3. **Integration Tests:** Full E2E flows work end-to-end
4. **Performance Tests:** Meet latency/throughput targets
5. **Security Tests:** No critical vulnerabilities

**Principle:** No phase marked "Done" until ALL tests pass in production-like environment.

---

## Phase 0: Immediate Cleanup (30 minutes)

**Goal:** Remove technical debt and establish clean baseline for E2E testing.

**Owner:** Platform Team
**Priority:** P0 (blocks all other work)
**Estimated Time:** 30 minutes
**Start Date:** Today (2026-01-16)

### Tasks

| Task | Owner | Time | DoD |
|------|-------|------|-----|
| Delete 51 temp files (`tmpclaude-*`) | Platform | 5min | `git status` shows clean untracked files |
| Remove orphaned Markdown files | Platform | 5min | Only valid docs in `/docs` |
| Clean up `bgcstack/tmpclaude-*` dirs | Platform | 5min | `bgcstack/` directory has no temp folders |
| Update `.gitignore` to prevent future temp files | Platform | 5min | Add `tmpclaude-*`, `nul`, `NUL` patterns |
| Commit cleanup | Platform | 10min | "chore: cleanup temp files and update gitignore" |

### E2E Tests Required

**Smoke Test:**
```bash
# Verify cleanup didn't break anything
docker compose -f bgcstack/docker-compose.yml up -d
curl http://localhost:8080/healthz
curl http://localhost:3000
```

**Success Criteria:**
- All services start successfully
- Health checks pass
- No functional regression

### Definition of Done

- [ ] Zero temporary files in repository (`git status` clean)
- [ ] `.gitignore` updated with temp file patterns
- [ ] All services start and pass health checks
- [ ] Commit pushed to `cleanup/temp-files` branch
- [ ] No regressions (verified by smoke tests)

### Risks & Mitigations

**Risk:** Accidentally delete important files
**Probability:** Low
**Impact:** High
**Mitigation:** Review each file before deletion, keep backup branch

---

## Phase 1: Backend Production-Ready (10 hours)

**Goal:** Bring backend to 100% production-ready state with comprehensive E2E coverage.

**Owner:** Backend Team
**Priority:** P0
**Estimated Time:** 10 hours (1.5 days)
**Dependencies:** Phase 0 complete
**Start Date:** 2026-01-16 afternoon
**End Date:** 2026-01-17 EOD

### Tasks

#### 1.1 Graceful Shutdown (2h)

**Problem:** Kubernetes sends SIGTERM but API doesn't handle it gracefully, causing dropped connections.

**Implementation:**
```go
// cmd/api/main.go
func main() {
    // ... existing setup ...

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    <-quit
    log.Info("shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("server forced to shutdown:", err)
    }

    log.Info("server exited")
}
```

**E2E Test:**
```go
// api/tests/e2e/graceful_shutdown_test.go
func TestGracefulShutdown(t *testing.T) {
    // Start API
    cmd := exec.Command("go", "run", "cmd/api/main.go")
    err := cmd.Start()
    require.NoError(t, err)

    time.Sleep(2 * time.Second) // Wait for startup

    // Start long-running request
    done := make(chan bool)
    go func() {
        resp, _ := http.Get("http://localhost:8080/v1/simulator/destinations?ncm=17011400")
        assert.Equal(t, 200, resp.StatusCode)
        done <- true
    }()

    time.Sleep(100 * time.Millisecond)

    // Send SIGTERM
    cmd.Process.Signal(syscall.SIGTERM)

    // Verify request completes
    select {
    case <-done:
        t.Log("Request completed successfully")
    case <-time.After(35 * time.Second):
        t.Fatal("Request timed out - graceful shutdown failed")
    }

    cmd.Wait()
}
```

**Success Criteria:**
- [ ] SIGTERM handled correctly
- [ ] In-flight requests complete (up to 30s grace period)
- [ ] New requests rejected during shutdown
- [ ] E2E test passes 5 consecutive runs

#### 1.2 Structured Logging (3h)

**Problem:** Logs are inconsistent, making debugging in production difficult.

**Implementation:**
```go
// internal/platform/logger/logger.go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func New(env string) (*zap.Logger, error) {
    var config zap.Config

    if env == "production" {
        config = zap.NewProductionConfig()
        config.EncoderConfig.TimeKey = "timestamp"
        config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    } else {
        config = zap.NewDevelopmentConfig()
    }

    return config.Build()
}

// Usage in handlers
func (h *SimulatorHandler) Handle(c *gin.Context) {
    h.logger.Info("simulator request received",
        zap.String("ncm", req.NCM),
        zap.Int("volume_kg", req.VolumeKg),
        zap.String("user_tier", req.UserTier),
        zap.String("request_id", c.GetString("request_id")),
    )
}
```

**E2E Test:**
```go
func TestStructuredLogging(t *testing.T) {
    // Capture logs
    logFile := "/tmp/bgc-api-test.log"
    os.Setenv("LOG_FILE", logFile)

    // Make API request
    resp, _ := http.Post("http://localhost:8080/v1/simulator/destinations", ...)
    require.Equal(t, 200, resp.StatusCode)

    // Parse log file
    logs, _ := ioutil.ReadFile(logFile)
    var logEntry map[string]interface{}
    json.Unmarshal(logs, &logEntry)

    // Verify structured fields
    assert.Contains(t, logEntry, "timestamp")
    assert.Contains(t, logEntry, "level")
    assert.Contains(t, logEntry, "msg")
    assert.Contains(t, logEntry, "ncm")
    assert.Contains(t, logEntry, "request_id")
}
```

**Success Criteria:**
- [ ] All logs in JSON format in production
- [ ] Request ID propagated through entire request lifecycle
- [ ] Zero `fmt.Println` or `log.Print` in codebase
- [ ] E2E test validates log structure

#### 1.3 Database Migration Automation (3h)

**Problem:** Migrations run manually, error-prone in production.

**Implementation:**
```go
// cmd/migrate/main.go
package main

import (
    "database/sql"
    "log"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    driver, _ := postgres.WithInstance(db, &postgres.Config{})

    m, err := migrate.NewWithDatabaseInstance(
        "file://db/migrations",
        "postgres",
        driver,
    )

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }

    log.Println("Migrations applied successfully")
}
```

**E2E Test:**
```bash
#!/bin/bash
# tests/e2e/migration_test.sh

echo "Testing migration automation..."

# Start fresh database
docker run -d --name test-db -e POSTGRES_PASSWORD=test postgres:16
sleep 5

# Run migrations
export DATABASE_URL="postgresql://postgres:test@localhost:5432/postgres"
go run cmd/migrate/main.go

# Verify schema
docker exec test-db psql -U postgres -c "\dt" | grep "countries_metadata"
docker exec test-db psql -U postgres -c "\dt" | grep "stg.exportacao"

# Cleanup
docker rm -f test-db

echo "✓ Migration test passed"
```

**Success Criteria:**
- [ ] Migrations run idempotently (safe to run multiple times)
- [ ] Migration version tracked in `schema_migrations` table
- [ ] Rollback capability implemented (`migrate down`)
- [ ] E2E test validates schema state after migrations

#### 1.4 Health Check Enhancement (1h)

**Problem:** `/healthz` only checks if API is running, not dependencies.

**Implementation:**
```go
// internal/handler/health.go
func (h *HealthHandler) Check(c *gin.Context) {
    health := map[string]interface{}{
        "status": "ok",
        "timestamp": time.Now().Unix(),
        "checks": map[string]interface{}{},
    }

    // Database check
    ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
    defer cancel()

    if err := h.db.PingContext(ctx); err != nil {
        health["checks"]["database"] = "unhealthy"
        health["status"] = "degraded"
    } else {
        health["checks"]["database"] = "healthy"
    }

    // Redis check
    if err := h.redis.Ping(ctx).Err(); err != nil {
        health["checks"]["redis"] = "unhealthy"
        health["status"] = "degraded"
    } else {
        health["checks"]["redis"] = "healthy"
    }

    statusCode := 200
    if health["status"] == "degraded" {
        statusCode = 503
    }

    c.JSON(statusCode, health)
}
```

**E2E Test:**
```go
func TestHealthCheckDependencies(t *testing.T) {
    tests := []struct{
        name string
        setup func()
        expectedStatus int
        expectedHealth string
    }{
        {
            name: "all healthy",
            setup: func() { /* normal state */ },
            expectedStatus: 200,
            expectedHealth: "ok",
        },
        {
            name: "database down",
            setup: func() {
                exec.Command("docker", "stop", "bgc_db").Run()
            },
            expectedStatus: 503,
            expectedHealth: "degraded",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setup()
            defer func() {
                exec.Command("docker", "start", "bgc_db").Run()
            }()

            resp, _ := http.Get("http://localhost:8080/healthz")
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)

            var health map[string]interface{}
            json.NewDecoder(resp.Body).Decode(&health)
            assert.Equal(t, tt.expectedHealth, health["status"])
        })
    }
}
```

**Success Criteria:**
- [ ] Health check validates database connectivity
- [ ] Health check validates Redis connectivity
- [ ] Returns 503 when dependencies unhealthy
- [ ] E2E test covers all scenarios

#### 1.5 Rate Limiting Enhancement (1h)

**Problem:** Current rate limiting is in-memory, not distributed.

**Implementation:**
```go
// internal/middleware/ratelimit.go
func RateLimitMiddleware(redis *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := "ratelimit:" + c.ClientIP() + ":" + time.Now().Format("2006-01-02")

        count, _ := redis.Incr(c.Request.Context(), key).Result()

        if count == 1 {
            redis.Expire(c.Request.Context(), key, 24*time.Hour)
        }

        c.Header("X-RateLimit-Limit", "5")
        c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, 5-count)))

        if count > 5 {
            c.JSON(429, gin.H{
                "error": "rate_limit_exceeded",
                "message": "Free tier: 5 simulations per day",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

**E2E Test:**
```go
func TestDistributedRateLimiting(t *testing.T) {
    // Start 3 API instances
    for i := 0; i < 3; i++ {
        startAPIInstance(8080 + i)
    }

    client := &http.Client{}

    // Make 5 requests across different instances
    for i := 0; i < 5; i++ {
        port := 8080 + (i % 3)
        resp, _ := client.Post(fmt.Sprintf("http://localhost:%d/v1/simulator/destinations", port), ...)
        assert.Equal(t, 200, resp.StatusCode)
    }

    // 6th request should be rate limited on ANY instance
    resp, _ := client.Post("http://localhost:8082/v1/simulator/destinations", ...)
    assert.Equal(t, 429, resp.StatusCode)
}
```

**Success Criteria:**
- [ ] Rate limit shared across multiple API instances
- [ ] Redis stores rate limit state
- [ ] Headers correctly reflect remaining quota
- [ ] E2E test validates distributed behavior

### E2E Test Suite (Phase 1)

**Location:** `api/tests/e2e/backend_readiness_test.go`

**Test Matrix:**

| Test Name | Scenario | Expected Result | Performance Target |
|-----------|----------|----------------|-------------------|
| `TestHealthCheckComplete` | Call `/healthz` | 200 OK with all checks passing | < 50ms |
| `TestGracefulShutdown` | Send SIGTERM during request | Request completes successfully | < 30s shutdown time |
| `TestStructuredLogs` | Make API request | Logs in JSON with request_id | N/A |
| `TestMigrationIdempotency` | Run migrations 3x | No errors, schema consistent | < 10s |
| `TestDistributedRateLimit` | 6 requests across 3 instances | 5 succeed, 1 rate limited | < 100ms per request |
| `TestDatabaseFailover` | Restart database mid-request | Requests retry and succeed | < 5s recovery |
| `TestRedisFailover` | Stop Redis | API degraded but functional | < 2s detection |

**Run Command:**
```bash
cd api
go test -tags=e2e ./tests/e2e/backend_readiness_test.go -v
```

**CI/CD Integration:**
```yaml
# .github/workflows/backend-e2e.yml
name: Backend E2E Tests

on:
  push:
    paths:
      - 'api/**'
  pull_request:
    paths:
      - 'api/**'

jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Start infrastructure
        run: docker compose -f bgcstack/docker-compose.yml up -d

      - name: Wait for services
        run: |
          timeout 60 bash -c 'until curl -f http://localhost:8080/healthz; do sleep 2; done'

      - name: Run E2E tests
        run: |
          cd api
          go test -tags=e2e -v -count=1 ./tests/e2e/backend_readiness_test.go

      - name: Upload logs on failure
        if: failure()
        uses: actions/upload-artifact@v3
        with:
          name: test-logs
          path: /tmp/*.log
```

### Definition of Done (Phase 1)

**Code Quality:**
- [ ] All 5 tasks implemented and merged to `main`
- [ ] Code review approved by 2+ engineers
- [ ] No linter warnings (`golangci-lint run`)
- [ ] Test coverage > 80% for new code

**Testing:**
- [ ] All 7 E2E tests passing locally
- [ ] E2E tests passing in CI/CD pipeline
- [ ] Load test: 100 concurrent requests without errors
- [ ] Soak test: 1 hour continuous traffic without degradation

**Documentation:**
- [ ] `RUNBOOK.md` updated with graceful shutdown procedure
- [ ] Migration guide documented in `docs/MIGRATIONS.md`
- [ ] Logging format documented with examples

**Performance:**
- [ ] API latency P95 < 200ms (target: 100ms)
- [ ] Health check < 50ms
- [ ] Graceful shutdown < 30s
- [ ] Database queries < 50ms

**Deployment:**
- [ ] Deployed to staging environment
- [ ] Smoke tests passing in staging
- [ ] Rollback plan documented and tested

### Risks & Mitigations

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| E2E tests flaky in CI | Medium | High | Implement retry logic, isolate tests, use test containers |
| Graceful shutdown too long (> 30s) | Low | Medium | Implement timeout for shutdown, force kill after 30s |
| Migration breaks existing data | Low | Critical | Test migrations on production snapshot, implement rollback |
| Distributed rate limit adds latency | Medium | Low | Cache rate limit check, use Redis pipelining |

---

## Phase 2: Data Expansion 8% → 28% (3-4 days)

**Goal:** Expand NCM coverage from chapter 17 (Sugars) to chapters 02 (Meats) and 08 (Fruits).

**Owner:** Data Engineering Team
**Priority:** P0 (blocks beta launch)
**Estimated Time:** 3-4 days
**Dependencies:** Phase 1 complete
**Start Date:** 2026-01-18
**End Date:** 2026-01-21

### Why This Matters (Product Context)

**Current Problem:**
- 92% of users entering NCMs get "no data available" error
- Bounce rate: 80%+ for non-sugar exporters
- NPS: 10 (poor) - users feel misled by "all products" messaging

**Target State:**
- 28% data coverage (3.5x improvement)
- Support top 3 Brazilian export sectors: Sugar, Meat, Fruits
- Projected bounce rate: 40% (50% improvement)
- Projected NPS: 40 (acceptable)

### Tasks

#### 2.1 Data Discovery & Validation (1 day)

**Goal:** Understand existing data and identify gaps.

**Tasks:**
```sql
-- Task 2.1.1: Audit current coverage
SELECT
    SUBSTRING(co_ncm, 1, 2) as chapter,
    COUNT(DISTINCT co_ncm) as unique_ncms,
    COUNT(*) as total_records,
    MIN(co_ano) as min_year,
    MAX(co_ano) as max_year
FROM stg.exportacao
GROUP BY chapter
ORDER BY chapter;

-- Expected: Only chapter 17

-- Task 2.1.2: Check source data availability
SELECT
    ncm_chapter,
    COUNT(*) as records,
    SUM(total_usd) as total_value
FROM trade_ncm_year
WHERE fluxo = 'exportacao'
    AND ncm_chapter IN ('02', '08', '17')
    AND year >= 2023
GROUP BY ncm_chapter
ORDER BY total_value DESC;

-- Expected: Data exists for 02 and 08

-- Task 2.1.3: Identify top NCMs by export value
SELECT
    co_ncm,
    no_ncm_por,
    SUM(total_usd) as total_exports
FROM trade_ncm_year
WHERE fluxo = 'exportacao'
    AND ncm_chapter IN ('02', '08')
    AND year >= 2023
GROUP BY co_ncm, no_ncm_por
ORDER BY total_exports DESC
LIMIT 20;
```

**E2E Validation:**
```bash
# tests/e2e/data_audit_test.sh
#!/bin/bash

echo "Running data audit..."

# Connect to database
export PGPASSWORD=bgc
PSQL="docker exec bgc_db psql -U bgc -d bgc -t -c"

# Test 1: Verify chapter 17 has data
count_17=$($PSQL "SELECT COUNT(*) FROM stg.exportacao WHERE SUBSTRING(co_ncm, 1, 2) = '17'")
if [ "$count_17" -gt 0 ]; then
    echo "✓ Chapter 17 has $count_17 records"
else
    echo "✗ Chapter 17 has no data!"
    exit 1
fi

# Test 2: Check source data for 02 and 08
count_source=$($PSQL "SELECT COUNT(*) FROM trade_ncm_year WHERE ncm_chapter IN ('02', '08')")
if [ "$count_source" -gt 0 ]; then
    echo "✓ Source data available for chapters 02 and 08: $count_source records"
else
    echo "✗ No source data for chapters 02 and 08!"
    exit 1
fi

echo "Data audit passed"
```

**Success Criteria:**
- [ ] Chapter 17 has 50+ records in `stg.exportacao`
- [ ] Chapters 02 and 08 have source data in `trade_ncm_year`
- [ ] Top 20 NCMs identified for each chapter
- [ ] Data quality issues documented (missing countries, zero values, etc.)

#### 2.2 ETL Pipeline Implementation (2 days)

**Goal:** Create automated ETL to populate `stg.exportacao` for chapters 02 and 08.

**Implementation:**
```go
// cmd/jobs/populate_exports/main.go
package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

type ExportRecord struct {
    NCM            string
    CountryCode    string
    Year           int
    Month          int
    TotalUSD       float64
    TotalKG        float64
    AvgPricePerKG  float64
}

func main() {
    db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    defer db.Close()

    chapters := []string{"02", "08"}

    for _, chapter := range chapters {
        log.Printf("Processing chapter %s...", chapter)

        records, err := fetchSourceData(db, chapter)
        if err != nil {
            log.Fatal(err)
        }

        log.Printf("Found %d records for chapter %s", len(records), chapter)

        if err := insertToStaging(db, records); err != nil {
            log.Fatal(err)
        }

        log.Printf("✓ Chapter %s populated successfully", chapter)
    }
}

func fetchSourceData(db *sql.DB, chapter string) ([]ExportRecord, error) {
    query := `
        SELECT
            co_ncm,
            co_pais,
            co_ano,
            co_mes,
            SUM(vl_fob) as total_usd,
            SUM(kg_liquido) as total_kg,
            CASE
                WHEN SUM(kg_liquido) > 0 THEN SUM(vl_fob) / SUM(kg_liquido)
                ELSE 0
            END as avg_price_per_kg
        FROM trade_ncm_year
        WHERE fluxo = 'exportacao'
            AND SUBSTRING(co_ncm, 1, 2) = $1
            AND co_ano >= 2023
        GROUP BY co_ncm, co_pais, co_ano, co_mes
        HAVING SUM(vl_fob) > 0
        ORDER BY total_usd DESC
    `

    rows, err := db.Query(query, chapter)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var records []ExportRecord
    for rows.Next() {
        var r ExportRecord
        rows.Scan(&r.NCM, &r.CountryCode, &r.Year, &r.Month, &r.TotalUSD, &r.TotalKG, &r.AvgPricePerKG)
        records = append(records, r)
    }

    return records, nil
}

func insertToStaging(db *sql.DB, records []ExportRecord) error {
    tx, _ := db.Begin()
    defer tx.Rollback()

    stmt, _ := tx.Prepare(`
        INSERT INTO stg.exportacao
            (co_ncm, co_pais, co_ano, co_mes, vl_fob, kg_liquido, vl_fob_per_kg)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (co_ncm, co_pais, co_ano, co_mes) DO UPDATE SET
            vl_fob = EXCLUDED.vl_fob,
            kg_liquido = EXCLUDED.kg_liquido,
            vl_fob_per_kg = EXCLUDED.vl_fob_per_kg
    `)
    defer stmt.Close()

    for _, r := range records {
        if _, err := stmt.Exec(r.NCM, r.CountryCode, r.Year, r.Month, r.TotalUSD, r.TotalKG, r.AvgPricePerKG); err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

**Kubernetes Job Definition:**
```yaml
# k8s/jobs/populate-exports.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: populate-exports
  namespace: bgc
spec:
  backoffLimit: 3
  template:
    spec:
      restartPolicy: OnFailure
      containers:
      - name: populate
        image: bgc-api:latest
        command: ["./populate_exports"]
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: bgc-secrets
              key: database-url
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "2000m"
```

**E2E Test:**
```go
// api/tests/e2e/data_population_test.go
func TestDataPopulation(t *testing.T) {
    // Setup: Clean staging table
    db := setupTestDB()
    db.Exec("TRUNCATE TABLE stg.exportacao")

    // Run ETL job
    cmd := exec.Command("go", "run", "cmd/jobs/populate_exports/main.go")
    cmd.Env = append(os.Environ(), "DATABASE_URL="+testDBURL)
    output, err := cmd.CombinedOutput()
    require.NoError(t, err, "ETL job failed: %s", output)

    // Verify: Chapter 02 has data
    var count02 int
    db.QueryRow("SELECT COUNT(*) FROM stg.exportacao WHERE SUBSTRING(co_ncm, 1, 2) = '02'").Scan(&count02)
    assert.Greater(t, count02, 100, "Chapter 02 should have at least 100 records")

    // Verify: Chapter 08 has data
    var count08 int
    db.QueryRow("SELECT COUNT(*) FROM stg.exportacao WHERE SUBSTRING(co_ncm, 1, 2) = '08'").Scan(&count08)
    assert.Greater(t, count08, 100, "Chapter 08 should have at least 100 records")

    // Verify: No duplicate records
    var duplicates int
    db.QueryRow(`
        SELECT COUNT(*) FROM (
            SELECT co_ncm, co_pais, co_ano, co_mes, COUNT(*)
            FROM stg.exportacao
            GROUP BY co_ncm, co_pais, co_ano, co_mes
            HAVING COUNT(*) > 1
        ) dup
    `).Scan(&duplicates)
    assert.Equal(t, 0, duplicates, "Should have no duplicate records")

    // Verify: All records have positive values
    var invalidRecords int
    db.QueryRow("SELECT COUNT(*) FROM stg.exportacao WHERE vl_fob <= 0 OR kg_liquido <= 0").Scan(&invalidRecords)
    assert.Equal(t, 0, invalidRecords, "Should have no invalid records")
}
```

**Success Criteria:**
- [ ] ETL job completes successfully for chapters 02 and 08
- [ ] Chapter 02 has 500+ records in `stg.exportacao`
- [ ] Chapter 08 has 500+ records in `stg.exportacao`
- [ ] No duplicate records (unique on NCM + country + year + month)
- [ ] All records have positive USD and KG values
- [ ] E2E test passes 3 consecutive runs

#### 2.3 API Validation (1 day)

**Goal:** Verify simulator endpoint works correctly with new data.

**Test NCMs (Top exports by sector):**

**Chapter 02 (Meats):**
- `02013000` - Carne bovina fresca/refrigerada (desossada)
- `02023000` - Carne bovina congelada (desossada)
- `02071200` - Carne de frango congelada (cortes)

**Chapter 08 (Fruits):**
- `08011100` - Cocos dessecados
- `08012200` - Castanha-do-pará com casca
- `08030019` - Bananas frescas tipo cavendish

**E2E Test Suite:**
```go
// api/tests/e2e/simulator_expanded_data_test.go
func TestSimulatorChapter02Meats(t *testing.T) {
    tests := []struct{
        ncm string
        description string
        minDestinations int
        expectedTopCountry string
    }{
        {"02013000", "Carne bovina desossada", 5, "CN"}, // China is top importer
        {"02023000", "Carne bovina congelada", 4, "US"},
        {"02071200", "Frango congelado", 6, "JP"},      // Japan loves chicken
    }

    for _, tt := range tests {
        t.Run(tt.description, func(t *testing.T) {
            payload := map[string]interface{}{
                "ncm": tt.ncm,
                "volume_kg": 10000,
            }
            body, _ := json.Marshal(payload)

            resp, err := http.Post(
                "http://localhost:8080/v1/simulator/destinations",
                "application/json",
                bytes.NewBuffer(body),
            )
            require.NoError(t, err)
            defer resp.Body.Close()

            assert.Equal(t, 200, resp.StatusCode, "Should return 200 for %s", tt.ncm)

            var result map[string]interface{}
            json.NewDecoder(resp.Body).Decode(&result)

            destinations := result["destinations"].([]interface{})
            assert.GreaterOrEqual(t, len(destinations), tt.minDestinations,
                "Should return at least %d destinations for %s", tt.minDestinations, tt.description)

            // Verify top destination
            if len(destinations) > 0 {
                top := destinations[0].(map[string]interface{})
                assert.Equal(t, tt.expectedTopCountry, top["country_code"],
                    "Top destination for %s should be %s", tt.description, tt.expectedTopCountry)
            }

            // Verify performance
            processingTime := result["metadata"].(map[string]interface{})["processing_time_ms"].(float64)
            assert.Less(t, processingTime, 200.0, "Should process in < 200ms")
        })
    }
}

func TestSimulatorChapter08Fruits(t *testing.T) {
    tests := []struct{
        ncm string
        description string
        minDestinations int
    }{
        {"08011100", "Cocos dessecados", 3},
        {"08012200", "Castanha-do-pará", 4},
        {"08030019", "Bananas frescas", 5},
    }

    for _, tt := range tests {
        t.Run(tt.description, func(t *testing.T) {
            payload := map[string]interface{}{
                "ncm": tt.ncm,
                "volume_kg": 5000,
            }
            body, _ := json.Marshal(payload)

            resp, err := http.Post(
                "http://localhost:8080/v1/simulator/destinations",
                "application/json",
                bytes.NewBuffer(body),
            )
            require.NoError(t, err)
            defer resp.Body.Close()

            assert.Equal(t, 200, resp.StatusCode)

            var result map[string]interface{}
            json.NewDecoder(resp.Body).Decode(&result)

            destinations := result["destinations"].([]interface{})
            assert.GreaterOrEqual(t, len(destinations), tt.minDestinations)

            // Verify all destinations have required fields
            for _, dest := range destinations {
                d := dest.(map[string]interface{})
                assert.NotEmpty(t, d["country_name"])
                assert.NotZero(t, d["score"])
                assert.NotEmpty(t, d["demand_level"])
                assert.NotZero(t, d["market_size_usd"])
            }
        })
    }
}

func TestSimulatorCoverageRegression(t *testing.T) {
    // Verify chapter 17 (sugar) still works after adding 02 and 08
    payload := map[string]interface{}{
        "ncm": "17011400",
        "volume_kg": 1000,
    }
    body, _ := json.Marshal(payload)

    resp, err := http.Post(
        "http://localhost:8080/v1/simulator/destinations",
        "application/json",
        bytes.NewBuffer(body),
    )
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, 200, resp.StatusCode)

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    destinations := result["destinations"].([]interface{})
    assert.GreaterOrEqual(t, len(destinations), 6, "Chapter 17 should still return 6+ destinations")
}
```

**Success Criteria:**
- [ ] All 6 test NCMs return valid results
- [ ] Response time < 200ms for all NCMs
- [ ] Top destinations match expected markets (China for beef, Japan for chicken, etc.)
- [ ] Chapter 17 regression test passes (old data still works)
- [ ] E2E test suite passes 5 consecutive runs

### Frontend Updates (Phase 2)

**Goal:** Update frontend to reflect new data coverage.

#### Task 2.4: Update NCM Autocomplete (2h)

**File:** `web-next/lib/data/ncm-list.ts`

**Changes:**
```typescript
// Add chapter 02 and 08 NCMs
export const AVAILABLE_NCMS: NCMItem[] = [
  // Chapter 02 - Meats
  {
    code: '02013000',
    description: 'Carne bovina fresca, desossada',
    chapter: '02',
    sector: 'Carnes'
  },
  {
    code: '02023000',
    description: 'Carne bovina congelada, desossada',
    chapter: '02',
    sector: 'Carnes'
  },
  {
    code: '02071200',
    description: 'Carne de frango congelada (cortes)',
    chapter: '02',
    sector: 'Carnes'
  },

  // Chapter 08 - Fruits
  {
    code: '08011100',
    description: 'Cocos dessecados',
    chapter: '08',
    sector: 'Frutas'
  },
  {
    code: '08012200',
    description: 'Castanha-do-pará com casca',
    chapter: '08',
    sector: 'Frutas'
  },
  {
    code: '08030019',
    description: 'Bananas frescas tipo cavendish',
    chapter: '08',
    sector: 'Frutas'
  },

  // Chapter 17 - Sugars (existing)
  {
    code: '17011400',
    description: 'Outros açúcares de cana',
    chapter: '17',
    sector: 'Açúcares'
  },
  // ... rest of chapter 17
];

// Add sector grouping for better UX
export const SECTORS = [
  { id: '02', name: 'Carnes', icon: '🥩' },
  { id: '08', name: 'Frutas', icon: '🍎' },
  { id: '17', name: 'Açúcares', icon: '🍬' },
];
```

#### Task 2.5: Update Coverage Banner (1h)

**File:** `web-next/components/home/SimulatorSection.tsx`

**Changes:**
```tsx
<Alert severity="success" sx={{ mb: 3 }}>
  <AlertTitle>Cobertura de Dados - v0.5.0</AlertTitle>
  Simulamos destinos para produtos dos setores:{' '}
  <strong>Açúcares 🍬, Carnes 🥩, Frutas 🍎</strong>
  <br />
  <Typography variant="caption">
    28% de cobertura • 10+ NCMs disponíveis • Dados reais 2023-2026
  </Typography>
</Alert>
```

### E2E Test Suite (Phase 2 Complete)

**Smoke Tests (run on every deployment):**
```bash
#!/bin/bash
# tests/e2e/phase2_smoke.sh

echo "Running Phase 2 smoke tests..."

# Test 1: Chapter 02 (Meat)
curl -X POST http://localhost:8080/v1/simulator/destinations \
  -d '{"ncm":"02013000","volume_kg":10000}' | grep '"country_code"'

# Test 2: Chapter 08 (Fruit)
curl -X POST http://localhost:8080/v1/simulator/destinations \
  -d '{"ncm":"08011100","volume_kg":5000}' | grep '"country_code"'

# Test 3: Chapter 17 (Sugar) - regression
curl -X POST http://localhost:8080/v1/simulator/destinations \
  -d '{"ncm":"17011400","volume_kg":1000}' | grep '"country_code"'

# Test 4: Frontend loads updated NCM list
curl http://localhost:3000/simulator | grep 'Carne bovina'

echo "✓ All smoke tests passed"
```

**Regression Tests:**
```go
// Verify old functionality still works
func TestChapter17StillWorksAfterExpansion(t *testing.T) {
    // Same test as before, validates no breaking changes
}
```

**Integration Tests (Full E2E):**
```typescript
// web-next/tests/e2e/simulator.spec.ts (Playwright)
import { test, expect } from '@playwright/test';

test.describe('Simulator with Expanded Coverage', () => {
  test('should show NCMs from 3 chapters in autocomplete', async ({ page }) => {
    await page.goto('http://localhost:3000/simulator');

    await page.click('input[name="ncm"]');

    // Verify chapter 02 appears
    await expect(page.locator('text=Carne bovina')).toBeVisible();

    // Verify chapter 08 appears
    await expect(page.locator('text=Cocos dessecados')).toBeVisible();

    // Verify chapter 17 still appears
    await expect(page.locator('text=Açúcar de cana')).toBeVisible();
  });

  test('should simulate beef export successfully', async ({ page }) => {
    await page.goto('http://localhost:3000/simulator');

    // Select beef NCM
    await page.fill('input[name="ncm"]', '02013000');
    await page.fill('input[name="volume_kg"]', '10000');

    await page.click('button:has-text("Simular")');

    // Wait for results
    await page.waitForSelector('[data-testid="destination-card"]');

    // Verify results
    const cards = page.locator('[data-testid="destination-card"]');
    await expect(cards).toHaveCount({ minimum: 5 });

    // Verify China is likely top destination for beef
    const firstCard = cards.first();
    await expect(firstCard.locator('text=China')).toBeVisible();
  });

  test('should handle all 3 sectors without errors', async ({ page }) => {
    const testCases = [
      { ncm: '02013000', sector: 'Carnes', product: 'Carne bovina' },
      { ncm: '08011100', sector: 'Frutas', product: 'Cocos' },
      { ncm: '17011400', sector: 'Açúcares', product: 'Açúcar' },
    ];

    for (const tc of testCases) {
      await page.goto('http://localhost:3000/simulator');

      await page.fill('input[name="ncm"]', tc.ncm);
      await page.fill('input[name="volume_kg"]', '5000');
      await page.click('button:has-text("Simular")');

      await page.waitForSelector('[data-testid="destination-card"]');

      const cards = page.locator('[data-testid="destination-card"]');
      await expect(cards).toHaveCount({ minimum: 3 });

      console.log(`✓ ${tc.sector} (${tc.product}) test passed`);
    }
  });
});
```

**Performance Tests:**
```bash
# tests/e2e/phase2_load_test.sh
#!/bin/bash

echo "Running Phase 2 load test..."

# Use Apache Bench to simulate 100 concurrent users
ab -n 1000 -c 100 -p payload_beef.json -T application/json \
  http://localhost:8080/v1/simulator/destinations

# Verify P95 latency < 200ms
echo "Checking latency metrics..."
curl http://localhost:9090/api/v1/query?query='histogram_quantile(0.95,rate(http_request_duration_seconds_bucket[5m]))'
```

### Definition of Done (Phase 2)

**Data Quality:**
- [ ] Chapter 02 has 500+ records covering top 20 beef/poultry NCMs
- [ ] Chapter 08 has 500+ records covering top 20 fruit NCMs
- [ ] Data freshness: all records from 2023-2026
- [ ] Zero duplicate records
- [ ] Zero records with invalid values (negative USD/KG)

**API Validation:**
- [ ] 6 test NCMs (3 meats, 3 fruits) return valid results
- [ ] All responses < 200ms (P95)
- [ ] Chapter 17 regression tests pass
- [ ] Rate limiting still works correctly

**Frontend:**
- [ ] Autocomplete shows 10+ NCMs from 3 chapters
- [ ] Coverage banner updated to "28% • 3 setores"
- [ ] UI handles all 3 sectors without errors
- [ ] No console errors or warnings

**E2E Testing:**
- [ ] Smoke tests pass (3 chapters tested)
- [ ] Regression tests pass (chapter 17 still works)
- [ ] Integration tests pass (Playwright suite)
- [ ] Load test passes (100 concurrent users, P95 < 200ms)

**Documentation:**
- [ ] `RELATORIO-NCM-COVERAGE-V0.5.0.md` published
- [ ] `CHANGELOG.md` updated with new NCMs
- [ ] API docs updated with example NCMs from all 3 chapters

**Deployment:**
- [ ] ETL job runs successfully in staging
- [ ] API deployed to staging and smoke tested
- [ ] Frontend deployed to staging
- [ ] Full E2E suite passes in staging environment

### Risks & Mitigations

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| ETL job fails for chapter 02/08 | Medium | High | Implement robust error handling, log all failures, add retry logic |
| Data quality issues (bad countries, zero values) | Medium | Medium | Add validation step in ETL, reject invalid records |
| Performance degrades with 3.5x more data | Low | Medium | Add database indexes on `co_ncm` and `co_pais`, monitor query performance |
| Frontend autocomplete too slow with 10+ NCMs | Low | Low | Implement virtual scrolling, limit autocomplete results to 20 |

---

## Phase 3: Authentication & User Management (2 days)

**Goal:** Implement Clerk authentication to enable user accounts, premium tiers, and personalized experiences.

**Owner:** Frontend Team + Backend Team
**Priority:** P1 (required for monetization)
**Estimated Time:** 2 days
**Dependencies:** Phase 1 and 2 complete
**Start Date:** 2026-01-22
**End Date:** 2026-01-23

### Why This Matters (Product Context)

**Current Problem:**
- Rate limiting by IP address (easy to bypass with VPN)
- No way to identify returning users
- Cannot offer premium tiers
- No user analytics (cohort analysis, retention)

**Target State:**
- User accounts with email/password or social login
- Rate limiting by user ID (accurate tracking)
- Premium tier with unlimited simulations
- User analytics dashboard

### Tasks

#### 3.1 Clerk Integration - Frontend (4h)

**Implementation:**
```bash
cd web-next
pnpm add @clerk/nextjs
```

**File:** `web-next/app/layout.tsx`
```typescript
import { ClerkProvider } from '@clerk/nextjs'

export default function RootLayout({ children }) {
  return (
    <ClerkProvider>
      <html lang="pt-BR">
        <body>{children}</body>
      </html>
    </ClerkProvider>
  )
}
```

**File:** `web-next/middleware.ts`
```typescript
import { authMiddleware } from "@clerk/nextjs";

export default authMiddleware({
  publicRoutes: ["/", "/simulator"],
  ignoredRoutes: ["/api/public"],
});

export const config = {
  matcher: ["/((?!.+\\.[\\w]+$|_next).*)", "/", "/(api|trpc)(.*)"],
};
```

**File:** `web-next/components/Header.tsx`
```typescript
import { UserButton, SignInButton } from '@clerk/nextjs';
import { useUser } from '@clerk/nextjs';

export function Header() {
  const { isSignedIn, user } = useUser();

  return (
    <header>
      <Logo />
      <nav>
        <Link href="/simulator">Simulador</Link>
        <Link href="/dashboard">Dashboard</Link>
      </nav>
      {isSignedIn ? (
        <UserButton afterSignOutUrl="/" />
      ) : (
        <SignInButton mode="modal">
          <Button>Login</Button>
        </SignInButton>
      )}
    </header>
  );
}
```

**E2E Test (Playwright):**
```typescript
// web-next/tests/e2e/auth.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
  test('should show sign in button when not authenticated', async ({ page }) => {
    await page.goto('http://localhost:3000');

    const signInButton = page.locator('button:has-text("Login")');
    await expect(signInButton).toBeVisible();
  });

  test('should sign in with email', async ({ page }) => {
    await page.goto('http://localhost:3000');

    // Click sign in
    await page.click('button:has-text("Login")');

    // Wait for Clerk modal
    await page.waitForSelector('[data-clerk-id]');

    // Fill email and password
    await page.fill('input[name="identifier"]', 'test@example.com');
    await page.click('button:has-text("Continue")');

    await page.fill('input[name="password"]', 'TestPassword123!');
    await page.click('button:has-text("Continue")');

    // Verify signed in
    await page.waitForSelector('[data-testid="user-button"]');
    await expect(page.locator('[data-testid="user-button"]')).toBeVisible();
  });

  test('should access protected dashboard after sign in', async ({ page }) => {
    // Sign in first
    await signIn(page, 'test@example.com', 'TestPassword123!');

    // Navigate to dashboard
    await page.goto('http://localhost:3000/dashboard');

    // Verify dashboard loads (not redirected to sign in)
    await expect(page).toHaveURL(/\/dashboard/);
    await expect(page.locator('h1:has-text("Dashboard")')).toBeVisible();
  });

  test('should sign out successfully', async ({ page }) => {
    await signIn(page, 'test@example.com', 'TestPassword123!');

    // Click user button
    await page.click('[data-testid="user-button"]');

    // Click sign out
    await page.click('button:has-text("Sign out")');

    // Verify signed out
    await page.waitForSelector('button:has-text("Login")');
    await expect(page.locator('button:has-text("Login")')).toBeVisible();
  });
});

// Helper function
async function signIn(page, email, password) {
  await page.goto('http://localhost:3000');
  await page.click('button:has-text("Login")');
  await page.fill('input[name="identifier"]', email);
  await page.click('button:has-text("Continue")');
  await page.fill('input[name="password"]', password);
  await page.click('button:has-text("Continue")');
  await page.waitForSelector('[data-testid="user-button"]');
}
```

#### 3.2 Backend JWT Validation (4h)

**Implementation:**
```go
// internal/middleware/auth.go
package middleware

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

type ClerkJWKS struct {
    Keys []JWK `json:"keys"`
}

type JWK struct {
    Kid string   `json:"kid"`
    Kty string   `json:"kty"`
    Use string   `json:"use"`
    N   string   `json:"n"`
    E   string   `json:"e"`
}

func AuthMiddleware(clerkDomain string) gin.HandlerFunc {
    // Fetch JWKS from Clerk
    jwks := fetchJWKS(clerkDomain)

    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "missing_auth_token"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // Parse and validate JWT
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Get public key from JWKS
            kid := token.Header["kid"].(string)
            publicKey := getPublicKey(jwks, kid)
            return publicKey, nil
        })

        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "invalid_token"})
            c.Abort()
            return
        }

        // Extract claims
        claims := token.Claims.(jwt.MapClaims)
        userID := claims["sub"].(string)

        // Store user ID in context
        c.Set("user_id", userID)
        c.Next()
    }
}

func fetchJWKS(domain string) ClerkJWKS {
    url := fmt.Sprintf("https://%s/.well-known/jwks.json", domain)
    resp, _ := http.Get(url)
    defer resp.Body.Close()

    var jwks ClerkJWKS
    json.NewDecoder(resp.Body).Decode(&jwks)
    return jwks
}
```

**Update Rate Limiting:**
```go
// internal/middleware/ratelimit.go
func RateLimitMiddleware(redis *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get user ID from context (set by AuthMiddleware)
        userID, exists := c.Get("user_id")

        var key string
        var limit int

        if exists {
            // Authenticated user
            key = "ratelimit:user:" + userID.(string) + ":" + time.Now().Format("2006-01-02")

            // Check user tier
            tier := getUserTier(userID.(string)) // Query from database
            if tier == "premium" {
                limit = 10000 // Effectively unlimited
            } else {
                limit = 5 // Free tier
            }
        } else {
            // Anonymous user (fallback to IP)
            key = "ratelimit:ip:" + c.ClientIP() + ":" + time.Now().Format("2006-01-02")
            limit = 5
        }

        count, _ := redis.Incr(c.Request.Context(), key).Result()

        if count == 1 {
            redis.Expire(c.Request.Context(), key, 24*time.Hour)
        }

        c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
        c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, limit-count)))

        if count > int64(limit) {
            c.JSON(429, gin.H{
                "error": "rate_limit_exceeded",
                "message": "Upgrade to premium for unlimited simulations",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

**E2E Test:**
```go
// api/tests/e2e/auth_test.go
func TestAuthenticatedRateLimit(t *testing.T) {
    // Create test user and get JWT
    token := createTestUserAndGetToken()

    client := &http.Client{}

    // Make 5 requests (free tier limit)
    for i := 0; i < 5; i++ {
        req, _ := http.NewRequest("POST", "http://localhost:8080/v1/simulator/destinations", ...)
        req.Header.Set("Authorization", "Bearer "+token)
        resp, _ := client.Do(req)
        assert.Equal(t, 200, resp.StatusCode)
        resp.Body.Close()
    }

    // 6th request should be rate limited
    req, _ := http.NewRequest("POST", "http://localhost:8080/v1/simulator/destinations", ...)
    req.Header.Set("Authorization", "Bearer "+token)
    resp, _ := client.Do(req)
    assert.Equal(t, 429, resp.StatusCode)
}

func TestPremiumUserUnlimited(t *testing.T) {
    // Create premium user
    token := createPremiumUserAndGetToken()

    client := &http.Client{}

    // Make 100 requests (should all succeed)
    for i := 0; i < 100; i++ {
        req, _ := http.NewRequest("POST", "http://localhost:8080/v1/simulator/destinations", ...)
        req.Header.Set("Authorization", "Bearer "+token)
        resp, _ := client.Do(req)
        assert.Equal(t, 200, resp.StatusCode)
        resp.Body.Close()
    }
}

func TestInvalidToken(t *testing.T) {
    req, _ := http.NewRequest("POST", "http://localhost:8080/v1/simulator/destinations", ...)
    req.Header.Set("Authorization", "Bearer invalid-token-12345")

    resp, _ := http.DefaultClient.Do(req)
    assert.Equal(t, 401, resp.StatusCode)

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    assert.Equal(t, "invalid_token", result["error"])
}
```

#### 3.3 User Database Schema (2h)

**Migration:**
```sql
-- db/migrations/0013_users_schema.sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clerk_user_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    tier VARCHAR(50) DEFAULT 'free' CHECK (tier IN ('free', 'premium')),
    simulations_used_today INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_clerk_id ON users(clerk_user_id);
CREATE INDEX idx_users_tier ON users(tier);

-- Reset counter daily
CREATE OR REPLACE FUNCTION reset_daily_counter()
RETURNS void AS $$
BEGIN
    UPDATE users SET simulations_used_today = 0;
END;
$$ LANGUAGE plpgsql;

-- Schedule daily reset (via cron job or Kubernetes CronJob)
```

**E2E Test:**
```bash
# tests/e2e/user_schema_test.sh
#!/bin/bash

echo "Testing user schema..."

export PGPASSWORD=bgc
PSQL="docker exec bgc_db psql -U bgc -d bgc -c"

# Test 1: Insert test user
$PSQL "
    INSERT INTO users (clerk_user_id, email, full_name, tier)
    VALUES ('user_123', 'test@example.com', 'Test User', 'free')
    RETURNING id;
"

# Test 2: Verify tier constraint
$PSQL "
    INSERT INTO users (clerk_user_id, email, tier)
    VALUES ('user_456', 'premium@example.com', 'premium')
    RETURNING id;
"

# Test 3: Try invalid tier (should fail)
$PSQL "
    INSERT INTO users (clerk_user_id, email, tier)
    VALUES ('user_789', 'invalid@example.com', 'enterprise');
" && exit 1 || echo "✓ Invalid tier rejected correctly"

# Test 4: Reset counter function
$PSQL "UPDATE users SET simulations_used_today = 10;"
$PSQL "SELECT reset_daily_counter();"
count=$($PSQL "SELECT SUM(simulations_used_today) FROM users;" -t)
if [ "$count" -eq 0 ]; then
    echo "✓ Counter reset successfully"
else
    echo "✗ Counter reset failed"
    exit 1
fi

echo "User schema tests passed"
```

### E2E Test Suite (Phase 3 Complete)

**Smoke Tests:**
```bash
# tests/e2e/phase3_smoke.sh
#!/bin/bash

# Test 1: Sign in flow
curl -X POST http://localhost:3000/api/auth/signin \
  -d '{"email":"test@example.com","password":"Test123!"}' | grep 'token'

# Test 2: Protected endpoint without auth (should fail)
curl -X POST http://localhost:8080/v1/simulator/destinations \
  -d '{"ncm":"17011400"}' -w "%{http_code}" | grep 401

# Test 3: Protected endpoint with auth (should succeed)
TOKEN=$(get_test_token)
curl -X POST http://localhost:8080/v1/simulator/destinations \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"ncm":"17011400"}' -w "%{http_code}" | grep 200

echo "✓ Auth smoke tests passed"
```

**Regression Tests:**
```typescript
// Verify anonymous users can still use simulator (backwards compatibility)
test('anonymous user can use simulator with IP-based rate limit', async ({ page }) => {
  await page.goto('http://localhost:3000/simulator');

  // No sign in
  await page.fill('input[name="ncm"]', '17011400');
  await page.click('button:has-text("Simular")');

  // Should work (anonymous allowed)
  await page.waitForSelector('[data-testid="destination-card"]');
  await expect(page.locator('[data-testid="destination-card"]')).toHaveCount({ minimum: 1 });
});
```

**Integration Tests:**
```typescript
test.describe('Full Auth Flow E2E', () => {
  test('free user hits rate limit, upgrades, gets unlimited access', async ({ page }) => {
    // Sign in as free user
    await signIn(page, 'free@example.com', 'Test123!');

    // Use all 5 free simulations
    for (let i = 0; i < 5; i++) {
      await page.goto('http://localhost:3000/simulator');
      await page.fill('input[name="ncm"]', '17011400');
      await page.click('button:has-text("Simular")');
      await page.waitForSelector('[data-testid="destination-card"]');
    }

    // 6th attempt shows upgrade modal
    await page.goto('http://localhost:3000/simulator');
    await page.fill('input[name="ncm"]', '17011400');
    await page.click('button:has-text("Simular")');

    await page.waitForSelector('[data-testid="upgrade-modal"]');
    await expect(page.locator('text=Limite Atingido')).toBeVisible();

    // Click upgrade
    await page.click('button:has-text("Fazer Upgrade")');

    // Mock payment success
    await mockStripePayment(page);

    // Verify upgraded
    await page.waitForSelector('text=Premium');

    // Try simulation again (should work now)
    await page.goto('http://localhost:3000/simulator');
    await page.fill('input[name="ncm"]', '17011400');
    await page.click('button:has-text("Simular")');
    await page.waitForSelector('[data-testid="destination-card"]');

    // Verify no rate limit banner
    await expect(page.locator('[data-testid="rate-limit-banner"]')).not.toBeVisible();
  });
});
```

### Definition of Done (Phase 3)

**Frontend:**
- [ ] Clerk integration complete
- [ ] Sign in/up flow works (email + social)
- [ ] Protected routes redirect to sign in
- [ ] User button with profile/settings
- [ ] Upgrade modal appears at rate limit

**Backend:**
- [ ] JWT validation middleware implemented
- [ ] Rate limiting by user ID (not IP)
- [ ] Premium users have unlimited access
- [ ] User tier stored in database
- [ ] `/me` endpoint returns user info

**Database:**
- [ ] `users` table created
- [ ] Tier constraint enforced
- [ ] Daily counter reset function works
- [ ] Indexes on `clerk_user_id` and `tier`

**E2E Testing:**
- [ ] Auth smoke tests pass
- [ ] Regression tests pass (anonymous still works)
- [ ] Full upgrade flow test passes
- [ ] JWT validation test passes
- [ ] Rate limit by user ID test passes

**Documentation:**
- [ ] Auth flow documented in `docs/AUTHENTICATION.md`
- [ ] API docs updated with Bearer token examples
- [ ] User tier migration guide

**Deployment:**
- [ ] Clerk API keys in Kubernetes secrets
- [ ] Frontend env vars configured
- [ ] Database migration applied in staging
- [ ] E2E tests pass in staging

### Risks & Mitigations

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Clerk integration breaks existing users | Low | Critical | Keep anonymous access, gradual migration |
| JWT validation too slow | Low | Medium | Cache JWKS, use in-memory key store |
| Daily counter reset misses some users | Medium | Low | Use reliable cron (Kubernetes CronJob), log all resets |

---

## Phase 4: Private Beta Launch (2-3 weeks)

**Goal:** Launch private beta with 20-50 SME exporters, collect feedback, validate product-market fit.

**Owner:** Product Manager + Growth Team
**Priority:** P0 (go-to-market milestone)
**Estimated Time:** 2-3 weeks
**Dependencies:** Phases 1, 2, 3 complete
**Start Date:** 2026-01-24
**End Date:** 2026-02-14

### Tasks

#### 4.1 Beta User Recruitment (1 week)

**Channels:**
- LinkedIn outreach (export associations, SME groups)
- Cold email to top 1000 Brazilian exporters
- Partner with APEX Brasil
- Referrals from existing network

**Target:** 50 signups, 20 active users

**E2E Test (Analytics):**
```typescript
// Track signup funnel
test('beta signup flow has < 30% drop-off', async ({ page }) => {
  // Visit landing page
  await page.goto('http://localhost:3000');

  // Track funnel
  await page.click('button:has-text("Solicitar Convite")');
  await page.fill('input[name="email"]', 'beta@example.com');
  await page.fill('input[name="company"]', 'Test Exporter Ltd');
  await page.fill('input[name="ncm"]', '17011400');
  await page.click('button:has-text("Enviar Solicitação")');

  // Verify confirmation
  await expect(page.locator('text=Convite enviado')).toBeVisible();

  // Check analytics event
  const events = await getAnalyticsEvents(page);
  expect(events).toContain('beta_signup_completed');
});
```

#### 4.2 Onboarding Experience (3 days)

**First-time user flow:**
1. Welcome email with login link
2. Onboarding tour (4 steps)
3. First simulation guided
4. Feedback prompt

**E2E Test:**
```typescript
test('new beta user completes onboarding', async ({ page }) => {
  // Click magic link from email
  await page.goto('http://localhost:3000/onboarding?token=xyz');

  // Step 1: Welcome
  await expect(page.locator('h1:has-text("Bem-vindo")')).toBeVisible();
  await page.click('button:has-text("Começar")');

  // Step 2: Enter NCM
  await page.fill('input[name="ncm"]', '17011400');
  await page.click('button:has-text("Próximo")');

  // Step 3: Review results
  await page.waitForSelector('[data-testid="destination-card"]');
  await page.click('button:has-text("Entendi")');

  // Step 4: Feedback
  await page.fill('textarea[name="feedback"]', 'Great tool!');
  await page.click('button:has-text("Enviar Feedback")');

  // Verify completed
  await expect(page).toHaveURL(/\/simulator/);
});
```

#### 4.3 Feedback Collection (ongoing)

**Methods:**
- In-app feedback widget
- Weekly email survey (NPS, satisfaction)
- 1:1 interviews (5 users/week)
- Usage analytics (Amplitude/Mixpanel)

**Key Metrics:**
- NPS (target: > 40)
- Simulation completion rate (target: > 70%)
- Weekly Active Users / Monthly Signups (target: > 60%)
- Feature requests categorized by RICE score

**E2E Test (Analytics Pipeline):**
```bash
# tests/e2e/analytics_pipeline_test.sh
#!/bin/bash

# Test 1: Event tracking
curl -X POST http://localhost:8080/analytics/track \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"event":"simulation_completed","properties":{"ncm":"17011400"}}'

# Test 2: Verify event in database
count=$(docker exec bgc_db psql -U bgc -d bgc -t -c "SELECT COUNT(*) FROM analytics_events WHERE event='simulation_completed'")
if [ "$count" -gt 0 ]; then
    echo "✓ Event tracked successfully"
else
    echo "✗ Event not tracked"
    exit 1
fi

# Test 3: Aggregate metrics
curl http://localhost:8080/analytics/metrics | grep 'simulations_today'

echo "Analytics pipeline test passed"
```

#### 4.4 Incident Response & Monitoring (ongoing)

**SLA for Beta:**
- P0 (system down): 1 hour response, 4 hour resolution
- P1 (major feature broken): 4 hour response, 24 hour resolution
- P2 (minor bug): 24 hour response, 1 week resolution

**E2E Test (Alerting):**
```bash
# tests/e2e/alerting_test.sh
#!/bin/bash

# Test 1: High error rate triggers alert
# Simulate 100 errors
for i in {1..100}; do
    curl -X POST http://localhost:8080/v1/simulator/destinations -d '{"ncm":"invalid"}'
done

# Wait for alert (check Prometheus Alertmanager)
sleep 60
alert=$(curl -s http://localhost:9093/api/v2/alerts | grep 'HighErrorRate')
if [ -n "$alert" ]; then
    echo "✓ Alert triggered for high error rate"
else
    echo "✗ Alert not triggered"
    exit 1
fi

# Test 2: High latency triggers alert
# Simulate slow queries (add artificial delay)
echo "Testing latency alert..."

# Test 3: Database down triggers alert
docker stop bgc_db
sleep 30
alert=$(curl -s http://localhost:9093/api/v2/alerts | grep 'DatabaseDown')
docker start bgc_db
if [ -n "$alert" ]; then
    echo "✓ Alert triggered for database down"
else
    echo "✗ Alert not triggered"
fi

echo "Alerting tests passed"
```

### E2E Test Suite (Phase 4)

**Beta User Journey (Critical Path):**
```typescript
test.describe('Beta User Critical Path', () => {
  test('complete user journey from signup to second simulation', async ({ page, context }) => {
    // 1. Signup
    await page.goto('http://localhost:3000');
    await page.click('button:has-text("Solicitar Convite Beta")');
    await page.fill('input[name="email"]', 'beta-user@example.com');
    await page.fill('input[name="company"]', 'Acme Exports');
    await page.click('button:has-text("Solicitar")');

    // 2. Receive invite email (mock)
    const inviteLink = await getInviteLink('beta-user@example.com');

    // 3. Click invite and onboard
    await page.goto(inviteLink);
    await page.fill('input[name="password"]', 'SecurePass123!');
    await page.click('button:has-text("Criar Conta")');

    // 4. Onboarding tour
    await expect(page.locator('h1:has-text("Bem-vindo")')).toBeVisible();
    await page.click('button:has-text("Começar Tour")');

    // Skip tour steps
    for (let i = 0; i < 3; i++) {
      await page.click('button:has-text("Próximo")');
    }

    // 5. First simulation
    await page.fill('input[name="ncm"]', '02013000'); // Beef
    await page.fill('input[name="volume_kg"]', '10000');
    await page.click('button:has-text("Simular")');

    // 6. View results
    await page.waitForSelector('[data-testid="destination-card"]');
    const cards = page.locator('[data-testid="destination-card"]');
    await expect(cards).toHaveCount({ minimum: 5 });

    // 7. Click first destination for details
    await cards.first().click();
    await expect(page.locator('[data-testid="destination-details"]')).toBeVisible();

    // 8. Save simulation
    await page.click('button:has-text("Salvar Simulação")');
    await expect(page.locator('text=Salvo com sucesso')).toBeVisible();

    // 9. Return to simulator
    await page.goto('http://localhost:3000/simulator');

    // 10. Second simulation (different NCM)
    await page.fill('input[name="ncm"]', '08011100'); // Coconuts
    await page.fill('input[name="volume_kg"]', '5000');
    await page.click('button:has-text("Simular")');

    await page.waitForSelector('[data-testid="destination-card"]');

    // 11. View saved simulations
    await page.goto('http://localhost:3000/dashboard/simulations');
    const savedSims = page.locator('[data-testid="saved-simulation"]');
    await expect(savedSims).toHaveCount(2); // Both simulations saved

    // 12. Give feedback
    await page.click('[data-testid="feedback-button"]');
    await page.fill('textarea[name="feedback"]', 'Love the tool! Very helpful.');
    await page.click('button:has-text("Enviar")');

    // Verify entire flow tracked in analytics
    const analyticsEvents = await getAnalyticsEvents(context);
    expect(analyticsEvents).toContain('signup_completed');
    expect(analyticsEvents).toContain('onboarding_completed');
    expect(analyticsEvents).toContain('simulation_completed');
    expect(analyticsEvents).toContain('simulation_saved');
    expect(analyticsEvents).toContain('feedback_submitted');
  });
});
```

**Performance Under Beta Load:**
```bash
# tests/e2e/beta_load_test.sh
#!/bin/bash

echo "Simulating beta user load (20 concurrent users)..."

# Use k6 for load testing
k6 run - <<EOF
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 20, // 20 concurrent users
  duration: '10m',
  thresholds: {
    http_req_duration: ['p(95)<200'], // 95% of requests under 200ms
    http_req_failed: ['rate<0.01'],   // Less than 1% errors
  },
};

export default function () {
  // Simulate user behavior
  let res = http.post('http://localhost:8080/v1/simulator/destinations',
    JSON.stringify({
      ncm: '17011400',
      volume_kg: 10000,
    }),
    {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer \${__ENV.TEST_TOKEN}',
      },
    }
  );

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });

  sleep(5); // 5 seconds between requests (realistic user behavior)
}
EOF

echo "✓ Beta load test completed"
```

### Definition of Done (Phase 4)

**User Acquisition:**
- [ ] 50 beta signups
- [ ] 20 active users (used simulator 3+ times)
- [ ] 10 power users (used simulator 10+ times)

**Product Metrics:**
- [ ] NPS: > 40 (acceptable)
- [ ] Simulation completion rate: > 70%
- [ ] Weekly Active Users / Monthly Signups: > 60%
- [ ] Average simulations per user: > 3

**Quality Metrics:**
- [ ] Uptime: > 99.5% during beta
- [ ] P95 latency: < 200ms
- [ ] Error rate: < 0.5%
- [ ] Zero P0 incidents lasting > 4 hours

**E2E Testing:**
- [ ] Beta user journey test passes
- [ ] Load test (20 concurrent users) passes
- [ ] Analytics pipeline test passes
- [ ] Alerting test passes

**Feedback:**
- [ ] 10+ user interviews completed
- [ ] 3+ feature requests prioritized in roadmap
- [ ] 5+ usability issues fixed
- [ ] Product-market fit score: > 40%

**Documentation:**
- [ ] User guide published
- [ ] FAQ updated
- [ ] Known issues documented
- [ ] Beta feedback report published

### Risks & Mitigations

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Low user engagement (< 50% activation) | Medium | High | Improve onboarding, proactive outreach, weekly tips email |
| Technical issues during beta | Medium | High | 24/7 monitoring, fast incident response, rollback plan |
| Negative feedback (NPS < 20) | Low | Critical | Pivot messaging, fix top issues, re-recruit users |
| Data quality complaints | Medium | Medium | Document coverage, set expectations, prioritize data expansion |

---

## Phase 5: Kubernetes Production Deployment (Parallel Track)

**Goal:** Deploy to production Kubernetes cluster with HA, auto-scaling, and zero-downtime deployments.

**Owner:** DevOps Team
**Priority:** P1 (parallel with beta)
**Estimated Time:** 2 weeks (parallel)
**Dependencies:** Phase 1 complete
**Start Date:** 2026-01-24
**End Date:** 2026-02-07

*(This is a parallel track - runs alongside Phase 4)*

### Tasks

#### 5.1 Kubernetes Cluster Setup (3 days)

**Infrastructure:**
- EKS (AWS) or GKE (Google Cloud) or AKS (Azure)
- 3 availability zones
- Node pools: 2 t3.medium (API) + 1 r5.large (Database)

**E2E Test (Cluster Health):**
```bash
# tests/e2e/k8s_cluster_test.sh
#!/bin/bash

echo "Testing Kubernetes cluster health..."

# Test 1: Nodes ready
node_count=$(kubectl get nodes --no-headers | wc -l)
if [ "$node_count" -ge 3 ]; then
    echo "✓ All nodes ready ($node_count nodes)"
else
    echo "✗ Insufficient nodes"
    exit 1
fi

# Test 2: Core services running
kubectl get pods -n kube-system | grep -E 'coredns|kube-proxy' | grep Running

# Test 3: Load balancer provisioned
kubectl get svc -n ingress-nginx | grep LoadBalancer

echo "Cluster health test passed"
```

#### 5.2 GitOps with ArgoCD (2 days)

**E2E Test:**
```bash
# tests/e2e/argocd_test.sh
#!/bin/bash

# Test 1: ArgoCD syncs app
kubectl patch app bgc-app -n argocd --type merge -p '{"spec":{"source":{"targetRevision":"main"}}}'
kubectl wait --for=condition=Synced app/bgc-app -n argocd --timeout=300s

# Test 2: Rolling update
kubectl set image deployment/bgc-api bgc-api=bgc-api:v0.5.0 -n bgc
kubectl rollout status deployment/bgc-api -n bgc --timeout=300s

# Test 3: Rollback
kubectl rollout undo deployment/bgc-api -n bgc
kubectl rollout status deployment/bgc-api -n bgc --timeout=300s

echo "ArgoCD test passed"
```

#### 5.3 Blue-Green Deployment (3 days)

**E2E Test:**
```bash
# tests/e2e/blue_green_test.sh
#!/bin/bash

# Deploy green version
kubectl apply -f k8s/blue-green/deployment-green.yaml

# Wait for green to be ready
kubectl wait --for=condition=Available deployment/bgc-api-green -n bgc --timeout=300s

# Run smoke tests against green
GREEN_URL=$(kubectl get svc bgc-api-green -n bgc -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
curl -f http://$GREEN_URL/healthz

# Switch traffic to green
kubectl patch svc bgc-api -n bgc -p '{"spec":{"selector":{"version":"green"}}}'

# Verify traffic switched
curl -f http://api.bgc.com/healthz

# Cleanup blue
kubectl delete deployment bgc-api-blue -n bgc

echo "Blue-green deployment test passed"
```

### E2E Test Suite (Phase 5)

**Full Deployment Pipeline:**
```yaml
# .github/workflows/deploy-production.yml
name: Deploy to Production

on:
  push:
    tags:
      - 'v*'

jobs:
  e2e-pre-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run E2E tests (staging)
        run: |
          cd api
          go test -tags=e2e -v ./tests/e2e/...

      - name: Build Docker images
        run: |
          docker build -t bgc-api:${{ github.ref_name }} api/
          docker build -t bgc-web:${{ github.ref_name }} web-next/

      - name: Push to registry
        run: |
          docker push bgc-api:${{ github.ref_name }}
          docker push bgc-web:${{ github.ref_name }}

  deploy-green:
    needs: e2e-pre-deploy
    runs-on: ubuntu-latest
    steps:
      - name: Deploy green version
        run: |
          kubectl set image deployment/bgc-api-green bgc-api=bgc-api:${{ github.ref_name }} -n bgc
          kubectl rollout status deployment/bgc-api-green -n bgc --timeout=300s

  e2e-post-deploy:
    needs: deploy-green
    runs-on: ubuntu-latest
    steps:
      - name: Run smoke tests (green)
        run: |
          bash tests/e2e/smoke_tests.sh green.bgc.com

      - name: Run load test (green)
        run: |
          k6 run tests/e2e/load_test.js --env URL=https://green.bgc.com

      - name: Check error rate
        run: |
          # Query Prometheus for error rate
          error_rate=$(curl -s 'http://prometheus:9090/api/v1/query?query=rate(http_requests_total{status=~"5.."}[5m])' | jq '.data.result[0].value[1]')
          if (( $(echo "$error_rate > 0.01" | bc -l) )); then
            echo "✗ Error rate too high: $error_rate"
            exit 1
          fi

  switch-traffic:
    needs: e2e-post-deploy
    runs-on: ubuntu-latest
    steps:
      - name: Switch traffic to green
        run: |
          kubectl patch svc bgc-api -n bgc -p '{"spec":{"selector":{"version":"green"}}}'

      - name: Verify traffic switched
        run: |
          sleep 30
          curl -f https://api.bgc.com/healthz

  cleanup:
    needs: switch-traffic
    runs-on: ubuntu-latest
    steps:
      - name: Delete blue version
        run: |
          kubectl delete deployment bgc-api-blue -n bgc
```

### Definition of Done (Phase 5)

**Infrastructure:**
- [ ] Kubernetes cluster running in production
- [ ] 3 availability zones configured
- [ ] Auto-scaling enabled (2-10 pods)
- [ ] Load balancer provisioned

**CI/CD:**
- [ ] ArgoCD deployed and syncing
- [ ] Blue-green deployment automated
- [ ] Rollback tested and works in < 5 minutes
- [ ] Zero-downtime deployment verified

**E2E Testing:**
- [ ] Cluster health test passes
- [ ] ArgoCD sync test passes
- [ ] Blue-green deployment test passes
- [ ] Full deployment pipeline runs successfully
- [ ] Load test passes in production (100 concurrent users)

**Monitoring:**
- [ ] Prometheus scraping all pods
- [ ] Grafana dashboards accessible
- [ ] Alerts configured for critical metrics
- [ ] On-call rotation set up

**Documentation:**
- [ ] Runbook updated with K8s commands
- [ ] Deployment guide for new releases
- [ ] Rollback procedure documented
- [ ] Architecture diagram updated

### Risks & Mitigations

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Deployment fails in production | Low | Critical | Test in staging first, have rollback plan, use blue-green |
| Database migration breaks production | Low | Critical | Test on production snapshot, have rollback SQL, use transactions |
| Auto-scaling too aggressive (cost spike) | Medium | Medium | Set conservative limits, monitor costs daily |
| Traffic not switching to green | Low | High | Test traffic switching in staging, have manual override |

---

## CI/CD Pipeline Integration (All Phases)

### GitHub Actions Workflow

**File:** `.github/workflows/e2e-main.yml`

```yaml
name: E2E Test Suite

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]
  schedule:
    - cron: '0 */6 * * *' # Every 6 hours

jobs:
  backend-e2e:
    name: Backend E2E Tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Start infrastructure
        run: docker compose -f bgcstack/docker-compose.yml up -d

      - name: Wait for services
        run: |
          timeout 120 bash -c 'until curl -f http://localhost:8080/healthz; do sleep 2; done'
          timeout 120 bash -c 'until curl -f http://localhost:3000; do sleep 2; done'

      - name: Run backend E2E tests
        working-directory: ./api
        run: go test -tags=e2e -v -count=1 ./tests/e2e/...

      - name: Upload logs on failure
        if: failure()
        uses: actions/upload-artifact@v3
        with:
          name: backend-logs
          path: |
            /tmp/*.log
            bgcstack/logs/**

  frontend-e2e:
    name: Frontend E2E Tests (Playwright)
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Node
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install pnpm
        run: npm install -g pnpm

      - name: Install dependencies
        working-directory: ./web-next
        run: pnpm install

      - name: Install Playwright
        working-directory: ./web-next
        run: pnpm exec playwright install --with-deps

      - name: Start infrastructure
        run: docker compose -f bgcstack/docker-compose.yml up -d

      - name: Wait for services
        run: |
          timeout 120 bash -c 'until curl -f http://localhost:8080/healthz; do sleep 2; done'
          timeout 120 bash -c 'until curl -f http://localhost:3000; do sleep 2; done'

      - name: Run Playwright tests
        working-directory: ./web-next
        run: pnpm exec playwright test

      - name: Upload Playwright report
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: web-next/playwright-report/

  integration-e2e:
    name: Full Stack Integration Tests
    runs-on: ubuntu-latest
    needs: [backend-e2e, frontend-e2e]
    steps:
      - uses: actions/checkout@v4

      - name: Start full stack
        run: docker compose -f bgcstack/docker-compose.yml up -d

      - name: Wait for services
        run: |
          timeout 120 bash -c 'until curl -f http://localhost:8080/healthz; do sleep 2; done'
          timeout 120 bash -c 'until curl -f http://localhost:3000; do sleep 2; done'

      - name: Run integration tests
        run: bash tests/e2e/full_stack_integration.sh

      - name: Check Prometheus metrics
        run: |
          # Verify no error spikes
          curl -s 'http://localhost:9090/api/v1/query?query=rate(http_requests_total{status=~"5.."}[5m])' | jq '.data.result'

      - name: Check logs for errors
        run: |
          docker logs bgc_api 2>&1 | grep -i error || true
          docker logs bgc_web 2>&1 | grep -i error || true

  load-test:
    name: Load Test (100 concurrent users)
    runs-on: ubuntu-latest
    needs: integration-e2e
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4

      - name: Install k6
        run: |
          sudo gpg -k
          sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update
          sudo apt-get install k6

      - name: Start infrastructure
        run: docker compose -f bgcstack/docker-compose.yml up -d

      - name: Run load test
        run: k6 run tests/e2e/load_test.js --vus 100 --duration 5m

      - name: Check performance metrics
        run: |
          # P95 latency should be < 200ms
          p95=$(curl -s 'http://localhost:9090/api/v1/query?query=histogram_quantile(0.95,rate(http_request_duration_seconds_bucket[5m]))' | jq '.data.result[0].value[1]')
          if (( $(echo "$p95 > 0.2" | bc -l) )); then
            echo "✗ P95 latency too high: ${p95}s"
            exit 1
          fi
```

### Test Execution Frequency

| Test Type | Trigger | Frequency | Timeout |
|-----------|---------|-----------|---------|
| Smoke Tests | Every commit | On push to any branch | 5 min |
| Backend E2E | Every commit to main/develop | On push + PR | 15 min |
| Frontend E2E (Playwright) | Every commit to main/develop | On push + PR | 20 min |
| Integration E2E | After backend + frontend pass | On push to main/develop | 30 min |
| Load Test | After integration tests pass | On push to main only | 10 min |
| Full Suite | Scheduled | Every 6 hours | 60 min |
| Pre-deploy | Before production deploy | Manual trigger | 60 min |

---

## Metrics & Success Criteria (Overall)

### Technical Metrics

| Metric | Phase 0 | Phase 1 | Phase 2 | Phase 3 | Phase 4 | Phase 5 | Target |
|--------|---------|---------|---------|---------|---------|---------|--------|
| **E2E Test Coverage** | 0% | 30% | 50% | 70% | 85% | 95% | > 90% |
| **E2E Test Reliability** | N/A | 80% | 90% | 95% | 98% | 99% | > 95% |
| **Deployment Success Rate** | N/A | N/A | N/A | N/A | 90% | 98% | > 95% |
| **Mean Time to Detect (MTTD)** | N/A | 10min | 5min | 2min | 1min | 30s | < 2min |
| **Mean Time to Recovery (MTTR)** | N/A | 2h | 1h | 30min | 15min | 10min | < 30min |
| **API Latency P95** | N/A | 200ms | 180ms | 150ms | 120ms | 100ms | < 200ms |
| **Error Rate** | N/A | 1% | 0.5% | 0.3% | 0.1% | 0.05% | < 0.5% |
| **Uptime** | N/A | 99% | 99.5% | 99.7% | 99.8% | 99.9% | > 99.5% |

### Product Metrics

| Metric | Phase 0 | Phase 1 | Phase 2 | Phase 3 | Phase 4 | Phase 5 | Target |
|--------|---------|---------|---------|---------|---------|---------|--------|
| **Data Coverage** | 8% | 8% | 28% | 28% | 28% | 28% | 100% (v0.6.0) |
| **Active Users** | 0 | 0 | 0 | 0 | 20 | 50 | 100 (v1.0) |
| **NPS** | N/A | N/A | N/A | N/A | 40 | 50 | > 60 |
| **Simulation Completion Rate** | N/A | N/A | N/A | N/A | 70% | 80% | > 80% |
| **Free → Premium Conversion** | N/A | N/A | N/A | 0% | 3% | 5% | > 5% |

---

## Roadmap Summary Table

| Phase | Duration | Priority | Key Deliverables | E2E Tests | DoD Criteria |
|-------|----------|----------|------------------|-----------|--------------|
| **0: Cleanup** | 30 min | P0 | Clean repo, updated .gitignore | Smoke tests | Zero temp files, all services start |
| **1: Backend Ready** | 10h | P0 | Graceful shutdown, logging, migrations, health check, distributed rate limit | 7 backend E2E tests | 100% production-ready, all tests pass |
| **2: Data 28%** | 3-4 days | P0 | Chapters 02+08 populated, API validated, frontend updated | 6 NCM tests + load test | 28% coverage, < 200ms P95 |
| **3: Auth** | 2 days | P1 | Clerk integration, JWT validation, user tiers, rate limit by user | Auth flow + rate limit tests | Sign in works, premium unlimited |
| **4: Beta** | 2-3 weeks | P0 | 20 active users, feedback collected, NPS > 40 | Critical path + load test | Product-market fit validated |
| **5: Kubernetes** | 2 weeks (parallel) | P1 | Production K8s, ArgoCD, blue-green, zero-downtime | Full deployment pipeline | Production deployment automated |

---

## Conclusion

This strategic roadmap treats **E2E testing as a first-class citizen**, ensuring every phase is validated end-to-end before being marked complete. By integrating tests into CI/CD and defining clear Definition of Done criteria, we guarantee production-ready quality at every milestone.

**Key Principles:**
1. **No phase is "done" until all E2E tests pass**
2. **Tests run automatically on every commit**
3. **Performance targets are enforced, not optional**
4. **Regressions are caught immediately, not in production**

**Next Actions:**
1. Review and approve this roadmap
2. Start Phase 0 (cleanup) immediately
3. Set up CI/CD pipeline with E2E tests
4. Begin Phase 1 (backend production-ready)

---

**Version:** 3.0
**Last Updated:** 2026-01-16
**Next Review:** 2026-01-20
**Owner:** BGC Product Management + Engineering
