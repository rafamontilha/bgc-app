
#!/usr/bin/env bash
set -euo pipefail

STAGE_DIR="${STAGE_DIR:-./stage}"

EXP="${STAGE_DIR}/exports_ncm_year_sample.csv"
IMP="${STAGE_DIR}/imports_ncm_year_sample.csv"
NCM="${STAGE_DIR}/lookup_ncm8.csv"

for f in "$EXP" "$IMP" "$NCM"; do
  if [ ! -f "$f" ]; then
    echo "Arquivo não encontrado: $f"
    exit 1
  fi
done

echo "⏳ Aguardando DB ficar saudável..."
docker compose ps db >/dev/null 2>&1 || (echo "Suba com: docker compose up -d"; exit 1)
for i in {1..60}; do
  state=$(docker inspect -f '{{.State.Health.Status}}' bgc_db 2>/dev/null || echo "starting")
  if [ "$state" = "healthy" ]; then break; fi
  sleep 2
done

echo "🧹 Limpando dados anteriores..."
docker compose exec -T db psql -U bgc -d bgc -c "TRUNCATE trade_ncm_year;"

echo "📚 Carregando lookup NCM..."
docker compose exec -T db psql -U bgc -d bgc -c "TRUNCATE ncm_lookup;"
cat "$NCM" | docker compose exec -T db psql -U bgc -d bgc -c \
  "COPY ncm_lookup (co_ncm,no_ncm_por,co_sh2,no_sh2_por) FROM STDIN WITH CSV HEADER;"

echo "⬆️  Carregando EXPORTAÇÃO por NCM/ano..."
cat "$EXP" | docker compose exec -T db psql -U bgc -d bgc -c \
  "COPY trade_ncm_year (ncm,ncm_desc,unidade,ano,valor_usd_fob,quantidade_estat,ncm_chapter,fluxo) FROM STDIN WITH CSV HEADER;"

echo "⬆️  Carregando IMPORTAÇÃO por NCM/ano..."
cat "$IMP" | docker compose exec -T db psql -U bgc -d bgc -c \
  "COPY trade_ncm_year (ncm,ncm_desc,unidade,ano,valor_usd_fob,quantidade_estat,ncm_chapter,fluxo) FROM STDIN WITH CSV HEADER;"

echo "🔄 Atualizando materialized view..."
docker compose exec -T db psql -U bgc -d bgc -c "REFRESH MATERIALIZED VIEW v_tam_by_year_chapter;"

echo "✅ Seed finalizado! Dicas:"
echo " - SELECT * FROM v_tam_by_year_chapter ORDER BY ano, ncm_chapter LIMIT 20;"
echo " - SELECT * FROM v_quality_orphans_ncm LIMIT 20;"

