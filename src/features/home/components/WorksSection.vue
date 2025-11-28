<script setup lang="ts">
import { type Project } from "@/config/site";

defineProps<{
  works: Project[];
  onNavigate: (link: string) => void;
}>();
</script>

<template>
  <section class="scroll-section py-20">
    <div class="text-center mb-16">
      <h2 class="text-4xl md:text-5xl font-bold text-foreground mb-4">
        {{ $t("home.works.title") }}
      </h2>
      <div class="w-24 h-1.5 bg-primary mx-auto rounded-full"></div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
      <div
        v-for="work in works"
        :key="work.name"
        class="bg-card border border-border shadow-sm p-8 rounded-2xl group cursor-pointer relative overflow-hidden transition-all duration-300 hover:-translate-y-2 hover:border-ring hover:shadow-lg"
        @click="onNavigate(work.link)"
      >
        <div
          class="absolute inset-0 opacity-0 group-hover:opacity-10 transition-opacity duration-500 bg-linear-to-br"
          :class="work.color"
        ></div>

        <div class="relative z-10 flex flex-col h-full">
          <div
            class="mb-6 p-4 rounded-xl bg-primary/10 w-fit group-hover:scale-110 transition-transform duration-300"
          >
            <component :is="work.icon" class="w-8 h-8 text-primary" />
          </div>

          <h3
            class="text-2xl font-bold mb-3 text-foreground group-hover:text-primary transition-colors"
          >
            {{ $t("home.works." + work.name + ".title") }}
          </h3>

          <p class="text-muted-foreground mb-8 grow leading-relaxed">
            {{ $t("home.works." + work.name + ".description") }}
          </p>

          <div class="flex items-center text-primary font-medium group/link">
            <span class="mr-2 group-hover/link:mr-3 transition-all">{{
              $t("home.works.explore")
            }}</span>
            <IconLucideExternalLink v-if="work.external" class="w-4 h-4" />
            <IconLucideArrowRight v-else class="w-4 h-4" />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
