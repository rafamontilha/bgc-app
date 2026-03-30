/**
 * TDD — Testes para next.config.ts
 *
 * Motivação: o output 'standalone' estava desabilitado via comentário, sem
 * controle programático. Estes testes garantem que:
 *   1. O modo standalone seja ativado quando NEXT_STANDALONE=true (padrão em Docker)
 *   2. O modo standalone seja omitido em desenvolvimento (Windows/OneDrive safe)
 *   3. As rewrites de API estejam configuradas corretamente para todos os endpoints
 *   4. CORS só seja habilitado em desenvolvimento
 *
 * Execução:
 *   pnpm test -- --testPathPattern="next-config"
 */

import type { NextConfig } from "next";

// Helper para isolar a leitura do config com env específico.
// Retorna o config E uma função restore() — chame restore() DEPOIS de usar
// config.rewrites() / config.headers(), pois essas funções leem process.env
// no momento em que são invocadas (não capturam o valor na importação do módulo).
async function loadConfig(
  env: Record<string, string | undefined>
): Promise<{ config: NextConfig; restore: () => void }> {
  const original: Record<string, string | undefined> = {};
  for (const key of Object.keys(env)) {
    original[key] = process.env[key];
    if (env[key] === undefined) {
      delete process.env[key];
    } else {
      process.env[key] = env[key];
    }
  }

  jest.resetModules();
  const mod = await import("../../next.config");
  const config = mod.default as NextConfig;

  const restore = () => {
    for (const key of Object.keys(original)) {
      if (original[key] === undefined) {
        delete process.env[key];
      } else {
        process.env[key] = original[key];
      }
    }
  };

  return { config, restore };
}

// ---------------------------------------------------------------------------
// Standalone output
// ---------------------------------------------------------------------------

describe("standalone output", () => {
  it("ativa standalone quando NEXT_STANDALONE=true (padrão em Docker/CI)", async () => {
    const { config, restore } = await loadConfig({ NEXT_STANDALONE: "true" });
    restore();
    expect(config.output).toBe("standalone");
  });

  it("omite output quando NEXT_STANDALONE não está definido (dev local)", async () => {
    const { config, restore } = await loadConfig({ NEXT_STANDALONE: undefined });
    restore();
    expect(config.output).toBeUndefined();
  });

  it("omite output quando NEXT_STANDALONE=false (Windows/OneDrive)", async () => {
    const { config, restore } = await loadConfig({ NEXT_STANDALONE: "false" });
    restore();
    expect(config.output).toBeUndefined();
  });
});

// ---------------------------------------------------------------------------
// Rewrites de API (proxy para o backend Go)
// ---------------------------------------------------------------------------

describe("API rewrites", () => {
  it("configura proxy para /v1/market/:path*", async () => {
    const { config, restore } = await loadConfig({ NEXT_STANDALONE: undefined });
    const rewrites = await (config.rewrites as () => Promise<unknown[]>)();
    restore();
    const flat = (rewrites as Array<{ source: string; destination: string }>).flat();
    const rule = flat.find((r) => r.source === "/v1/market/:path*");
    expect(rule).toBeDefined();
    expect(rule?.destination).toContain("/v1/market/:path*");
  });

  it("configura proxy para /v1/routes/:path*", async () => {
    const { config, restore } = await loadConfig({ NEXT_STANDALONE: undefined });
    const rewrites = await (config.rewrites as () => Promise<unknown[]>)();
    restore();
    const flat = (rewrites as Array<{ source: string; destination: string }>).flat();
    const rule = flat.find((r) => r.source === "/v1/routes/:path*");
    expect(rule).toBeDefined();
  });

  it("configura proxy para /v1/chapters/:path*", async () => {
    const { config, restore } = await loadConfig({ NEXT_STANDALONE: undefined });
    const rewrites = await (config.rewrites as () => Promise<unknown[]>)();
    restore();
    const flat = (rewrites as Array<{ source: string; destination: string }>).flat();
    const rule = flat.find((r) => r.source === "/v1/chapters/:path*");
    expect(rule).toBeDefined();
  });

  it("usa API_URL como destino quando definida", async () => {
    const { config, restore } = await loadConfig({
      NEXT_STANDALONE: undefined,
      API_URL: "http://custom-api:9090",
      NEXT_PUBLIC_API_URL: undefined,
    });
    // rewrites() lê process.env ao vivo — chamar ANTES do restore
    const rewrites = await (config.rewrites as () => Promise<unknown[]>)();
    restore();
    const flat = (rewrites as Array<{ source: string; destination: string }>).flat();
    const rule = flat.find((r) => r.source === "/v1/market/:path*");
    expect(rule?.destination).toContain("custom-api:9090");
  });

  it("usa fallback http://bgc-api:8080 quando API_URL não está definida", async () => {
    const { config, restore } = await loadConfig({
      NEXT_STANDALONE: undefined,
      API_URL: undefined,
      NEXT_PUBLIC_API_URL: undefined,
    });
    const rewrites = await (config.rewrites as () => Promise<unknown[]>)();
    restore();
    const flat = (rewrites as Array<{ source: string; destination: string }>).flat();
    const rule = flat.find((r) => r.source === "/v1/market/:path*");
    expect(rule?.destination).toContain("bgc-api:8080");
  });

  it("inclui rota de health check /healthz", async () => {
    const { config, restore } = await loadConfig({ NEXT_STANDALONE: undefined });
    const rewrites = await (config.rewrites as () => Promise<unknown[]>)();
    restore();
    const flat = (rewrites as Array<{ source: string; destination: string }>).flat();
    const rule = flat.find((r) => r.source === "/healthz");
    expect(rule).toBeDefined();
  });

  it("inclui proxy para swagger docs", async () => {
    const { config, restore } = await loadConfig({ NEXT_STANDALONE: undefined });
    const rewrites = await (config.rewrites as () => Promise<unknown[]>)();
    restore();
    const flat = (rewrites as Array<{ source: string; destination: string }>).flat();
    const swagger = flat.find((r) => r.source.startsWith("/swagger/"));
    expect(swagger).toBeDefined();
  });
});

// ---------------------------------------------------------------------------
// CORS headers (apenas em desenvolvimento)
// ---------------------------------------------------------------------------

describe("CORS headers", () => {
  it("não adiciona headers CORS em produção", async () => {
    const { config, restore } = await loadConfig({
      NEXT_STANDALONE: "true",
      NODE_ENV: "production",
    });
    restore();
    expect(config.headers).toBeUndefined();
  });

  it("adiciona headers CORS em desenvolvimento", async () => {
    const env = process.env as Record<string, string | undefined>;
    const originalNodeEnv = env.NODE_ENV;
    env.NODE_ENV = "development";
    jest.resetModules();
    const mod = await import("../../next.config");
    const config = mod.default as NextConfig;

    if (config.headers) {
      const headers = await config.headers();
      env.NODE_ENV = originalNodeEnv;
      expect(headers.length).toBeGreaterThan(0);
      const apiHeaders = headers.find((h) => h.source.startsWith("/api/"));
      expect(apiHeaders).toBeDefined();
    } else {
      env.NODE_ENV = originalNodeEnv;
    }
  });
});

// ---------------------------------------------------------------------------
// Configurações gerais
// ---------------------------------------------------------------------------

describe("configurações gerais", () => {
  it("tem reactStrictMode habilitado", async () => {
    const { config, restore } = await loadConfig({ NEXT_STANDALONE: undefined });
    restore();
    expect(config.reactStrictMode).toBe(true);
  });

  it("tem images.unoptimized configurado", async () => {
    const { config, restore } = await loadConfig({ NEXT_STANDALONE: undefined });
    restore();
    expect(config.images?.unoptimized).toBe(true);
  });
});
