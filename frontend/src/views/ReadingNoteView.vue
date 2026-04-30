<script setup lang="ts">
import DotGrid from "@/components/DotGrid.vue";
import { getReadingNote } from "@/features/reading-notes/api";
import type { ReadingNoteDetail } from "@/features/reading-notes/types";
import { appPath } from "@/utils/app-base";

const props = defineProps<{
	slug: string;
}>();

const note = ref<ReadingNoteDetail | null>(null);
const isLoading = ref(true);
const hasError = ref(false);
const themeStore = useThemeStore();
const dotBaseColor = ref("rgba(128, 128, 128, 0.2)");
const dotActiveColor = ref("#000000");

const coverSrc = computed(() =>
	note.value?.coverUrl ? appPath(note.value.coverUrl) : "",
);

const updateThemeColors = () => {
	const style = getComputedStyle(document.documentElement);
	const dotBase = style.getPropertyValue("--dot-base").trim();
	const border = style.getPropertyValue("--border").trim();
	const primary = style.getPropertyValue("--primary").trim();

	dotBaseColor.value = dotBase || border || dotBaseColor.value;
	if (primary) dotActiveColor.value = primary;
};

watch(
	[() => themeStore.themeName, () => themeStore.mode],
	() => {
		setTimeout(updateThemeColors, 50);
	},
	{ immediate: true },
);

watch(
	() => props.slug,
	async (slug) => {
		isLoading.value = true;
		hasError.value = false;
		try {
			note.value = await getReadingNote(slug);
		} catch {
			note.value = null;
			hasError.value = true;
		} finally {
			isLoading.value = false;
		}
	},
	{ immediate: true },
);

onMounted(updateThemeColors);
</script>

<template>
  <div class="relative w-full min-h-screen overflow-hidden">
    <div class="fixed inset-0 z-0">
      <DotGrid
        :dotSize="3"
        :gap="18"
        :baseColor="dotBaseColor"
        :activeColor="dotActiveColor"
        :proximity="100"
        :shockRadius="150"
        :shockStrength="3"
      />
    </div>
    <div class="relative z-10 mx-auto w-full max-w-4xl px-6 py-28">
      <a
        :href="appPath('/notes')"
        class="mb-10 inline-flex text-sm font-medium text-primary hover:text-foreground"
      >
        {{ $t("notes.back") }}
      </a>

      <div
        v-if="isLoading"
        class="rounded-lg border border-border bg-card/80 p-8 text-muted-foreground"
      >
        {{ $t("notes.loading") }}
      </div>

      <div
        v-else-if="hasError || !note"
        class="rounded-lg border border-border bg-card/80 p-8"
      >
        <h1 class="mb-3 text-3xl font-bold text-foreground">
          {{ $t("notes.notFoundTitle") }}
        </h1>
        <p class="text-muted-foreground">
          {{ $t("notes.notFoundDescription") }}
        </p>
      </div>

      <article v-else>
        <header class="mb-12 border-b border-border pb-10">
          <div class="grid gap-8 md:grid-cols-[180px_1fr]">
            <div
              v-if="note.coverUrl"
              class="aspect-[2/3] self-start overflow-hidden rounded-md border border-border bg-card/70"
            >
              <img
                :src="coverSrc"
                :alt="note.bookTitle"
                class="h-full w-full object-cover"
                referrerpolicy="no-referrer"
              />
            </div>

            <div>
              <p class="mb-4 text-sm font-medium text-primary">
                {{ note.bookTitle }}
                <span v-if="note.author" class="text-muted-foreground">
                  / {{ note.author }}
                </span>
              </p>
              <h1 class="mb-6 text-4xl font-black leading-tight text-foreground md:text-6xl">
                {{ note.title }}
              </h1>
              <div class="mb-5 flex flex-wrap gap-x-4 gap-y-2 text-sm text-muted-foreground">
                <span v-if="note.publisher">{{ note.publisher }}</span>
                <span v-if="note.bookPubDate">{{ note.bookPubDate }}</span>
                <span v-if="note.bookPages">{{ note.bookPages }} pages</span>
                <a
                  v-if="note.doubanUrl"
                  :href="note.doubanUrl"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-primary hover:text-foreground"
                >
                  Douban
                </a>
              </div>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="tag in note.tags"
                  :key="tag"
                  class="rounded-full bg-primary/10 px-3 py-1 text-xs font-medium text-primary"
                >
                  {{ tag }}
                </span>
                <span
                  v-if="note.publishedAt"
                  class="rounded-full border border-border px-3 py-1 text-xs text-muted-foreground"
                >
                  {{ note.publishedAt.slice(0, 10) }}
                </span>
              </div>
            </div>
          </div>
        </header>

        <div
          class="reading-note-content max-w-none text-lg leading-8 text-muted-foreground"
          v-html="note.contentHtml"
        ></div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.reading-note-content :deep(h1),
.reading-note-content :deep(h2),
.reading-note-content :deep(h3) {
	color: var(--foreground);
	font-weight: 800;
	line-height: 1.25;
	margin: 2rem 0 1rem;
}

.reading-note-content :deep(h2) {
	font-size: 1.75rem;
}

.reading-note-content :deep(p) {
	margin: 1rem 0;
}

.reading-note-content :deep(blockquote) {
	border-left: 3px solid var(--primary);
	color: var(--foreground);
	margin: 1.5rem 0;
	padding-left: 1rem;
}

.reading-note-content :deep(ul),
.reading-note-content :deep(ol) {
	margin: 1rem 0;
	padding-left: 1.5rem;
}

.reading-note-content :deep(a) {
	color: var(--primary);
}
</style>
