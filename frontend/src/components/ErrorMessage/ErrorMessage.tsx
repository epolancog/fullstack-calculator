interface ErrorMessageProps {
  message: string | null;
}

export function ErrorMessage({ message }: ErrorMessageProps) {
  if (!message) return null;

  return (
    <div
      role="alert"
      className="animate-shake flex items-center gap-2 rounded-xl bg-red-500/10 border border-red-400/20 backdrop-blur-sm px-3 py-2 text-sm text-red-300"
    >
      <span aria-hidden="true">⚠</span>
      <span>{message}</span>
    </div>
  );
}
