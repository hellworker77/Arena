import {BrowserRouter, Route, Routes} from "react-router-dom";
import {Layout} from "../../tamplates/layout/layout.tsx";
import {FirstRouteRedirectOrNotFound} from "../../tamplates/firstOrNotFound/firstOrNotFound.tsx";
import {app_routes, renderRoutes} from "../../../routing/routes.tsx";

const App = ()=> {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Layout/>}>
                    {renderRoutes(app_routes)}
                    <Route path="/" element={<FirstRouteRedirectOrNotFound />} />
                    <Route path="/*" element={<FirstRouteRedirectOrNotFound />} />
                </Route>
            </Routes>
        </BrowserRouter>
    )
}

export default App
