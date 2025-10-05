import {useTranslation} from "react-i18next";
import {Page} from "../../tamplates/page/page.tsx";
import {Button} from "../../tamplates/button/button.tsx";

export const Home = () => {
    const {t, i18n} = useTranslation();

    return (
        <Page>
            <Page.Title title={"home"}/>
            <Page.Body>
                <Button onClick={() => i18n.changeLanguage(i18n.language === "en" ? "ru" : "en")}>
                    {t("switch_language")}
                </Button>
            </Page.Body>
        </Page>
    );
};
