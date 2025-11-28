import pinia from "../stores";
import i18n from "./i18n";

import type { App } from "vue";

export function registerPlugins(app: App) {
	app.use(pinia).use(i18n);
	app.config.globalProperties.$t = i18n.global.t;
}
