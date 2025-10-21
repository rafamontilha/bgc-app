interface ErrorMessageProps {
  message: string;
}

export function ErrorMessage({ message }: ErrorMessageProps) {
  return (
    <div className="text-error font-semibold text-sm mt-2">
      {message}
    </div>
  );
}
