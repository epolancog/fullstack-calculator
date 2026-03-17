interface DisplayProps {
  expression: string;
  currentInput: string;
  result: string | null;
  isLoading: boolean;
}

function formatNumber(value: string): string {
  // Don't format if it contains operators or special chars (expression string)
  if (/[+\-*/^%() a-zA-Z]/.test(value)) return value;

  const num = Number(value);
  if (isNaN(num)) return value;

  // Use toPrecision for large numbers, toLocaleString for comma formatting
  if (Math.abs(num) >= 1e12 || (Math.abs(num) < 1e-10 && num !== 0)) {
    return num.toExponential(6);
  }

  // Cap at 12 significant digits
  const formatted = parseFloat(num.toPrecision(12));
  return formatted.toLocaleString("en-US", { maximumFractionDigits: 10 });
}

export function Display({
  expression,
  currentInput,
  result,
  isLoading,
}: DisplayProps) {
  const mainDisplay = result !== null ? formatNumber(result) : currentInput || "0";

  return (
    <div className="w-full rounded-xl bg-black/40 border border-white/5 p-4 shadow-inner">
      <div
        className="text-right text-sm text-gray-400 min-h-5 overflow-hidden text-ellipsis whitespace-nowrap tracking-wide"
        aria-label="expression"
      >
        {expression || "\u00A0"}
      </div>
      <div
        className={`text-right text-3xl font-light text-white min-h-10 overflow-hidden text-ellipsis whitespace-nowrap tracking-tight ${result === null ? "" : "animate-fade-in"}`}
        aria-live="polite"
        aria-label="display value"
      >
        {isLoading ? (
          <span className="animate-pulse text-gray-500">...</span>
        ) : (
          mainDisplay
        )}
      </div>
    </div>
  );
}
