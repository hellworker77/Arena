import {NavLink, type NavLinkProps} from "react-router-dom";
import clsx from "clsx";

type LinkVariant = "primary" | "secondary" | "success" | "warning" | "danger" | "info" | "light" | "dark";

interface NavBarLinkProps extends Omit<NavLinkProps, "className"> {
    variant?: LinkVariant;
    activeGlow?: boolean
}

const VARIANT_CLASSES: Record<LinkVariant, string> = {
    primary: "text-indigo-400 hover:text-indigo-300",
    secondary: "text-gray-300 hover:text-gray-100",
    success: "text-green-400 hover:text-green-300",
    danger: "text-red-400 hover:text-red-300",
    warning: "text-yellow-400 hover:text-yellow-300",
    info: "text-blue-400 hover:text-blue-300",
    light: "text-gray-200 hover:text-white",
    dark: "text-gray-400 hover:text-gray-200",
}

export const NavBarLink = ({
                                variant = "primary",
                                activeGlow = false,
                                ...props
                            }: NavBarLinkProps) => {
    return (
        <NavLink
            {...props}
            className={({ isActive }) =>
                clsx(
                    "block px-4 py-2 rounded-lg font-medium transition-colors duration-200 cursor-pointer",
                    VARIANT_CLASSES[variant],
                    isActive &&
                    (activeGlow
                        ? "shadow-[0_0_8px_rgba(99,102,241,0.6)] text-white"
                        : "text-white")
                )
            }
        />
    );
};