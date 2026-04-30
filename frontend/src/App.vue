<script setup lang="ts">
import DefaultLayout from "@/layouts/DefaultLayout.vue";
import HomeView from "@/views/HomeView.vue";
import ReadingNotesView from "@/views/ReadingNotesView.vue";
import ReadingNoteView from "@/views/ReadingNoteView.vue";
import { appBase, routePath } from "@/utils/app-base";

const currentPath = ref(routePath());

const noteSlug = computed(() => {
	const match = currentPath.value.match(/^\/notes\/([^/]+)\/?$/);
	return match?.[1] ?? "";
});

const isNotesIndex = computed(() => /^\/notes\/?$/.test(currentPath.value));

const syncPath = () => {
	currentPath.value = routePath();
};

const navigateTo = (url: URL) => {
	window.history.pushState({}, "", `${url.pathname}${url.search}${url.hash}`);
	syncPath();
	window.dispatchEvent(new CustomEvent("lumina:navigated"));
};

const handleInternalLinkClick = (event: MouseEvent) => {
	if (
		event.defaultPrevented ||
		event.button !== 0 ||
		event.metaKey ||
		event.ctrlKey ||
		event.shiftKey ||
		event.altKey
	) {
		return;
	}

	const link = (event.target as Element | null)?.closest("a[href]");
	if (!link) return;

	const target = link.getAttribute("target");
	const href = link.getAttribute("href");
	if (!href || target === "_blank" || link.hasAttribute("download")) return;

	const url = new URL(href, window.location.origin);
	if (url.origin !== window.location.origin) return;
	if (
		appBase !== "/" &&
		url.pathname !== appBase.slice(0, -1) &&
		!url.pathname.startsWith(appBase)
	) {
		return;
	}

	event.preventDefault();
	navigateTo(url);
};

onMounted(() => {
	window.addEventListener("popstate", syncPath);
	document.addEventListener("click", handleInternalLinkClick);
});

onUnmounted(() => {
	window.removeEventListener("popstate", syncPath);
	document.removeEventListener("click", handleInternalLinkClick);
});
</script>

<template>
  <default-layout>
    <ReadingNoteView v-if="noteSlug" :slug="noteSlug" />
    <ReadingNotesView v-else-if="isNotesIndex" />
    <home-view v-else />
  </default-layout>
</template>

<style scoped></style>
