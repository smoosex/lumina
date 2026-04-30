<script setup lang="ts">
import DotGrid from "@/components/DotGrid.vue";
import { listReadingNotes } from "@/features/reading-notes/api";
import type { ReadingNoteListItem } from "@/features/reading-notes/types";

const notes = ref<ReadingNoteListItem[]>([]);
const selectedTag = ref("all");
const isLoading = ref(true);
const hasError = ref(false);
const themeStore = useThemeStore();
const dotBaseColor = ref("rgba(128, 128, 128, 0.2)");
const dotActiveColor = ref("#000000");

const tags = computed(() => {
	const counts = new Map<string, number>();
	for (const note of notes.value) {
		for (const tag of note.tags) {
			counts.set(tag, (counts.get(tag) ?? 0) + 1);
		}
	}
	return [...counts.entries()]
		.map(([name, count]) => ({ name, count }))
		.sort((a, b) => a.name.localeCompare(b.name));
});

const visibleNotes = computed(() => {
	if (selectedTag.value === "all") return notes.value;
	return notes.value.filter((note) => note.tags.includes(selectedTag.value));
});

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

onMounted(async () => {
	updateThemeColors();
	try {
		const response = await listReadingNotes();
		notes.value = response.items;
	} catch {
		hasError.value = true;
	} finally {
		isLoading.value = false;
	}
});
</script>

<template>
  <div class="relative min-h-screen w-full overflow-hidden">
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

    <main class="relative z-10 mx-auto w-full max-w-7xl px-6 py-28">
      <header class="mb-12 max-w-3xl">
        <p class="mb-4 text-sm font-bold uppercase tracking-wider text-primary">
          {{ $t("notes.eyebrow") }}
        </p>
        <h1 class="text-5xl font-black leading-tight text-foreground md:text-7xl">
          {{ $t("home.notes.title") }}
        </h1>
        <div class="mt-5 h-1.5 w-24 rounded-full bg-primary"></div>
        <p class="mt-6 max-w-2xl text-lg leading-8 text-muted-foreground">
          {{ $t("home.notes.description") }}
        </p>
      </header>

      <div class="grid gap-8 lg:grid-cols-[260px_1fr]">
        <aside class="space-y-8">
          <div class="rounded-lg border border-border bg-card/45 p-3 backdrop-blur-sm">
            <button
              class="flex w-full items-center justify-between rounded-md px-4 py-3 text-left text-sm font-semibold transition"
              :class="selectedTag === 'all'
                ? 'bg-primary/10 text-primary'
                : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground'"
              @click="selectedTag = 'all'"
            >
              <span class="inline-flex items-center gap-3">
                <IconLucideLayoutGrid class="h-4 w-4" />
                {{ $t("notes.all") }}
              </span>
              <span>{{ notes.length }}</span>
            </button>

            <button
              v-for="tag in tags"
              :key="tag.name"
              class="mt-1 flex w-full items-center justify-between rounded-md px-4 py-3 text-left text-sm font-semibold transition"
              :class="selectedTag === tag.name
                ? 'bg-primary/10 text-primary'
                : 'text-muted-foreground hover:bg-muted/60 hover:text-foreground'"
              @click="selectedTag = tag.name"
            >
              <span class="inline-flex items-center gap-3">
                <IconLucideTag class="h-4 w-4" />
                {{ tag.name }}
              </span>
              <span>{{ tag.count }}</span>
            </button>
          </div>

          <div class="border-t border-border pt-8 text-muted-foreground">
            <IconLucideQuote class="mb-4 h-7 w-7 text-primary" />
            <p class="text-base leading-7">
              {{ $t("notes.quote") }}
            </p>
            <p class="mt-4 text-sm">
              {{ $t("notes.quoteAuthor") }}
            </p>
          </div>
        </aside>

        <section>
          <div v-if="isLoading" class="grid gap-5 xl:grid-cols-2">
            <div
              v-for="index in 4"
              :key="index"
              class="h-64 rounded-lg border border-border bg-card/50"
            ></div>
          </div>

          <div
            v-else-if="hasError || notes.length === 0"
            class="rounded-lg border border-border bg-card/65 p-10 backdrop-blur-sm"
          >
            <h2 class="text-2xl font-bold text-foreground">
              {{ $t("notes.emptyTitle") }}
            </h2>
            <p class="mt-3 max-w-xl text-muted-foreground">
              {{ $t("home.notes.empty") }}
            </p>
          </div>

          <div v-else class="grid gap-5 xl:grid-cols-2">
            <a
              v-for="note in visibleNotes"
              :key="note.slug"
              :href="`/notes/${note.slug}`"
              class="group grid gap-5 rounded-lg border border-border bg-card/55 p-5 backdrop-blur-sm transition hover:border-primary/60 hover:bg-card/75 md:grid-cols-[132px_1fr]"
            >
              <div class="relative flex aspect-[2/3] w-full items-center justify-center self-start overflow-hidden rounded-md border border-border bg-primary/10 text-primary">
                <img
                  v-if="note.coverUrl"
                  :src="note.coverUrl"
                  :alt="note.bookTitle"
                  class="h-full w-full object-cover"
                  loading="lazy"
                  referrerpolicy="no-referrer"
                />
                <IconLucideBookOpen v-else class="h-12 w-12" />
              </div>

              <div class="min-w-0">
                <p class="text-sm text-muted-foreground">
                  {{ note.bookTitle }}
                  <span v-if="note.author"> / {{ note.author }}</span>
                </p>
                <h2 class="mt-2 text-2xl font-bold leading-snug text-foreground transition group-hover:text-primary">
                  {{ note.title }}
                </h2>
                <p class="mt-4 line-clamp-3 text-sm leading-6 text-muted-foreground">
                  {{ note.excerpt || $t("notes.noExcerpt") }}
                </p>

                <div class="mt-5 flex flex-wrap items-center gap-2 border-t border-border/70 pt-4">
                  <span
                    v-for="tag in note.tags"
                    :key="tag"
                    class="rounded-md border border-border bg-muted/40 px-2.5 py-1 text-xs text-muted-foreground"
                  >
                    {{ tag }}
                  </span>
                  <span class="ml-auto inline-flex items-center gap-1 text-sm font-semibold text-primary">
                    {{ $t("home.notes.read") }}
                    <IconLucideArrowRight class="h-4 w-4 transition group-hover:translate-x-1" />
                  </span>
                </div>
              </div>
            </a>
          </div>
        </section>
      </div>
    </main>
  </div>
</template>
