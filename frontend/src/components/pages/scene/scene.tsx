import {Page} from "../../tamplates/page/page.tsx";
import {Canvas} from "../../tamplates/game/canvas.tsx";

export const Scene = () => {

    return (
        <Page>

            <Page.Body>
                <Canvas />
            </Page.Body>
            <Page.Footer>
                <div className="h-40">

                </div>
            </Page.Footer>
        </Page>
    )
}