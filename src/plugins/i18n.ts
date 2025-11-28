import { createI18n } from "vue-i18n";
import zhHans from "../locales/zh-hans.json";
import en from "../locales/en.json";

const savedLocale = localStorage.getItem("user-locale") || "zhHans";

export default createI18n({
	legacy: false,
	locale: savedLocale,
	fallbackLocale: "zhHans",
	messages: {
		zhHans: zhHans,
		en: en,
	},
});
