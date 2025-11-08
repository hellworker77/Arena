import {Nav} from "../nav/nav.tsx";
import background from "../../../assets/bg.jpeg"
import {BreadCrumb} from "../breadcrumb/breadCrumb.tsx";
import {Outlet} from "react-router-dom";
import {Side} from "../side/side.tsx";

export const Layout = () => {
    return (
        <div style={{
            backgroundImage: `url(${background})`,
            backgroundSize: 'cover',
            backgroundPosition: 'center',
            backgroundRepeat: 'no-repeat',
        }}
             className="flex flex-col min-h-screen bg-transparent text-gray-200 select-none">

            <Nav/>
            <div className="flex flex-1">
                <Side/>
                <div
                    className="flex-1 flex flex-col p-6 bg-black/10">
                    <div className="mb-4">
                        <BreadCrumb/>
                    </div>
                    <Outlet/>
                </div>
            </div>
        </div>
    );
};
