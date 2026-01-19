# BGC E2E Testing Strategy - Execution Guide

**Version:** 1.0
**Date:** 2026-01-16
**Owner:** QA + Engineering
**Status:** Active

---

## Executive Summary

This document provides **executable guides** for implementing E2E tests across all BGC platform phases. Each section includes:
- Test suite specifications
- Automation scripts (copy-paste ready)
- Success metrics
- Troubleshooting guides

---

## Test Infrastructure Setup

### Prerequisites

**Tools Required:**
- Go 1.24+ (backend tests)
- Node.js 20+ (frontend tests)
- Playwright (browser automation)
- Docker + Docker Compose
- k6 (load testing)
- curl, jq (API testing)

**Installation Script:**
```bash
#!/bin/bash
# setup_e2e_infrastructure.sh

echo "Setting up E2E testing infrastructure..."

# Install Playwright
cd web-next
pnpm add -D @playwright/test
pnpm exec playwright install --with-deps

# Install k6 (Ubuntu/Debian)
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg \
  --keyserver hkp://keyserver.ubuntu.com:80 \
  --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | \
  sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6

# Install Go test dependencies
cd ../api
go install github.com/stretchr/testify/assert@latest
go install github.com/stretchr/testify/require@latest

echo "✓ E2E infrastructure setup complete"
```

---

## Test Suite Organization

### Directory Structure

```
bgc-app/
├── api/
│   └── tests/
│       ├── e2e/
│       │   ├── simulator_test.go
│       │   ├── backend_readiness_test.go
│       │   ├── auth_test.go
│       │   └── data_population_test.go
│       └── load/
│           └── load_test.go
├── web-next/
│   └── tests/
│       └── e2e/
│           ├── simulator.spec.ts
│           ├── auth.spec.ts
│           └── user_journey.spec.ts
├── tests/
│   └── e2e/
│       ├── smoke_tests.sh
│       ├── phase1_backend_ready.sh
│       ├── phase2_data_expansion.sh
│       ├── phase3_auth.sh
│       ├── phase4_beta.sh
│       └── full_stack_integration.sh
└── .github/
    └── workflows/
        ├── e2e-main.yml
        ├── backend-e2e.yml
        └── frontend-e2e.yml
```

---

## Phase-Specific Test Suites

### Phase 0: Cleanup Tests

**File:** `tests/e2e/phase0_cleanup.sh`

```bash
#!/bin/bash

set -e

echo "╔════════════════════════════════════════════════════════════╗"
echo "║  Phase 0: Cleanup Validation Tests                        ║"
echo "╚════════════════════════════════════════════════════════════╝"

PASSED=0
FAILED=0

# Test 1: No temp files in repository
echo "Test 1: Verifying no temporary files..."
temp_files=$(find . -name "tmpclaude-*" -o -name "nul" -o -name "NUL" 2>/dev/null | wc -l)
if [ "$temp_files" -eq 0 ]; then
    echo "✓ No temporary files found"
    ((PASSED++))
else
    echo "✗ Found $temp_files temporary files"
    ((FAILED++))
fi

# Test 2: .gitignore includes temp patterns
echo "Test 2: Verifying .gitignore updated..."
if grep -q "tmpclaude-" .gitignore && grep -q "nul" .gitignore; then
    echo "✓ .gitignore properly configured"
    ((PASSED++))
else
    echo "✗ .gitignore missing temp file patterns"
    ((FAILED++))
fi

# Test 3: Services start successfully
echo "Test 3: Starting services..."
docker compose -f bgcstack/docker-compose.yml up -d
sleep 10

# Test 4: Health checks pass
echo "Test 4: Checking health endpoints..."
if curl -f -s http://localhost:8080/healthz > /dev/null; then
    echo "✓ API health check passed"
    ((PASSED++))
else
    echo "✗ API health check failed"
    ((FAILED++))
fi

if curl -f -s http://localhost:3000 > /dev/null; then
    echo "✓ Frontend health check passed"
    ((PASSED++))
else
    echo "✗ Frontend health check failed"
    ((FAILED++))
fi

# Summary
echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║  Phase 0 Test Results                                     ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo "Passed: $PASSED"
echo "Failed: $FAILED"

if [ $FAILED -eq 0 ]; then
    echo "✓ Phase 0 cleanup validated - READY FOR PHASE 1"
    exit 0
else
    echo "✗ Phase 0 validation failed - FIX ISSUES BEFORE PROCEEDING"
    exit 1
fi
```

**Run Command:**
```bash
bash tests/e2e/phase0_cleanup.sh
```

---

### Phase 1: Backend Production-Ready Tests

**File:** `api/tests/e2e/backend_readiness_test.go`

```go
// +build e2e

package e2e

import (
    "context"
    "net/http"
    "os"
    "os/exec"
    "syscall"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

func TestHealthCheckComplete(t *testing.T) {
    resp, err := http.Get(baseURL + "/healthz")
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, 200, resp.StatusCode)

    var health map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&health)

    assert.Equal(t, "ok", health["status"])
    assert.Contains(t, health, "checks")

    checks := health["checks"].(map[string]interface{})
    assert.Equal(t, "healthy", checks["database"])
    assert.Equal(t, "healthy", checks["redis"])

    // Performance check
    assert.NotZero(t, health["timestamp"])
    assert.Less(t, time.Since(time.Unix(health["timestamp"].(int64), 0)), 2*time.Second)
}

func TestGracefulShutdown(t *testing.T) {
    // Start API in separate process
    cmd := exec.Command("go", "run", "../../cmd/api/main.go")
    err := cmd.Start()
    require.NoError(t, err)

    // Wait for startup
    time.Sleep(3 * time.Second)

    // Start long-running request
    done := make(chan bool)
    go func() {
        resp, _ := http.Get(baseURL + "/v1/simulator/destinations?ncm=17011400&volume_kg=1000")
        assert.Equal(t, 200, resp.StatusCode)
        done <- true
    }()

    time.Sleep(200 * time.Millisecond)

    // Send SIGTERM
    cmd.Process.Signal(syscall.SIGTERM)

    // Verify request completes
    select {
    case <-done:
        t.Log("✓ Request completed successfully during shutdown")
    case <-time.After(35 * time.Second):
        t.Fatal("✗ Request timed out - graceful shutdown failed")
    }

    cmd.Wait()
}

func TestStructuredLogging(t *testing.T) {
    // Make API request
    resp, err := http.Post(
        baseURL+"/v1/simulator/destinations",
        "application/json",
        bytes.NewBuffer([]byte(`{"ncm":"17011400","volume_kg":1000}`)),
    )
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, 200, resp.StatusCode)

    // Check logs (assumes logs written to file or stdout captured)
    // In production, query from centralized logging (e.g., ELK, Datadog)
    // For now, verify log format from Docker logs

    logs := exec.Command("docker", "logs", "bgc_api", "--tail", "10")
    output, _ := logs.CombinedOutput()

    // Verify JSON log format
    assert.Contains(t, string(output), `"level"`)
    assert.Contains(t, string(output), `"timestamp"`)
    assert.Contains(t, string(output), `"msg"`)
    assert.Contains(t, string(output), `"ncm"`)
}

func TestDatabaseMigrations(t *testing.T) {
    // Run migrations
    cmd := exec.Command("go", "run", "../../cmd/migrate/main.go")
    cmd.Env = append(os.Environ(), "DATABASE_URL="+os.Getenv("TEST_DATABASE_URL"))
    output, err := cmd.CombinedOutput()
    require.NoError(t, err, "Migration failed: %s", output)

    // Verify schema
    db := setupTestDB()
    defer db.Close()

    // Check critical tables exist
    tables := []string{"users", "countries_metadata", "stg.exportacao"}
    for _, table := range tables {
        var exists bool
        err := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE tablename = $1)", table).Scan(&exists)
        require.NoError(t, err)
        assert.True(t, exists, "Table %s should exist", table)
    }

    // Run migrations again (idempotency test)
    cmd = exec.Command("go", "run", "../../cmd/migrate/main.go")
    cmd.Env = append(os.Environ(), "DATABASE_URL="+os.Getenv("TEST_DATABASE_URL"))
    output, err = cmd.CombinedOutput()
    require.NoError(t, err, "Idempotent migration failed: %s", output)
}

func TestDistributedRateLimiting(t *testing.T) {
    // Start 3 API instances
    ports := []int{8080, 8081, 8082}
    for i, port := range ports {
        cmd := exec.Command("go", "run", "../../cmd/api/main.go")
        cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", port))
        cmd.Start()
        defer cmd.Process.Kill()
        if i == 0 {
            time.Sleep(3 * time.Second) // Wait for first instance
        } else {
            time.Sleep(1 * time.Second)
        }
    }

    client := &http.Client{}

    // Make 5 requests distributed across instances
    for i := 0; i < 5; i++ {
        port := ports[i%3]
        req, _ := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/v1/simulator/destinations", port), bytes.NewBuffer([]byte(`{"ncm":"17011400"}`)))
        req.Header.Set("Content-Type", "application/json")
        resp, err := client.Do(req)
        require.NoError(t, err)
        assert.Equal(t, 200, resp.StatusCode, "Request %d should succeed", i+1)
        resp.Body.Close()
    }

    // 6th request should be rate limited on ANY instance
    port := ports[2] // Try third instance
    req, _ := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/v1/simulator/destinations", port), bytes.NewBuffer([]byte(`{"ncm":"17011400"}`)))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := client.Do(req)
    assert.Equal(t, 429, resp.StatusCode, "Request should be rate limited across instances")
}

func TestPerformanceTarget(t *testing.T) {
    // Warmup
    for i := 0; i < 5; i++ {
        http.Post(baseURL+"/v1/simulator/destinations", "application/json", bytes.NewBuffer([]byte(`{"ncm":"17011400","volume_kg":1000}`)))
        time.Sleep(100 * time.Millisecond)
    }

    // Measure P95 latency over 20 requests
    var latencies []time.Duration
    for i := 0; i < 20; i++ {
        start := time.Now()
        resp, err := http.Post(baseURL+"/v1/simulator/destinations", "application/json", bytes.NewBuffer([]byte(`{"ncm":"17011400","volume_kg":1000}`)))
        latency := time.Since(start)
        require.NoError(t, err)
        resp.Body.Close()

        latencies = append(latencies, latency)
        time.Sleep(50 * time.Millisecond)
    }

    // Sort and get P95
    sort.Slice(latencies, func(i, j int) bool {
        return latencies[i] < latencies[j]
    })
    p95 := latencies[int(float64(len(latencies))*0.95)]

    t.Logf("P95 latency: %v", p95)
    assert.Less(t, p95, 200*time.Millisecond, "P95 latency should be < 200ms")
}
```

**Run Command:**
```bash
cd api
go test -tags=e2e -v -timeout 10m ./tests/e2e/backend_readiness_test.go
```

---

### Phase 2: Data Expansion Tests

**File:** `tests/e2e/phase2_data_expansion.sh`

```bash
#!/bin/bash

set -e

echo "╔════════════════════════════════════════════════════════════╗"
echo "║  Phase 2: Data Expansion (8% → 28%) Tests                 ║"
echo "╚════════════════════════════════════════════════════════════╝"

PASSED=0
FAILED=0

# Database connection
export PGPASSWORD=bgc
PSQL="docker exec bgc_db psql -U bgc -d bgc -t -c"

# Test 1: Verify chapter 17 data still exists (regression)
echo "Test 1: Chapter 17 regression check..."
count_17=$($PSQL "SELECT COUNT(*) FROM stg.exportacao WHERE SUBSTRING(co_ncm, 1, 2) = '17'" | xargs)
if [ "$count_17" -gt 50 ]; then
    echo "✓ Chapter 17 has $count_17 records (maintained)"
    ((PASSED++))
else
    echo "✗ Chapter 17 has insufficient data: $count_17"
    ((FAILED++))
fi

# Test 2: Verify chapter 02 data populated
echo "Test 2: Chapter 02 (Meats) data check..."
count_02=$($PSQL "SELECT COUNT(*) FROM stg.exportacao WHERE SUBSTRING(co_ncm, 1, 2) = '02'" | xargs)
if [ "$count_02" -gt 100 ]; then
    echo "✓ Chapter 02 has $count_02 records"
    ((PASSED++))
else
    echo "✗ Chapter 02 has insufficient data: $count_02 (expected > 100)"
    ((FAILED++))
fi

# Test 3: Verify chapter 08 data populated
echo "Test 3: Chapter 08 (Fruits) data check..."
count_08=$($PSQL "SELECT COUNT(*) FROM stg.exportacao WHERE SUBSTRING(co_ncm, 1, 2) = '08'" | xargs)
if [ "$count_08" -gt 100 ]; then
    echo "✓ Chapter 08 has $count_08 records"
    ((PASSED++))
else
    echo "✗ Chapter 08 has insufficient data: $count_08 (expected > 100)"
    ((FAILED++))
fi

# Test 4: API endpoint works for chapter 02
echo "Test 4: API endpoint - Chapter 02 (Beef)..."
response=$(curl -s -X POST http://localhost:8080/v1/simulator/destinations \
    -H "Content-Type: application/json" \
    -d '{"ncm":"02013000","volume_kg":10000}')

destinations=$(echo "$response" | jq -r '.destinations | length')
if [ "$destinations" -ge 3 ]; then
    echo "✓ Chapter 02 API returns $destinations destinations"
    ((PASSED++))
else
    echo "✗ Chapter 02 API returned insufficient destinations: $destinations"
    ((FAILED++))
fi

# Test 5: API endpoint works for chapter 08
echo "Test 5: API endpoint - Chapter 08 (Coconuts)..."
response=$(curl -s -X POST http://localhost:8080/v1/simulator/destinations \
    -H "Content-Type: application/json" \
    -d '{"ncm":"08011100","volume_kg":5000}')

destinations=$(echo "$response" | jq -r '.destinations | length')
if [ "$destinations" -ge 3 ]; then
    echo "✓ Chapter 08 API returns $destinations destinations"
    ((PASSED++))
else
    echo "✗ Chapter 08 API returned insufficient destinations: $destinations"
    ((FAILED++))
fi

# Test 6: Performance check (< 200ms P95)
echo "Test 6: Performance validation..."
total_time=0
max_time=0
for i in {1..10}; do
    start=$(date +%s%3N)
    curl -s -X POST http://localhost:8080/v1/simulator/destinations \
        -H "Content-Type: application/json" \
        -d '{"ncm":"02013000","volume_kg":10000}' > /dev/null
    end=$(date +%s%3N)
    time=$((end - start))
    total_time=$((total_time + time))
    if [ $time -gt $max_time ]; then
        max_time=$time
    fi
done
avg_time=$((total_time / 10))

if [ $max_time -lt 200 ]; then
    echo "✓ Performance target met (avg: ${avg_time}ms, max: ${max_time}ms)"
    ((PASSED++))
else
    echo "✗ Performance target missed (max: ${max_time}ms > 200ms)"
    ((FAILED++))
fi

# Test 7: Data quality - no duplicates
echo "Test 7: Data quality check..."
duplicates=$($PSQL "
    SELECT COUNT(*) FROM (
        SELECT co_ncm, co_pais, co_ano, co_mes, COUNT(*)
        FROM stg.exportacao
        WHERE SUBSTRING(co_ncm, 1, 2) IN ('02', '08')
        GROUP BY co_ncm, co_pais, co_ano, co_mes
        HAVING COUNT(*) > 1
    ) dup
" | xargs)

if [ "$duplicates" -eq 0 ]; then
    echo "✓ No duplicate records found"
    ((PASSED++))
else
    echo "✗ Found $duplicates duplicate records"
    ((FAILED++))
fi

# Summary
echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║  Phase 2 Test Results                                     ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo "Passed: $PASSED"
echo "Failed: $FAILED"
echo ""
echo "Data Coverage Summary:"
echo "  Chapter 17 (Sugars):  $count_17 records"
echo "  Chapter 02 (Meats):   $count_02 records"
echo "  Chapter 08 (Fruits):  $count_08 records"
total_records=$((count_17 + count_02 + count_08))
echo "  Total:                $total_records records"

if [ $FAILED -eq 0 ]; then
    echo ""
    echo "✓ Phase 2 data expansion validated - 28% COVERAGE ACHIEVED"
    exit 0
else
    echo ""
    echo "✗ Phase 2 validation failed - FIX ISSUES BEFORE PROCEEDING"
    exit 1
fi
```

**Run Command:**
```bash
bash tests/e2e/phase2_data_expansion.sh
```

---

### Phase 3: Authentication Tests (Playwright)

**File:** `web-next/tests/e2e/auth.spec.ts`

```typescript
import { test, expect } from '@playwright/test';

test.describe('Authentication Flow - Phase 3', () => {
  test.beforeEach(async ({ page }) => {
    // Reset test database
    await page.goto('http://localhost:3000/api/test/reset-db');
  });

  test('should show sign in button when not authenticated', async ({ page }) => {
    await page.goto('http://localhost:3000');

    const signInButton = page.locator('button:has-text("Login")');
    await expect(signInButton).toBeVisible();
  });

  test('should complete full sign up flow', async ({ page }) => {
    await page.goto('http://localhost:3000');

    // Click sign up
    await page.click('button:has-text("Login")');
    await page.click('text=Criar Conta');

    // Fill sign up form
    await page.fill('input[name="email"]', `test-${Date.now()}@example.com`);
    await page.fill('input[name="password"]', 'SecurePassword123!');
    await page.fill('input[name="firstName"]', 'Test');
    await page.fill('input[name="lastName"]', 'User');

    await page.click('button:has-text("Criar Conta")');

    // Verify signed in
    await page.waitForSelector('[data-testid="user-button"]', { timeout: 10000 });
    await expect(page.locator('[data-testid="user-button"]')).toBeVisible();
  });

  test('should rate limit free users correctly', async ({ page }) => {
    // Sign in as free user
    await signIn(page, 'free-user@example.com', 'FreeUser123!');

    // Use all 5 free simulations
    for (let i = 0; i < 5; i++) {
      await page.goto('http://localhost:3000/simulator');
      await page.fill('input[name="ncm"]', '17011400');
      await page.fill('input[name="volume_kg"]', '1000');
      await page.click('button:has-text("Simular")');

      await page.waitForSelector('[data-testid="destination-card"]');

      // Verify rate limit header updates
      const banner = page.locator('[data-testid="rate-limit-banner"]');
      await expect(banner).toContainText(`${4 - i} de 5`);
    }

    // 6th attempt should show upgrade modal
    await page.goto('http://localhost:3000/simulator');
    await page.fill('input[name="ncm"]', '17011400');
    await page.click('button:has-text("Simular")');

    await page.waitForSelector('[data-testid="upgrade-modal"]');
    await expect(page.locator('text=Limite Atingido')).toBeVisible();
  });

  test('should allow unlimited simulations for premium users', async ({ page }) => {
    // Sign in as premium user
    await signIn(page, 'premium-user@example.com', 'PremiumUser123!');

    // Make 10 simulations (should all succeed)
    for (let i = 0; i < 10; i++) {
      await page.goto('http://localhost:3000/simulator');
      await page.fill('input[name="ncm"]', '17011400');
      await page.fill('input[name="volume_kg"]', '1000');
      await page.click('button:has-text("Simular")');

      await page.waitForSelector('[data-testid="destination-card"]');

      // Verify NO rate limit banner for premium
      const banner = page.locator('[data-testid="rate-limit-banner"]');
      await expect(banner).not.toBeVisible();
    }
  });

  test('should protect dashboard routes', async ({ page }) => {
    // Try to access protected route without auth
    await page.goto('http://localhost:3000/dashboard');

    // Should redirect to sign in
    await page.waitForURL(/sign-in/);
    await expect(page).toHaveURL(/sign-in/);
  });

  test('should allow access to dashboard after sign in', async ({ page }) => {
    await signIn(page, 'test-user@example.com', 'TestUser123!');

    await page.goto('http://localhost:3000/dashboard');

    // Should load dashboard (not redirect)
    await expect(page).toHaveURL(/dashboard/);
    await expect(page.locator('h1:has-text("Dashboard")')).toBeVisible();
  });

  test('should sign out successfully', async ({ page }) => {
    await signIn(page, 'test-user@example.com', 'TestUser123!');

    // Click user button
    await page.click('[data-testid="user-button"]');

    // Click sign out
    await page.click('button:has-text("Sign out")');

    // Verify signed out
    await page.waitForSelector('button:has-text("Login")');
    await expect(page.locator('button:has-text("Login")')).toBeVisible();

    // Verify cannot access protected routes
    await page.goto('http://localhost:3000/dashboard');
    await page.waitForURL(/sign-in/);
  });
});

// Helper function
async function signIn(page, email: string, password: string) {
  await page.goto('http://localhost:3000');
  await page.click('button:has-text("Login")');

  await page.waitForSelector('input[name="identifier"]');
  await page.fill('input[name="identifier"]', email);
  await page.click('button:has-text("Continue")');

  await page.waitForSelector('input[name="password"]');
  await page.fill('input[name="password"]', password);
  await page.click('button:has-text("Continue")');

  await page.waitForSelector('[data-testid="user-button"]', { timeout: 10000 });
}
```

**Run Command:**
```bash
cd web-next
pnpm exec playwright test tests/e2e/auth.spec.ts --headed
```

---

## Load Testing (k6)

**File:** `tests/e2e/load_test.js`

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '1m', target: 10 },   // Ramp up to 10 users
    { duration: '3m', target: 50 },   // Ramp up to 50 users
    { duration: '5m', target: 100 },  // Ramp up to 100 users
    { duration: '2m', target: 100 },  // Stay at 100 users
    { duration: '2m', target: 0 },    // Ramp down
  ],
  thresholds: {
    'http_req_duration': ['p(95)<200'], // 95% of requests under 200ms
    'http_req_failed': ['rate<0.01'],   // Less than 1% errors
    'errors': ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

const NCMs = [
  '17011400', // Sugar
  '02013000', // Beef
  '08011100', // Coconuts
];

export default function () {
  // Select random NCM
  const ncm = NCMs[Math.floor(Math.random() * NCMs.length)];

  const payload = JSON.stringify({
    ncm: ncm,
    volume_kg: Math.floor(Math.random() * 50000) + 1000,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${__ENV.TEST_TOKEN}`,
    },
  };

  // Make request
  const res = http.post(`${BASE_URL}/v1/simulator/destinations`, payload, params);

  // Validate response
  const success = check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
    'has destinations': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.destinations && body.destinations.length > 0;
      } catch {
        return false;
      }
    },
  });

  errorRate.add(!success);

  // Realistic user behavior: wait 3-10 seconds between requests
  sleep(Math.random() * 7 + 3);
}

export function handleSummary(data) {
  return {
    'stdout': textSummary(data, { indent: ' ', enableColors: true }),
    'load-test-report.json': JSON.stringify(data),
  };
}
```

**Run Command:**
```bash
# Phase 2 load test (30 users)
k6 run tests/e2e/load_test.js --vus 30 --duration 5m

# Phase 4 beta load test (50 users)
k6 run tests/e2e/load_test.js --vus 50 --duration 10m

# Phase 5 production load test (100 users)
k6 run tests/e2e/load_test.js --vus 100 --duration 15m
```

---

## Smoke Test Suite (Quick Validation)

**File:** `tests/e2e/smoke_tests.sh`

```bash
#!/bin/bash

# Quick smoke tests (run after every deployment)
# Should complete in < 2 minutes

set -e

echo "Running smoke tests..."

BASE_URL=${1:-http://localhost:8080}
FRONTEND_URL=${2:-http://localhost:3000}

PASSED=0
FAILED=0

# Test 1: API health
if curl -f -s "${BASE_URL}/healthz" > /dev/null; then
    echo "✓ API health check passed"
    ((PASSED++))
else
    echo "✗ API health check failed"
    ((FAILED++))
fi

# Test 2: Frontend loads
if curl -f -s "${FRONTEND_URL}" > /dev/null; then
    echo "✓ Frontend loads"
    ((PASSED++))
else
    echo "✗ Frontend failed to load"
    ((FAILED++))
fi

# Test 3: Simulator endpoint works
response=$(curl -s -X POST "${BASE_URL}/v1/simulator/destinations" \
    -H "Content-Type: application/json" \
    -d '{"ncm":"17011400","volume_kg":1000}')

if echo "$response" | jq -e '.destinations | length > 0' > /dev/null; then
    echo "✓ Simulator endpoint working"
    ((PASSED++))
else
    echo "✗ Simulator endpoint failed"
    ((FAILED++))
fi

# Test 4: Database connectivity
if docker exec bgc_db psql -U bgc -d bgc -c "SELECT 1" > /dev/null 2>&1; then
    echo "✓ Database connected"
    ((PASSED++))
else
    echo "✗ Database connection failed"
    ((FAILED++))
fi

# Test 5: Redis connectivity
if docker exec bgc_redis redis-cli ping | grep -q PONG; then
    echo "✓ Redis connected"
    ((PASSED++))
else
    echo "✗ Redis connection failed"
    ((FAILED++))
fi

echo ""
echo "Smoke Test Summary: $PASSED passed, $FAILED failed"

if [ $FAILED -eq 0 ]; then
    echo "✓ All smoke tests passed"
    exit 0
else
    echo "✗ Smoke tests failed"
    exit 1
fi
```

**Run Command:**
```bash
# Local
bash tests/e2e/smoke_tests.sh

# Staging
bash tests/e2e/smoke_tests.sh https://staging-api.bgc.com https://staging.bgc.com

# Production
bash tests/e2e/smoke_tests.sh https://api.bgc.com https://bgc.com
```

---

## Test Execution Cheat Sheet

### Local Development

```bash
# Start infrastructure
docker compose -f bgcstack/docker-compose.yml up -d

# Run all E2E tests
./run_all_e2e.sh

# Run specific phase tests
bash tests/e2e/phase0_cleanup.sh
bash tests/e2e/phase1_backend_ready.sh
bash tests/e2e/phase2_data_expansion.sh

# Run Go backend tests
cd api && go test -tags=e2e -v ./tests/e2e/...

# Run Playwright frontend tests
cd web-next && pnpm exec playwright test

# Run load test
k6 run tests/e2e/load_test.js --vus 50 --duration 5m

# Run smoke tests
bash tests/e2e/smoke_tests.sh
```

### CI/CD

```bash
# Trigger E2E suite in GitHub Actions
git push origin main

# Check status
gh run list --workflow=e2e-main.yml

# View logs
gh run view <run-id> --log

# Download artifacts
gh run download <run-id>
```

### Troubleshooting

```bash
# Check service logs
docker logs bgc_api --tail 100
docker logs bgc_db --tail 100
docker logs bgc_redis --tail 100

# Check Prometheus metrics
curl http://localhost:9090/api/v1/query?query=rate(http_requests_total[5m])

# Check database state
docker exec bgc_db psql -U bgc -d bgc -c "SELECT * FROM stg.exportacao LIMIT 10"

# Restart services
docker compose -f bgcstack/docker-compose.yml restart

# Clean slate
docker compose -f bgcstack/docker-compose.yml down -v
docker compose -f bgcstack/docker-compose.yml up -d
```

---

## Success Metrics Dashboard

### Key Metrics to Track

| Metric | Phase 1 Target | Phase 2 Target | Phase 4 Target | How to Measure |
|--------|---------------|---------------|---------------|----------------|
| **E2E Pass Rate** | > 90% | > 95% | > 98% | `(passed_tests / total_tests) * 100` |
| **Test Execution Time** | < 15 min | < 20 min | < 30 min | CI/CD pipeline duration |
| **Flaky Test Rate** | < 5% | < 3% | < 1% | Failed tests / runs without code changes |
| **API P95 Latency** | < 200ms | < 200ms | < 150ms | Load test report |
| **Error Rate** | < 1% | < 0.5% | < 0.1% | `rate(http_requests_total{status=~"5.."}[5m])` |
| **Test Coverage** | 30% | 50% | 85% | Go test coverage + Playwright coverage |

### Monitoring Queries (Prometheus)

```promql
# API error rate
rate(http_requests_total{status=~"5.."}[5m])

# P95 latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Request throughput
sum(rate(http_requests_total[5m]))

# Database query time
histogram_quantile(0.95, rate(db_query_duration_seconds_bucket[5m]))

# Rate limit hit rate
rate(rate_limit_exceeded_total[5m])
```

---

## Next Actions

1. **Today (Phase 0):** Run cleanup tests, fix issues, commit
2. **Tomorrow (Phase 1):** Implement graceful shutdown, run backend readiness tests
3. **Day 3 (Phase 2):** Populate data for chapters 02 and 08, run data expansion tests
4. **Day 5 (Phase 3):** Integrate Clerk auth, run auth tests
5. **Week 2 (Phase 4):** Launch beta, run user journey tests
6. **Week 3 (Phase 5):** Deploy to Kubernetes, run production tests

---

**Version:** 1.0
**Last Updated:** 2026-01-16
**Owner:** QA + Engineering Teams
**Status:** Active - Execute Phase 0 Today
