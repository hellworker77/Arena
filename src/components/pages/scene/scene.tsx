import {PixiCanvas} from "../../tamplates/pixi/pixiCanvas.tsx";
import {Page} from "../../tamplates/page/page.tsx";

export const Scene = () => {

    return (
        <Page>
            <Page.Title title={"scene"}/>
            <Page.Body>
                <PixiCanvas />
            </Page.Body>
        </Page>
    )
}