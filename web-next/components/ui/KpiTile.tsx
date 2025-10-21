interface KpiTileProps {
  label: string;
  value: string | number;
  className?: string;
}

export function KpiTile({ label, value, className = '' }: KpiTileProps) {
  return (
    <div
      className={`
        bg-[#0b1322] border border-outline rounded-[12px] p-3
        ${className}
      `}
    >
      <div className="text-xs text-on-surface-variant">{label}</div>
      <div className="text-lg font-extrabold mt-1">{value}</div>
    </div>
  );
}
