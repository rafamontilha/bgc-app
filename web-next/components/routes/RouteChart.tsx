'use client';

import { Bar } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ChartOptions,
} from 'chart.js';
import { formatCurrency } from '@/lib/formatters';
import type { RouteComparisonResult } from '@/types/api';

// Registrar componentes do Chart.js
ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

interface RouteChartProps {
  results: RouteComparisonResult[];
}

export function RouteChart({ results }: RouteChartProps) {
  const data = {
    labels: results.map((r) => r.partner),
    datasets: [
      {
        label: 'Estimado (USD)',
        data: results.map((r) => r.estimated_usd),
        backgroundColor: '#3b82f6',
        borderColor: '#3b82f6',
        borderWidth: 0,
      },
    ],
  };

  const options: ChartOptions<'bar'> = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        callbacks: {
          label: (context) => formatCurrency(context.parsed.y ?? 0),
        },
      },
    },
    scales: {
      x: {
        grid: {
          display: false,
        },
        ticks: {
          color: '#9fb0c3',
        },
      },
      y: {
        ticks: {
          callback: (value) => formatCurrency(Number(value)),
          color: '#9fb0c3',
        },
        grid: {
          color: 'rgba(255, 255, 255, 0.06)',
        },
      },
    },
  };

  return (
    <div className="bg-[#0b1322] border border-outline rounded-[12px] p-2 h-[280px]">
      <Bar data={data} options={options} />
    </div>
  );
}
