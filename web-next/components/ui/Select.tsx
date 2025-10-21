import { SelectHTMLAttributes, ReactNode } from 'react';

interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  children: ReactNode;
}

export function Select({ label, children, className = '', ...props }: SelectProps) {
  return (
    <div className="flex flex-col gap-1.5">
      {label && (
        <label className="text-xs text-on-surface-variant leading-tight">
          {label}
        </label>
      )}
      <select
        className={`
          rounded-[10px] border border-outline
          bg-[#0b1322] text-on-surface
          px-3 py-2.5 text-sm
          disabled:bg-outline-variant disabled:text-[#7e8ea4]
          focus:outline-none focus:ring-2 focus:ring-primary
          ${className}
        `}
        {...props}
      >
        {children}
      </select>
    </div>
  );
}
