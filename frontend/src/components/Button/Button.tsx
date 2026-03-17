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
    "focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white/50",
    "disabled:opacity-40 disabled:cursor-not-allowed disabled:active:scale-100",
    "min-h-[48px] text-lg",
  ],
  {
    variants: {
      variant: {
        number:
          "bg-white/10 border-white/20 text-white hover:bg-white/20",
        operator:
          "bg-indigo-500/30 border-indigo-400/30 text-indigo-200 hover:bg-indigo-500/40",
        action:
          "bg-white/5 border-white/10 text-gray-300 hover:bg-white/15",
        equals:
          "bg-indigo-600 border-indigo-500 text-white hover:bg-indigo-500",
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
