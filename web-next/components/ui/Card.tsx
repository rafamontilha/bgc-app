import { ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  className?: string;
}

export function Card({ children, className = '' }: CardProps) {
  return (
    <div
      className={`
        bg-gradient-to-b from-surface to-surface-variant
        border border-outline rounded-[16px] p-4
        shadow-[0_10px_30px_rgba(0,0,0,0.25)]
        ${className}
      `}
    >
      {children}
    </div>
  );
}
