import { Outlet } from "react-router-dom";
import { BreadCrumb } from "../breadcrumb/breadCrumb.tsx";
import { Side } from "../side/side.tsx";
import { Nav } from "../nav/nav.tsx";

export const Layout = () => {
    return (
        <div className="flex flex-col min-h-screen bg-black text-gray-200 select-none">
            <Nav />

            <div className="flex flex-1">
                <Side />

                <div className="flex-1 flex flex-col p-6 bg-gradient-to-b from-gray-900 via-black to-gray-950 border-l border-red-900 shadow-inner">
                    <div className="mb-4">
                        <BreadCrumb />
                    </div>
                    <Outlet />
                </div>
            </div>
        </div>
    );
};
