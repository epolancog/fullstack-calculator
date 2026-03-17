interface ErrorMessageProps {
  message: string | null;
}

export function ErrorMessage({ message }: ErrorMessageProps) {
  if (!message) return null;

  return (
    <div
      role="alert"
      className="flex items-center gap-2 rounded-lg bg-red-500/15 border border-red-400/30 px-3 py-2 text-sm text-red-300"
    >
      <span aria-hidden="true">⚠</span>
      <span>{message}</span>
    </div>
  );
}
