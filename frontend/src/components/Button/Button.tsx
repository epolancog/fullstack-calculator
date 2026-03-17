import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "../../lib/utils";
import type { ButtonSize, ButtonVariant } from "../../types";

const buttonVariants = cva(
  [
    "flex items-center justify-center",
    "rounded-xl border font-medium",
    "cursor-pointer select-none",
    "transition-all duration-150 ease-in-out",
    "active:scale-95",
    "focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white/40",
    "disabled:opacity-40 disabled:cursor-not-allowed disabled:active:scale-100",
    "min-h-12 text-lg",
    "backdrop-blur-sm",
  ],
  {
    variants: {
      variant: {
        number:
          "bg-glass-button border-glass-border text-white hover:bg-glass-button-hover",
        operator:
          "bg-accent-muted border-accent-border text-indigo-100 hover:bg-accent/40",
        action:
          "bg-white/5 border-white/8 text-gray-300 hover:bg-white/12",
        equals:
          "bg-accent border-accent-hover text-white font-semibold hover:bg-accent-hover shadow-lg shadow-accent/20",
      } satisfies Record<ButtonVariant, string>,
      size: {
        default: "col-span-1",
        wide: "col-span-2",
      } satisfies Record<ButtonSize, string>,
    },
    defaultVariants: {
      variant: "number",
      size: "default",
    },
  }
);

interface ButtonProps extends VariantProps<typeof buttonVariants> {
  label: string;
  onClick: () => void;
  disabled?: boolean;
  ariaLabel?: string;
}

export function Button({
  label,
  onClick,
  variant,
  size,
  disabled = false,
  ariaLabel,
}: ButtonProps) {
  return (
    <button
      className={cn(buttonVariants({ variant, size }))}
      onClick={onClick}
      disabled={disabled}
      aria-label={ariaLabel ?? label}
    >
      {label}
    </button>
  );
}
