import {PropsWithChildren} from "react";

interface ExocetTitleProps extends PropsWithChildren {
    className?: string;
    color?: string;
}

export const ExocetTitle = ({children, className, color="#fff"}: ExocetTitleProps) => {
    if (typeof children !== "string") {
        throw new Error("ExocetTitle must be a string");
    }

    const text = String(children);
    const first = text[0] ?? "";
    const rest = text.slice(1);

    return (
        <span className={`font-exocet ${className}`} style={{color}}>
            <span className="exocet-caps">{first}</span>
            {rest}
        </span>
    );
}