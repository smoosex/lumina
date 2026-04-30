import { createI18n } from "vue-i18n";
import en from "../locales/en.json";
import zhHans from "../locales/zh-hans.json";

const savedLocale = localStorage.getItem("user-locale") || "en";

export default createI18n({
	legacy: false,
	locale: savedLocale,
	fallbackLocale: "en",
	messages: {
		zhHans: zhHans,
		en: en,
	},
});
