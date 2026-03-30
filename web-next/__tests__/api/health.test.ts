/**
 * @jest-environment node
 *
 * Route handlers Next.js são código de servidor; precisam do ambiente Node (tem
 * Response/Request nativos no Node 18+). jsdom não os define.
 *
 * TDD — Testes para o endpoint /api/health
 *
 * Motivação: o health check é usado por:
 *   - Kubernetes readinessProbe e livenessProbe (web-next.yaml)
 *   - Docker HEALTHCHECK no Dockerfile
 *   - Monitoring do uptime da aplicação
 *
 * Uma falha silenciosa neste endpoint causa:
 *   - Pod nunca fica Ready → serviço fora do ar
 *   - Rolling update trava indefinidamente
 *
 * Execução:
 *   pnpm test -- --testPathPattern="health"
 */

import { GET } from "../../app/api/health/route";

describe("GET /api/health", () => {
  it("retorna status 200", async () => {
    const response = await GET();
    expect(response.status).toBe(200);
  });

  it("retorna Content-Type application/json", async () => {
    const response = await GET();
    const contentType = response.headers.get("content-type");
    expect(contentType).toContain("application/json");
  });

  it("retorna campo status='ok'", async () => {
    const response = await GET();
    const body = await response.json();
    expect(body.status).toBe("ok");
  });

  it("retorna campo timestamp com ISO 8601 válido", async () => {
    const response = await GET();
    const body = await response.json();
    expect(body.timestamp).toBeDefined();
    const parsed = new Date(body.timestamp);
    expect(parsed.toISOString()).toBe(body.timestamp);
  });

  it("retorna timestamp próximo ao momento atual (±5s)", async () => {
    const before = Date.now();
    const response = await GET();
    const after = Date.now();
    const body = await response.json();
    const ts = new Date(body.timestamp).getTime();
    expect(ts).toBeGreaterThanOrEqual(before - 5000);
    expect(ts).toBeLessThanOrEqual(after + 5000);
  });

  it("responde de forma síncrona (não depende de serviços externos)", async () => {
    // O health check deve ser self-contained; não pode chamar o banco ou API externa.
    // Verificamos que resolve em < 100ms
    const start = performance.now();
    await GET();
    const elapsed = performance.now() - start;
    expect(elapsed).toBeLessThan(100);
  });
});
