import { InputHTMLAttributes } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
}

export function Input({ label, className = '', ...props }: InputProps) {
  return (
    <div className="flex flex-col gap-1.5">
      {label && (
        <label className="text-xs text-on-surface-variant leading-tight">
          {label}
        </label>
      )}
      <input
        className={`
          rounded-[10px] border border-outline
          bg-[#0b1322] text-on-surface
          px-3 py-2.5 text-sm
          placeholder:text-[#5a6b80]
          disabled:bg-outline-variant disabled:text-[#7e8ea4]
          focus:outline-none focus:ring-2 focus:ring-primary
          ${className}
        `}
        {...props}
      />
    </div>
  );
}
