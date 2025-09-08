Param(
  [string]$StageDir    = "$PSScriptRoot\..\stage",
  [string]$ComposeFile = "$PSScriptRoot\..\bgcstack\docker-compose.yml"
)

function Invoke-Psql {
  param([string]$Sql)
  & docker compose -f $ComposeFile exec -T db psql -U bgc -d bgc -c $Sql
  if ($LASTEXITCODE -ne 0) { throw "psql falhou: $Sql" }
}

try {
  # PrÃ©-checks
  $need = @("lookup_ncm8.csv","exports_ncm_year_sample.csv","imports_ncm_year_sample.csv")
  foreach ($f in $need) {
    $p = Join-Path $StageDir $f
    if (-not (Test-Path $p)) { throw "Arquivo nÃ£o encontrado: $p" }
  }
  if (-not (Test-Path $ComposeFile)) { throw "Compose nÃ£o encontrado: $ComposeFile" }

  # Sobe/verifica stack
  & docker compose -f $ComposeFile up -d
  if ($LASTEXITCODE -ne 0) { throw "Falha ao subir compose." }

  # Espera Postgres saudÃ¡vel
  $ok = $false
  for ($i=0; $i -lt 60; $i++) {
    $state = (docker inspect -f '{{.State.Health.Status}}' bgc_db 2>$null)
    if ($state -eq "healthy") { $ok = $true; break }
    Start-Sleep -Seconds 2
  }
  if (-not $ok) { throw "Postgres nÃ£o ficou healthy a tempo." }

  # Aplica schema idempotente que jÃ¡ estÃ¡ montado em /docker-entrypoint-initdb.d
  & docker compose -f $ComposeFile exec -T db psql -U bgc -d bgc -f /docker-entrypoint-initdb.d/00_schema.sql
  if ($LASTEXITCODE -ne 0) { throw "Falha aplicando schema." }

  # Copia CSVs para o container
  & docker cp (Join-Path $StageDir "lookup_ncm8.csv")            bgc_db:/tmp/lookup_ncm8.csv
  & docker cp (Join-Path $StageDir "exports_ncm_year_sample.csv") bgc_db:/tmp/exports_ncm_year_sample.csv
  & docker cp (Join-Path $StageDir "imports_ncm_year_sample.csv") bgc_db:/tmp/imports_ncm_year_sample.csv

  # SQLs usando dollar-quoted strings (evita aspas simples)
  $truncateLookup = 'TRUNCATE ncm_lookup;'
  $truncateTrade  = 'TRUNCATE trade_ncm_year;'

  # o $csv$ NÃO é expandido e chega literal no psql.
  $copyLookup = 'COPY ncm_lookup (co_ncm,no_ncm_por,co_sh2,no_sh2_por) FROM $csv$/tmp/lookup_ncm8.csv$csv$ CSV HEADER;'
  $copyExp    = 'COPY trade_ncm_year (ncm,ncm_desc,unidade,ano,valor_usd_fob,quantidade_estat,ncm_chapter,fluxo) FROM $csv$/tmp/exports_ncm_year_sample.csv$csv$ CSV HEADER;'
  $copyImp    = 'COPY trade_ncm_year (ncm,ncm_desc,unidade,ano,valor_usd_fob,quantidade_estat,ncm_chapter,fluxo) FROM $csv$/tmp/imports_ncm_year_sample.csv$csv$ CSV HEADER;'

  $refreshMV  = 'REFRESH MATERIALIZED VIEW v_tam_by_year_chapter;'

  Invoke-Psql $truncateLookup
  Invoke-Psql $truncateTrade
  Invoke-Psql $copyLookup
  Invoke-Psql $copyExp
  Invoke-Psql $copyImp
  Invoke-Psql $refreshMV

  # EvidÃªncias
  Invoke-Psql "SELECT COUNT(*) AS linhas_trade FROM trade_ncm_year;"
  Invoke-Psql "SELECT * FROM v_tam_by_year_chapter ORDER BY ano, ncm_chapter LIMIT 20;"
  Invoke-Psql "SELECT * FROM v_quality_orphans_ncm LIMIT 20;"
}
catch {
  Write-Error $_
  exit 1
}
