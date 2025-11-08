import type {PropsWithChildren} from "react";
import {useTranslation} from "react-i18next";

export const Page = ({children}: PropsWithChildren) => {
    return (
        <main
            className="flex-1 flex flex-col items-center justify-center w-full h-full rounded-xl shadow-[0_0_15px_rgba(0,0,0,0.3)] bg-black/10 border border-black/20 select-none overflow-hidden">
            {children}
        </main>
    );
};

type TitleProps = {
    title: string;
};

const Title = ({title}: TitleProps) => {
    const {t} = useTranslation();
    return (
        <div className="flex items-start text-left justify-between p-4">
            <h2 className="text-2xl font-bold text-red-500 tracking-wider drop-shadow-[0_0_8px_rgba(255,0,0,0.7)]">
                {t(title)}
            </h2>
        </div>
    );
};

const Body = ({children}: PropsWithChildren) => {
    return (
        <div className="flex-1 w-full h-full m-0 p-0 flex items-center justify-center overflow-hidden">
            {children}
        </div>
    );
};

const Footer = ({children}: PropsWithChildren) => {
    return (
        <div className="text-right text-sm text-gray-500 border-t border-red-900 pt-2 px-4">
            {children}
        </div>
    );
};

Page.Title = Title;
Page.Body = Body;
Page.Footer = Footer;