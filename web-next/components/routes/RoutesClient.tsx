'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { Card } from '@/components/ui/Card';
import { Input } from '@/components/ui/Input';
import { Select } from '@/components/ui/Select';
import { Button } from '@/components/ui/Button';
import { KpiTile } from '@/components/ui/KpiTile';
import { ErrorMessage } from '@/components/ui/ErrorMessage';
import { RouteChart } from '@/components/routes/RouteChart';
import { useRoutesData } from '@/hooks/use-routes-data';
import { useScenarios } from '@/hooks/use-scenarios';
import { generateCSV, downloadCSV } from '@/lib/utils';
import { formatCurrency, formatPercent, formatNumber } from '@/lib/formatters';
import type { RouteComparisonParams } from '@/types/api';

export function RoutesClient() {
  const [year, setYear] = useState(2024);
  const [chapter, setChapter] = useState('84');
  const [from, setFrom] = useState('USA');
  const [alts, setAlts] = useState('CHN,ARE,IND');
  const [tariffScenario, setTariffScenario] = useState('base');
  const [queryParams, setQueryParams] = useState<RouteComparisonParams | null>(null);

  const { scenarios } = useScenarios();
  const { data, error, isLoading } = useRoutesData(queryParams);

  // Auto-load on mount
  useEffect(() => {
    handleComparar();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleComparar = () => {
    const params: RouteComparisonParams = {
      year,
      ncm_chapter: chapter.trim().padStart(2, '0'),
      from: from.trim().toUpperCase(),
      alts: alts.trim().toUpperCase(),
      tariff_scenario: tariffScenario,
    };

    setQueryParams(params);
  };

  const handleExport = () => {
    if (!data?.results) return;

    const rows = data.results.map((r) => [
      r.partner,
      r.share,
      r.factor ?? 1,
      r.estimated_usd,
    ]);
    const csv = generateCSV(['partner', 'share', 'factor', 'estimated_usd'], rows);
    downloadCSV(
      `bgc_routes_${data.ncm_chapter}_${data.year}_${data.tariff_scenario}.csv`,
      csv
    );
  };

  const sum = data?.results?.reduce((acc, r) => acc + r.estimated_usd, 0) || 0;
  const isCheckOk = Math.abs(sum - (data?.adjusted_total_usd || 0)) < 0.5;

  return (
    <div className="max-w-[1180px] mx-auto px-6 py-6">
      <header className="mb-6">
        <h1 className="text-[22px] font-bold mb-2">🧭 BGC — EUA vs Alternativos</h1>
        <p className="text-on-surface-variant text-sm">
          Compara um parceiro <strong>from</strong> (ex.:{' '}
          <code className="bg-surface-variant px-1 rounded">USA</code>) com alternativos (ex.:{' '}
          <code className="bg-surface-variant px-1 rounded">CHN,ARE,IND</code>) para um{' '}
          <strong>capítulo NCM</strong> e <strong>ano</strong>. Consome{' '}
          <code className="bg-surface-variant px-1 rounded">/routes/compare</code>.
        </p>
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-[340px_1fr] gap-4">
        {/* Card de filtros */}
        <Card>
          <div className="grid grid-cols-2 gap-3">
            <Input
              label="Ano"
              type="number"
              min={2020}
              max={2025}
              value={year}
              onChange={(e) => setYear(Number(e.target.value))}
            />
            <Input
              label="Capítulo NCM (2 dígitos)"
              placeholder="Ex.: 84"
              pattern="^[0-9]{2}$"
              value={chapter}
              onChange={(e) => setChapter(e.target.value)}
            />
          </div>

          <Input
            label="Parceiro principal"
            value={from}
            onChange={(e) => setFrom(e.target.value)}
            className="mt-2.5"
          />

          <Input
            label="Alternativos (separe por vírgulas)"
            value={alts}
            onChange={(e) => setAlts(e.target.value)}
            className="mt-2.5"
          />

          <div className="mt-2.5">
            <Select
              label="Cenário de Tarifa"
              value={tariffScenario}
              onChange={(e) => setTariffScenario(e.target.value)}
            >
              {scenarios.map((s) => (
                <option key={s} value={s}>
                  {s}
                </option>
              ))}
            </Select>
            <div className="text-[11px] text-on-surface-variant mt-1.5">
              Fator multiplicativo por parceiro/capítulo/ano (ver /healthz →
              available_scenarios).
            </div>
          </div>

          <div className="flex gap-2 flex-wrap mt-2.5">
            <Button onClick={handleComparar} disabled={isLoading}>
              {isLoading ? 'Carregando...' : 'Comparar'}
            </Button>
            <Button
              variant="secondary"
              onClick={handleExport}
              disabled={!data?.results || data.results.length === 0}
            >
              Exportar CSV
            </Button>
            <Link
              href="/"
              className="bg-[#0b1322] border border-outline rounded-[10px] px-3 py-2.5 text-sm inline-block hover:bg-surface-variant transition-colors"
            >
              ← Voltar Dashboard
            </Link>
          </div>

          {error && <ErrorMessage message={error.message} />}

          {data && (
            <div className="mt-2.5 text-on-surface-variant text-xs break-all">
              ano={data.year}, capítulo={data.ncm_chapter}, cenário={data.tariff_scenario}
            </div>
          )}
        </Card>

        {/* Card de resultados */}
        <Card>
          <div className="grid grid-cols-4 gap-3">
            <KpiTile
              label="TAM (USD)"
              value={data ? formatCurrency(data.tam_total_usd) : '—'}
            />
            <KpiTile
              label="Ajustado (USD)"
              value={data ? formatCurrency(data.adjusted_total_usd) : '—'}
            />
            <KpiTile
              label="Parceiros"
              value={data?.results?.length ? formatNumber(data.results.length) : '—'}
            />
            <KpiTile
              label="Checagem de soma"
              value={data ? (isCheckOk ? 'OK ✓' : 'ALERTA ⚠') : '—'}
              className={data ? (isCheckOk ? 'text-success' : 'text-error') : ''}
            />
          </div>

          {data?.results && <RouteChart results={data.results} />}

          <div className="mt-2 overflow-x-auto">
            <table className="w-full border-collapse">
              <thead>
                <tr>
                  <th className="text-left py-2.5 px-2 border-b border-outline text-xs text-on-surface-variant font-semibold uppercase tracking-wide">
                    Parceiro
                  </th>
                  <th className="text-right py-2.5 px-2 border-b border-outline text-xs text-on-surface-variant font-semibold uppercase tracking-wide">
                    Share
                  </th>
                  <th className="text-right py-2.5 px-2 border-b border-outline text-xs text-on-surface-variant font-semibold uppercase tracking-wide">
                    Factor
                  </th>
                  <th className="text-right py-2.5 px-2 border-b border-outline text-xs text-on-surface-variant font-semibold uppercase tracking-wide">
                    Estimado (USD)
                  </th>
                </tr>
              </thead>
              <tbody>
                {data?.results?.map((row) => (
                  <tr key={row.partner} className="hover:bg-white/[0.03]">
                    <td className="text-left py-2.5 px-2 border-b border-outline">
                      {row.partner}
                    </td>
                    <td className="text-right py-2.5 px-2 border-b border-outline">
                      {formatPercent(row.share)}
                    </td>
                    <td className="text-right py-2.5 px-2 border-b border-outline">
                      {(row.factor ?? 1).toFixed(2)}
                    </td>
                    <td className="text-right py-2.5 px-2 border-b border-outline">
                      {formatCurrency(row.estimated_usd)}
                    </td>
                  </tr>
                ))}
                {!data && (
                  <tr>
                    <td colSpan={4} className="text-center py-6 text-on-surface-variant">
                      Clique em &ldquo;Comparar&rdquo; para ver os resultados
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </Card>
      </div>
    </div>
  );
}
