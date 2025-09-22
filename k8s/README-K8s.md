# BGC — Deploy local em Kubernetes (k3d)

## Pré-requisitos
- Docker, kubectl, k3d, helm
- `hosts` (Windows):  
127.0.0.1 api.bgc.local
127.0.0.1 web.bgc.local

## 1 Cluster k3d
```powershell
k3d cluster create bgc --servers 1 --agents 0 --port "80:80@loadbalancer" --port "443:443@loadbalancer"
kubectl create ns data

## 2 Postgres (Binami Helm)

helm repo add bitnami https://charts.bitnami.com/bitnami
helm install pg bitnami/postgresql -n data `
  --set auth.username=bgc `
  --set auth.password=bgc `
  --set auth.database=bgc
kubectl -n data rollout status statefulset/pg-postgresql

Senha do app bgc: Secret pg-postgresql chave password.

## 3 Build & Import Images

# API
cd .\api
docker build -t bgc-api:v0.2.5 .
k3d image import -c bgc bgc-api:v0.2.5

# Web
cd ..\web
docker build -t bgc-web:v0.1 .
k3d image import -c bgc bgc-web:v0.1

## 4 Manifests

kubectl apply -f .\k8s\api.yaml
kubectl apply -f .\k8s\web.yaml
kubectl -n data get ingress

## 5 Seed (DDL + load + mview)

$raw = kubectl -n data get secret pg-postgresql -o jsonpath="{.data.password}"
$APP_PWD = [Text.Encoding]::UTF8.GetString([Convert]::FromBase64String($raw))

# copiar CSVs
kubectl -n data cp .\stage\exports_ncm_year_sample.csv pg-postgresql-0:/tmp/exports.csv
kubectl -n data cp .\stage\imports_ncm_year_sample.csv pg-postgresql-0:/tmp/imports.csv

# DDL e load
$sql = @'
TRUNCATE trade_ncm_year;
DROP TABLE IF EXISTS stage_raw;
CREATE TEMP TABLE stage_raw (
  co_ncm TEXT, no_ncm_por TEXT, unidade TEXT, ano INT,
  valor_usd NUMERIC, qtd NUMERIC, ncm_chapter TEXT, fluxo TEXT
);
\copy stage_raw FROM '/tmp/exports.csv' WITH (FORMAT csv, HEADER true);
\copy stage_raw FROM '/tmp/imports.csv' WITH (FORMAT csv, HEADER true);
INSERT INTO trade_ncm_year(ano, ncm_chapter, exp_valor_usd, imp_valor_usd)
SELECT ano, LPAD(TRIM(ncm_chapter),2,'0'),
       SUM(CASE WHEN lower(fluxo) LIKE 'export%' THEN COALESCE(valor_usd,0) ELSE 0 END),
       SUM(CASE WHEN lower(fluxo) LIKE 'import%' THEN COALESCE(valor_usd,0) ELSE 0 END)
FROM stage_raw GROUP BY ano, LPAD(TRIM(ncm_chapter),2,'0');
REFRESH MATERIALIZED VIEW m_tam_by_year_chapter;
'@
$sql | kubectl -n data exec -i pg-postgresql-0 -- `
  env PGPASSWORD=$APP_PWD /opt/bitnami/postgresql/bin/psql `
  -U bgc -d bgc -v ON_ERROR_STOP=1 -f -

## 6 Smoke Tests

iwr http://api.bgc.local/healthz
iwr "http://api.bgc.local/market/size?metric=TAM&year_from=2024&year_to=2025&ncm_chapter=84"
iwr "http://api.bgc.local/market/size?metric=SOM&scenario=aggressive&year_from=2024&year_to=2025"
iwr "http://api.bgc.local/routes/compare?from=USA&alts=CHN,ARE,IND&ncm_chapter=84&year=2024&tariff_scenario=tarifa10"

# UI
# http://web.bgc.local/index.html
