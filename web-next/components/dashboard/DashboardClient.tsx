'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Card } from '@/components/ui/Card';
import { Input } from '@/components/ui/Input';
import { Select } from '@/components/ui/Select';
import { Button } from '@/components/ui/Button';
import { KpiTile } from '@/components/ui/KpiTile';
import { ErrorMessage } from '@/components/ui/ErrorMessage';
import { useMarketData } from '@/hooks/use-market-data';
import { aggregateByYear, generateCSV, downloadCSV } from '@/lib/utils';
import { formatCurrency, formatNumber } from '@/lib/formatters';
import type { Metric, Scenario, MarketSizeParams } from '@/types/api';

export function DashboardClient() {
  const [metric, setMetric] = useState<Metric>('TAM');
  const [yearFrom, setYearFrom] = useState(2020);
  const [yearTo, setYearTo] = useState(2025);
  const [ncmChapter, setNcmChapter] = useState('');
  const [scenario, setScenario] = useState<Scenario>('base');
  const [queryParams, setQueryParams] = useState<MarketSizeParams | null>(null);

  const { data, error, isLoading } = useMarketData(queryParams);

  const aggregated = data ? aggregateByYear(data.items) : null;

  const isSOM = metric === 'SOM';

  const handleConsultar = () => {
    const params: MarketSizeParams = {
      metric,
      year_from: yearFrom,
      year_to: yearTo,
    };

    if (ncmChapter.trim()) {
      params.ncm_chapter = ncmChapter.trim().padStart(2, '0');
    }

    if (metric === 'SOM') {
      params.scenario = scenario;
    }

    setQueryParams(params);
  };

  const handleExport = () => {
    if (!aggregated) return;

    const rows = aggregated.rows.map((r) => [r.ano, metric, r.total]);
    const csv = generateCSV(['ano', 'metric', 'valor_usd'], rows);
    downloadCSV(`bgc_${metric.toLowerCase()}_${Date.now()}.csv`, csv);
  };

  const handleClear = () => {
    setMetric('TAM');
    setYearFrom(2020);
    setYearTo(2025);
    setNcmChapter('');
    setScenario('base');
    setQueryParams(null);
  };

  return (
    <div className="max-w-[1080px] mx-auto px-6 py-6">
      <header className="mb-6">
        <h1 className="text-[22px] font-bold mb-2">
          📊 BGC — Dashboard TAM / SAM / SOM
        </h1>
        <p className="text-on-surface-variant text-sm">
          Consome <code className="bg-surface-variant px-1 rounded">/market/size</code> da API
          local (porta 8080). Ajuste os filtros e exporte CSV.
        </p>
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-[320px_1fr] gap-4">
        {/* Card de filtros */}
        <Card>
          <Select
            label="Métrica"
            value={metric}
            onChange={(e) => setMetric(e.target.value as Metric)}
          >
            <option value="TAM">TAM — Mercado Total</option>
            <option value="SAM">SAM — Mercado Atendível</option>
            <option value="SOM">SOM — Mercado Obtenível</option>
          </Select>

          <div className="grid grid-cols-2 gap-3 mt-2.5">
            <Input
              label="Ano de"
              type="number"
              min={2020}
              max={2025}
              value={yearFrom}
              onChange={(e) => setYearFrom(Number(e.target.value))}
            />
            <Input
              label="Ano até"
              type="number"
              min={2020}
              max={2025}
              value={yearTo}
              onChange={(e) => setYearTo(Number(e.target.value))}
            />
          </div>

          <div className="mt-2.5">
            <Select
              label={isSOM ? 'Cenário (SOM)' : 'Cenário (somente SOM)'}
              value={scenario}
              onChange={(e) => setScenario(e.target.value as Scenario)}
              disabled={!isSOM}
            >
              <option value="base">base — 1,5%</option>
              <option value="aggressive">aggressive — 3%</option>
            </Select>
            <div
              className="text-[11px] text-on-surface-variant mt-1.5"
              style={{ opacity: isSOM ? 1 : 0.6 }}
            >
              O cenário só altera valores quando a métrica = SOM.
            </div>
          </div>

          <Input
            label="Capítulo NCM (2 dígitos)"
            placeholder="Ex.: 84"
            pattern="^[0-9]{2}$"
            value={ncmChapter}
            onChange={(e) => setNcmChapter(e.target.value)}
            className="mt-2.5"
          />

          <div className="flex gap-2 flex-wrap mt-2.5">
            <Button onClick={handleConsultar} disabled={isLoading}>
              {isLoading ? 'Carregando...' : 'Consultar'}
            </Button>
            <Button
              variant="secondary"
              onClick={handleExport}
              disabled={!aggregated || aggregated.rows.length === 0}
            >
              Exportar CSV
            </Button>
            <Button variant="secondary" onClick={handleClear}>
              Limpar
            </Button>
          </div>

          {error && <ErrorMessage message={error.message} />}

          {data && (
            <div className="mt-2.5 text-on-surface-variant text-xs break-all">
              Métrica: {data.metric}
              {data.metric === 'SOM' && ` (${data.scenario || 'base'})`}
            </div>
          )}
        </Card>

        {/* Card de resultados */}
        <Card>
          <div className="grid grid-cols-3 gap-3">
            <KpiTile
              label="Total de linhas (API)"
              value={data?.items?.length ? formatNumber(data.items.length) : '—'}
            />
            <KpiTile
              label="Capítulos únicos"
              value={aggregated ? formatNumber(aggregated.chapters) : '—'}
            />
            <KpiTile
              label="Soma (USD, período)"
              value={aggregated ? formatCurrency(aggregated.sum) : '—'}
            />
          </div>

          <div className="mt-2 overflow-x-auto">
            <table className="w-full border-collapse">
              <thead>
                <tr>
                  <th className="text-left py-2.5 px-2 border-b border-outline text-xs text-on-surface-variant font-semibold uppercase tracking-wide">
                    Ano
                  </th>
                  <th className="text-right py-2.5 px-2 border-b border-outline text-xs text-on-surface-variant font-semibold uppercase tracking-wide">
                    Total USD (agregado)
                  </th>
                </tr>
              </thead>
              <tbody>
                {aggregated?.rows.map((row) => (
                  <tr key={row.ano} className="hover:bg-white/[0.03]">
                    <td className="text-left py-2.5 px-2 border-b border-outline">
                      {row.ano}
                    </td>
                    <td className="text-right py-2.5 px-2 border-b border-outline">
                      {formatCurrency(row.total)}
                    </td>
                  </tr>
                ))}
                {!aggregated && (
                  <tr>
                    <td colSpan={2} className="text-center py-6 text-on-surface-variant">
                      Clique em &ldquo;Consultar&rdquo; para ver os resultados
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </Card>
      </div>

      <div className="mt-6 text-center">
        <Link
          href="/routes"
          className="inline-block bg-[#0b1322] border border-outline rounded-[10px] px-3 py-2.5 text-sm hover:bg-surface-variant transition-colors"
        >
          🧭 Ver Comparação de Rotas →
        </Link>
      </div>
    </div>
  );
}
