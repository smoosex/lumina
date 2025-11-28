import { defineStore } from "pinia";

export const useThemeStore = defineStore("theme", () => {
	const storedMode = localStorage.getItem("theme-mode");
	const storedThemeName = localStorage.getItem("theme-name");

	const mode = ref(
		storedMode
			? storedMode
			: window.matchMedia("(prefers-color-scheme: dark)").matches
				? "dark"
				: "light"
	);
	const themeName = ref(storedThemeName || "claude");

	const setMode = (next: string) => {
		mode.value = next;
		if (next === "dark") {
			document.documentElement.classList.add("dark");
		} else {
			document.documentElement.classList.remove("dark");
		}
		localStorage.setItem("theme-mode", next);
	};

	const setThemeName = (next: string) => {
		// Remove old theme class
		if (themeName.value) {
			document.documentElement.classList.remove(`theme-${themeName.value}`);
		}
		// Add new theme class
		themeName.value = next;
		document.documentElement.classList.add(`theme-${next}`);
		localStorage.setItem("theme-name", next);
	};

	// Initialize
	if (mode.value === "dark") {
		document.documentElement.classList.add("dark");
	} else {
		document.documentElement.classList.remove("dark");
	}
	document.documentElement.classList.add(`theme-${themeName.value}`);

	return {
		mode,
		themeName,
		setMode,
		setThemeName,
	};
});
