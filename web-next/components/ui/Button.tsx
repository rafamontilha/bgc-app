import { ButtonHTMLAttributes, ReactNode } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary';
  children: ReactNode;
}

export function Button({
  variant = 'primary',
  children,
  className = '',
  disabled,
  ...props
}: ButtonProps) {
  const baseStyles = `
    rounded-[10px] px-3 py-2.5 text-sm font-bold
    cursor-pointer transition-opacity
    disabled:opacity-60 disabled:cursor-not-allowed
  `;

  const variantStyles = {
    primary: 'bg-gradient-to-r from-primary to-secondary border-none text-on-primary',
    secondary: 'bg-[#0b1322] border border-outline text-on-surface',
  };

  return (
    <button
      className={`${baseStyles} ${variantStyles[variant]} ${className}`}
      disabled={disabled}
      {...props}
    >
      {children}
    </button>
  );
}
