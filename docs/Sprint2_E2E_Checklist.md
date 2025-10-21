
# Sprint 2 — Testes E2E & Aceite (S2-14)

Este roteiro valida ponta-a-ponta a Onda 1 (Sprint 2): ingestão → API → Front → qualidade.

## Pré-requisitos
- Containers ativos: `db`, `pgadmin`, `api`, `web`
- Banco populado via `scripts/seed.ps1`
- CORS ativo na API

Comandos úteis (PowerShell, na raiz `bgc`):
```powershell
docker compose -f .\bgcstack\docker-compose.yml up -d
.\scripts\seed.ps1
```

---

## 1) API — Health & Market

### Health
```powershell
curl.exe http://localhost:8080/healthz
```
**Espera:** `ok:true`, `partner_weights:true`, `tariffs_loaded:true`, `available_scenarios` (ex.: `["base","tarifa10"]`).

### TAM/SAM/SOM
```powershell
# TAM 2024–2025 (agregado)
curl.exe "http://localhost:8080/market/size?metric=TAM&year_from=2024&year_to=2025"

# SOM base vs aggressive (deve ~dobrar 1,5% -> 3%)
curl.exe "http://localhost:8080/market/size?metric=SOM&year_from=2024&year_to=2024&scenario=base"
curl.exe "http://localhost:8080/market/size?metric=SOM&year_from=2024&year_to=2024&scenario=aggressive"
```

---

## 2) API — Rotas EUA vs Alternativos
```powershell
# Sem tarifa
curl.exe "http://localhost:8080/routes/compare?from=USA&alts=CHN,ARE,IND&ncm_chapter=84&year=2024"

# Com tarifa (ex.: tarifa10)
curl.exe "http://localhost:8080/routes/compare?from=USA&alts=CHN,ARE,IND&ncm_chapter=84&year=2024&tariff_scenario=tarifa10"
```
**Espera:** `adjusted_total_usd` muda conforme `factor`. `results[].factor` exibido por parceiro. Soma de `results[].estimated_usd` = `adjusted_total_usd` (tolerância de arredondamento).

---

## 3) Front — Dashboard (S2-10)
Abrir `http://localhost:3000` e testar:
- **TAM 2020–2025**: tabela por ano e KPIs preenchem.
- **SOM**: trocar **base/aggressive** e clicar **Consultar** → valores diferem (~2x).
- **Exportar CSV** gera `bgc_<metric>.csv` com `ano,metric,valor_usd`.
- Campo **Cenário** só habilita quando **Métrica = SOM**.

---

## 4) Front — Rotas (S2-11)
Abrir `http://localhost:3000/routes` e testar:
- Inputs: **Ano, Capítulo, Parceiro principal, Alternativos (abaixo), Cenário**.
- Tabela + **gráfico de barras** carregam com `results`.
- KPI **Checagem de soma** = **OK**.
- Export CSV com `partner,share,factor,estimated_usd`.

---

## 5) Qualidade de dados (S2-15)
```powershell
.\scripts\checks.ps1
```
**Espera atual:** Órfãos/Nulos/Negativos/Duplicados = 0. Outliers podem existir (ALERTA).

---

## 6) Encerramento (DoD)
```powershell
git add .
git commit -m "Sprint 2: DoD (E2E & checks)"
# opcional: git tag s2-done
```
Pronto: S2-14 concluída com evidências reprodutíveis.
