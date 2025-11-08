import i18n from "i18next";
import {initReactI18next} from "react-i18next";

import enGlobal from "./en/global.json"
import enRouting from "./en/routing.json"
import enItem from "./en/item.json"

import ruGlobal from "./ru/global.json"
import ruRouting from "./ru/routing.json"
import ruItem from "./ru/item.json"


export const resources = {
    en: {
        translation: {...enGlobal, ...enRouting, ...enItem}
    },
    ru: {
        translation: {...ruGlobal, ...ruRouting, ...ruItem}
    }
};

i18n.use(initReactI18next).init({
    resources,
    lng: "en",
    fallbackLng: "en",
    interpolation: {escapeValue: false}
}).then(console.log);

export default i18n;
