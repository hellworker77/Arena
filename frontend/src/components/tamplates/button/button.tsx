import {type ButtonHTMLAttributes, forwardRef, type ReactNode} from "react";
import clsx from "clsx";

type ButtonVariant = "primary" | "warning" | "danger" | "secondary" | "success" | "info" | "light" | "dark" | "link";
type ButtonSize = "sm" | "md" | "lg";

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {

    variant?: ButtonVariant;

    size?: ButtonSize;

    leftIcon?: ReactNode;

    rightIcon?: ReactNode;

    loading?: boolean;
}

const VARIANT_CLASSES: Record<ButtonVariant, string> = {
    primary: "bg-indigo-600 text-white hover:bg-indigo-700 focus:ring-indigo-500",
    secondary: "bg-magenta-200 text-gray-800 hover:bg-gray-300 focus:ring-gray-400",
    success: "bg-green-600 text-white hover:bg-green-700 focus:ring-green-500",
    danger: "bg-red-600 text-white hover:bg-red-700 focus:ring-red-500",
    warning: "bg-yellow-400 text-gray-900 hover:bg-yellow-500 focus:ring-yellow-400",
    info: "bg-blue-400 text-white hover:bg-blue-500 focus:ring-blue-400",
    light: "bg-gray-100 text-gray-800 hover:bg-gray-200 focus:ring-gray-200",
    dark: "bg-gray-900 text-white hover:bg-gray-800 focus:ring-gray-900",
    link: "bg-transparent text-indigo-600 hover:underline focus:ring-0",
}

const SIZE_CLASSES: Record<ButtonSize, string> = {
    sm: "px-3 py-1.5 text-sm rounded-md",
    md: "px-4 py-2 text-base rounded-lg",
    lg: "px-6 py-3 text-lg rounded-xl",
};

export const Button = forwardRef<HTMLButtonElement, ButtonProps>((props, ref) => {

    const {
        variant = "primary",
        size = "md",
        leftIcon,
        rightIcon,
        loading = false,
        disabled,
        children,
        className = "",
        ...rest
    } = props;


    return (
        <button ref={ref}
                disabled={disabled || loading}
                className={clsx(VARIANT_CLASSES[variant],
                    SIZE_CLASSES[size],
                    className,
                    "inline-flex items-center justify-center gap-2 " +
                    "font-semibold transition focus:outline-none focus:ring-2 " +
                    "focus:ring-offset-1 disabled:opacity-50 disabled:cursor-not-allowed")}
                {...rest}>

            {loading && (
                <svg
                    className="animate-spin h-5 w-5 text-current"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                >
                    <circle
                        className="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        strokeWidth="4"
                    />
                    <path
                        className="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
                    />
                </svg>
            )}

            {leftIcon && !loading && <span>{leftIcon}</span>}

            {children}

            {rightIcon && !loading && <span>{rightIcon}</span>}
        </button>
    )
})