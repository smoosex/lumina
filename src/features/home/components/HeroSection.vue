<script setup lang="ts">
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";
const ParticleModel = defineAsyncComponent(() => import("./ParticleModel.vue"));

import { siteConfig } from "@/config/site";

gsap.registerPlugin(ScrollTrigger);

defineProps<{
  onNavigate: () => void;
  onWorksClick?: () => void;
}>();

const heroRef = ref<HTMLElement | null>(null);
const modelRef = ref<HTMLElement | null>(null);
const isVisible = ref(true);
let trigger: ScrollTrigger | null = null;

onMounted(() => {
  if (modelRef.value) {
    trigger = ScrollTrigger.create({
      trigger: modelRef.value,
      start: "top bottom",
      end: "bottom top",
      onUpdate: (self) => {
        if (self.progress > 0.3 && self.progress < 0.7) {
          isVisible.value = true;
        } else {
          isVisible.value = false;
        }
      },
    });
  }
});

onUnmounted(() => {
  if (trigger) {
    trigger.kill();
  }
});
</script>

<template>
  <section
    ref="heroRef"
    class="min-h-[80vh] grid grid-cols-1 lg:grid-cols-2 gap-12 items-center w-full"
  >
    <!-- Text Content -->
    <div class="hero-content flex flex-col justify-center items-start z-10">
      <h2
        class="text-xl md:text-2xl font-medium text-primary mb-4 tracking-wide uppercase"
      >
        <shuffle
          :key="$t('home.greeting')"
          :text="$t('home.greeting')"
          shuffle-direction="right"
          :duration="0.35"
          animation-mode="evenodd"
          :shuffle-times="1"
          ease="power3.out"
          :stagger="0.03"
          :threshold="0.1"
          :trigger-once="true"
          :trigger-on-hover="true"
          :respect-reduced-motion="true"
        />
      </h2>
      <h1
        class="text-4xl sm:text-5xl md:text-8xl font-black text-foreground leading-tight mb-8"
      >
        {{ $t("home.titlePrefix") }} <br />
        <span
          class="bg-linear-to-r from-primary to-accent bg-clip-text text-transparent"
        >
          {{ siteConfig.profile.name }}
        </span>
      </h1>
      <p class="text-xl text-muted-foreground max-w-2xl leading-relaxed mb-10">
        {{ $t("home.description") }}
      </p>
      <div class="flex flex-wrap gap-4">
        <button
          @click="onNavigate"
          class="px-6 py-3 md:px-8 md:py-4 rounded-full bg-primary text-primary-foreground font-bold text-base md:text-lg hover:opacity-90 transition shadow-lg hover:shadow-xl transform hover:-translate-y-1"
        >
          {{ $t("home.hero.primaryButton") }}
        </button>
        <button
          @click="onWorksClick"
          class="px-6 py-3 md:px-8 md:py-4 rounded-full bg-secondary/10 border-2 border-border text-foreground font-bold text-base md:text-lg hover:bg-secondary transition flex items-center gap-2 cursor-pointer"
        >
          <IconLucideDog class="w-5 h-5" />
          {{ $t("home.hero.secondaryButton") }}
        </button>
      </div>
    </div>

    <!-- 3D Model -->
    <div
      ref="modelRef"
      class="w-full h-[500px] flex items-center justify-center relative z-0"
    >
      <ParticleModel :active="isVisible" />
    </div>
  </section>
</template>
