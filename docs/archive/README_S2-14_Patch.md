
### Sprint 2 — Testes E2E & Aceite (S2-14)

Siga o roteiro em `docs/Sprint2_E2E_Checklist.md` ou use os comandos abaixo (PowerShell, na raiz `bgc`):

```powershell
docker compose -f .\bgcstack\docker-compose.yml up -d
.\scripts\seed.ps1
curl.exe http://localhost:8080/healthz
start http://localhost:3000
start http://localhost:3000/routes.html
.\scripts\checks.ps1
```

Critérios de aceite:
- `/healthz` ok com `partner_weights:true` e `tariffs_loaded:true`
- `/market/size` retorna TAM/SAM/SOM; **SOM base vs aggressive** difere (~2x)
- `/routes/compare` aplica `tariff_scenario` e soma `results` = `adjusted_total_usd`
- Dashboard e Rotas funcionais com export CSV
- Checks de qualidade **Aprovado** (ou **Aprovado com alertas** apenas para outliers)
