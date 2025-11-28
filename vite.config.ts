import { fileURLToPath, URL } from "node:url";
import tailwindcss from "@tailwindcss/vite";
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import Icons from "unplugin-icons/vite";
import { FileSystemIconLoader } from "unplugin-icons/loaders";
import IconsResolver from "unplugin-icons/resolver";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
	plugins: [
		vue(),
		tailwindcss(),
		AutoImport({
			imports: [
				"vue",
				{
					pinia: ["defineStore", "storeToRefs"],
				},
				{
					"vue-i18n": ["useI18n"],
				},
			],
			dirs: ["./src/stores/**"],
			dts: "src/types/auto-imports.d.ts",
			vueTemplate: true,
			resolvers: [
				IconsResolver({
					prefix: "Icon",
					enabledCollections: ["skill-icons", "lucide"],
					customCollections: ["sys-icons"],
				}),
			],
		}),
		Components({
			dirs: ["src/components"],
			dts: "src/types/components.d.ts",
			deep: true,
			extensions: ["vue"],
			resolvers: [
				IconsResolver({
					prefix: "Icon",
					enabledCollections: ["skill-icons", "lucide"],
					customCollections: ["sys-icons"],
				}),
			],
		}),
		Icons({
			autoInstall: true,
			customCollections: {
				"sys-icons": FileSystemIconLoader("./src/assets/icons", (svg) =>
					svg.replace(/^<svg /, '<svg fill="currentColor" '),
				),
			},
		}),
	],
	resolve: {
		alias: {
			"@": fileURLToPath(new URL("./src", import.meta.url)),
		},
	},
});
