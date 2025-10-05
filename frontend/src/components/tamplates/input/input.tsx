import React, {type ChangeEvent, forwardRef, type InputHTMLAttributes, useId, useState} from "react";
import clsx from "clsx";
import {Eye, EyeOff, X} from "lucide-react";

type InputSize = "sm" | "md" | "lg";

type InputVariant = "outline" | "solid" | "ghost";

const SIZE_CLASSES: Record<InputSize, string> = {
    sm: "px-2 py-1 text-sm rounded-md",
    md: "px-3 py-2 text-base rounded-lg",
    lg: "px-4 py-3 text-lg rounded-xl",
};

const VARIANT_CLASSES: Record<InputVariant, string> = {
    outline:
        "bg-white border shadow-sm focus:ring-2 focus:ring-offset-1 focus:ring-indigo-500 border-gray-300",
    solid: "bg-gray-100 border-0 focus:ring-2 focus:ring-indigo-500",
    ghost: "bg-transparent border-0 focus:ring-0",
};

export interface InputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, "size" | "prefix" | "suffix" | "onChange"> {

    label?: string;

    description?: string;

    error?: string;

    size?: InputSize;

    variant?: InputVariant;

    prefix?: React.ReactNode;

    suffix?: React.ReactNode;

    clearable?: boolean;

    showPasswordToggle?: boolean;

    onChange: (e: ChangeEvent<HTMLInputElement>) => void;
}

export const Input = forwardRef<HTMLInputElement, InputProps>((props, ref) => {

    const {
        id,
        label,
        description,
        error,
        className = "",
        size = "md",
        variant = "outline",
        prefix,
        suffix,
        clearable = false,
        type = "text",
        showPasswordToggle = false,
        value,
        onChange,
        name,
        disabled = false,
        placeholder,
        ...rest
    } = props;

    const autoId = useId();

    const inputId = id || `custom-input-${autoId}`;

    const descrId = description ? `${inputId}-desc` : undefined;

    const errId = error ? `${inputId}-err` : undefined;

    const controlled = value !== undefined;

    const [internalValue, setInternalValue] = useState<string>(value?.toString() ?? "");

    const computedValue = controlled ? value : internalValue;

    const [isPasswordVisible, setIsPasswordVisible] = useState(false);

    const inputType =
        showPasswordToggle && type === "password"
            ? isPasswordVisible
                ? "text"
                : "password"
            : type;

    const baseContainer =
        "w-full flex items-center gap-2 border transition-shadow duration-150 relative";

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
        if (!controlled) setInternalValue(e.target.value);

        onChange?.(e);
    }

    const clear = () => {
        if (!controlled) setInternalValue("")

        onChange?.({target: {value: ""}} as ChangeEvent<HTMLInputElement>)
    }

    return (
        <div className={clsx("w-full", className)}>
            {label &&
                <label htmlFor={inputId}
                       className="block text-sm font-medium text-gray-700 mb-1">
                    {label}{rest.required ? " *" : ""}
                </label>
            }

            <div className={clsx(
                error ? "border-red-400 shadow-red-50" : "",
                SIZE_CLASSES[size],
                VARIANT_CLASSES[variant],
                baseContainer)}>

                {prefix && <div className="flex items-center pl-1">{prefix}</div>}

                <input ref={ref}
                       id={inputId}
                       name={name}
                       type={inputType}
                       value={computedValue}
                       onChange={handleChange}
                       disabled={disabled}
                       placeholder={placeholder}
                       aria-describedby={`${descrId ?? ""} ${errId ?? ""}`.trim() || undefined}
                       aria-invalid={!!error}
                       className="flex-1 bg-transparent outline-none placeholder-gray-400 disabled:opacity-60"
                       {...rest}/>

                <div className="flex items-center gap-1 pr-1">
                    {clearable && computedValue && (
                        <button
                            type="button"
                            onClick={clear}
                            aria-label="Clear input"
                            className="p-1 rounded-md hover:bg-gray-100 focus:outline-none focus:ring-1 focus:ring-indigo-400"
                        >
                            <X size={16} />
                        </button>
                    )}

                    {showPasswordToggle && type === "password" && (
                        <button
                            type="button"
                            onClick={() => setIsPasswordVisible((s) => !s)}
                            aria-label={isPasswordVisible ? "Hide password" : "Show password"}
                            className="p-1 rounded-md hover:bg-gray-100 focus:outline-none focus:ring-1 focus:ring-indigo-400"
                        >
                            {isPasswordVisible ? <EyeOff size={16} /> : <Eye size={16} />}
                        </button>
                    )}

                    {suffix && <div className="flex items-center">{suffix}</div>}
                </div>

                {description && !error && (
                    <p id={descrId} className="mt-1 text-xs text-gray-500">
                        {description}
                    </p>
                )}

                {error && (
                    <p id={errId} className="mt-1 text-xs text-red-600">
                        {error}
                    </p>
                )}
            </div>
        </div>
    )
})