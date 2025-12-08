import {createContext, CSSProperties, ReactNode, RefObject, useContext, useEffect, useRef, useState} from "react";
import {createPortal} from "react-dom";

interface TooltipProps {
    children: ReactNode
}

type TooltipPosition = "top" | "bottom" | "left" | "right";

interface TooltipContextValue {
    show: boolean;
    setShow: (v: boolean) => void;
    triggerRef: RefObject<HTMLDivElement | null>
    position: TooltipPosition
}

const TooltipContext = createContext<TooltipContextValue | null>(null);

export const Tooltip = ({children}: TooltipProps) => {
    const [show, setShow] = useState(false);
    const [position, setPosition] = useState<TooltipPosition>("top")
    const triggerRef = useRef<HTMLDivElement>(null);

    const calculatePosition = () => {
        if (!triggerRef?.current) return;

        const rect = triggerRef.current.getBoundingClientRect();
        const space = {
            top: rect.top,
            bottom: window.innerHeight - rect.bottom,
            left: rect.left,
            right: window.innerWidth - rect.right,
        };

        const preferredOrder: TooltipPosition[] = ["top", "bottom", "left", "right"];

        let bestSide: TooltipPosition = "top";
        for (const side of preferredOrder) {
            if (space[side] > 40) { // например, минимум 40px для тултипа
                bestSide = side;
                break;
            }
        }

        setPosition(bestSide);
    };

    useEffect(() => {
        if (show) calculatePosition();
    }, [show]);

    useEffect(() => {
        window.addEventListener("resize", calculatePosition);
        return () => window.removeEventListener("resize", calculatePosition);
    }, []);

    return (
        <TooltipContext.Provider value={{show, setShow, triggerRef, position}}>
            {children}
        </TooltipContext.Provider>
    )
}

Tooltip.Trigger = ({ children }: { children: ReactNode }) => {
    const ctx = useContext(TooltipContext)!;
    return (
        <div
            ref={ctx.triggerRef}
            onMouseEnter={() => ctx.setShow(true)}
            onMouseLeave={() => ctx.setShow(false)}
            className="relative w-full h-full"
        >
            {children}
        </div>
    );
};

interface TooltipProps {
    children: ReactNode
    className?: string
    styles?: CSSProperties
}

Tooltip.Content = ({ children, className, styles }: TooltipProps) => {
    const ctx = useContext(TooltipContext)!;
    if (!ctx.show || !ctx.triggerRef.current) return null;

    const triggerRect = ctx.triggerRef.current.getBoundingClientRect();
    const OFFSET = 6;

    const style: React.CSSProperties = {
        position: "fixed",
        zIndex: 9999,
        whiteSpace: "nowrap",
    };

    if (ctx.position === "top") {
        style.left = triggerRect.left + triggerRect.width / 2;
        style.top = triggerRect.top - OFFSET;
        style.transform = "translate(-50%, -100%)";
    } else if (ctx.position === "bottom") {
        style.left = triggerRect.left + triggerRect.width / 2;
        style.top = triggerRect.bottom + OFFSET;
        style.transform = "translate(-50%, 0)";
    } else if (ctx.position === "left") {
        style.left = triggerRect.left - OFFSET;
        style.top = triggerRect.top + triggerRect.height / 2;
        style.transform = "translate(-100%, -50%)";
    } else if (ctx.position === "right") {
        style.left = triggerRect.right + OFFSET;
        style.top = triggerRect.top + triggerRect.height / 2;
        style.transform = "translate(0, -50%)";
    }

    return createPortal(
        <div className={className} style={{...style, ...styles}}>
            {children}
        </div>,
        document.body
    );
};

Tooltip.Title = ({ children, className, styles }: TooltipProps) => (
    <div className={className} style={styles}>{children}</div>
);

Tooltip.Footer = ({ children, className, styles}: TooltipProps) => (
    <div className={className} style={styles}>{children}</div>
);