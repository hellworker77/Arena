import type { PropsWithChildren } from "react";
import {useTranslation} from "react-i18next";

export const Page = ({ children }: PropsWithChildren) => {
    return (
        <main className="flex-1 flex flex-col items-center justify-center w-full rounded-xl shadow-[0_0_15px_rgba(255,0,0,0.3)] p-6 bg-gradient-to-b from-gray-900 via-black to-gray-950 border border-red-900 select-none">
            {children}
        </main>
    );
};

type TitleProps = {
    title: string;
};

const Title = ({ title }: TitleProps) => {
    const {t} = useTranslation()
    return (
        <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-bold text-red-500 tracking-wider drop-shadow-[0_0_8px_rgba(255,0,0,0.7)]">
                {t(title)}
            </h2>
        </div>
    );
};

const Body = ({ children }: PropsWithChildren) => {
    return (
        <div className="text-gray-300 mb-6 leading-relaxed">
            {children}
        </div>
    );
};

const Footer = ({ children }: PropsWithChildren) => {
    return (
        <div className="text-right text-sm text-gray-500 border-t border-red-900 pt-2">
            {children}
        </div>
    );
};

Page.Title = Title;
Page.Body = Body;
Page.Footer = Footer;
