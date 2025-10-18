import {NavLink, type NavLinkProps} from "react-router-dom";
import clsx from "clsx";

type LinkVariant = "primary" | "secondary" | "success" | "warning" | "danger" | "info" | "light" | "dark";

interface SideBarLinkProps extends Omit<NavLinkProps, "className"> {
    variant?: LinkVariant;
    activeGlow?: boolean
}

const VARIANT_CLASSES: Record<LinkVariant, string> = {
    primary: "text-indigo-400 hover:text-indigo-300 hover:bg-indigo-900/20",
    secondary: "text-gray-300 hover:text-gray-100 hover:bg-gray-800/50",
    success: "text-green-400 hover:text-green-300 hover:bg-green-900/20",
    danger: "text-red-400 hover:text-red-300 hover:bg-red-900/20",
    warning: "text-yellow-400 hover:text-yellow-300 hover:bg-yellow-900/20",
    info: "text-blue-400 hover:text-blue-300 hover:bg-blue-900/20",
    light: "text-gray-200 hover:text-white hover:bg-gray-700/40",
    dark: "text-gray-400 hover:text-gray-200 hover:bg-gray-800/40",
}

export const SideBarLink = ({
                                variant = "primary",
                                activeGlow = true,
                                ...props
                            }: SideBarLinkProps) => {
    return (
        <NavLink
            {...props}
            className={({ isActive }) =>
                clsx(
                    "block px-4 py-2 rounded-lg font-medium transition-colors duration-200 cursor-pointer",
                    VARIANT_CLASSES[variant],
                    isActive &&
                    (activeGlow
                        ? "bg-gray-800/60 shadow-[0_0_8px_rgba(99,102,241,0.6)] text-white"
                        : "bg-gray-800 text-white")
                )
            }
        />
    );
};