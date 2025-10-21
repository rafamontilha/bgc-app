/**
 * Utility functions
 */

import type { MarketSizeItem, AggregatedYearData } from '@/types/api';

/**
 * Agrega dados por ano
 */
export function aggregateByYear(items: MarketSizeItem[]): {
  rows: AggregatedYearData[];
  chapters: number;
  sum: number;
} {
  const map = new Map<number, number>();
  const chaptersSet = new Set<string>();
  let sum = 0;

  for (const item of items) {
    chaptersSet.add(item.ncm_chapter);
    sum += item.valor_usd;
    const ano = item.ano;
    map.set(ano, (map.get(ano) || 0) + item.valor_usd);
  }

  const rows = Array.from(map.entries())
    .sort((a, b) => a[0] - b[0])
    .map(([ano, total]) => ({ ano, total }));

  return {
    rows,
    chapters: chaptersSet.size,
    sum,
  };
}

/**
 * Gera CSV a partir de dados
 */
export function generateCSV(
  headers: string[],
  rows: (string | number)[][]
): string {
  const csvRows = [headers, ...rows];
  return csvRows.map((row) => row.join(',')).join('\n');
}

/**
 * Download de arquivo CSV
 */
export function downloadCSV(filename: string, content: string): void {
  const blob = new Blob([content], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();
  URL.revokeObjectURL(url);
}
