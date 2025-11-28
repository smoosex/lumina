<script setup lang="ts">
import gsap from "gsap";
import { siteConfig, aboutCode } from "@/config/site";

const codeContainer = ref<HTMLElement | null>(null);
const scrollArea = ref<HTMLElement | null>(null);
const displayedCode = ref("");
const cursorVisible = ref(true);

const fullCode = aboutCode(siteConfig.profile);

let typeTween: gsap.core.Tween | null = null;
let cursorInterval: number | null = null;
let observer: IntersectionObserver | null = null;

onMounted(() => {
  cursorInterval = window.setInterval(() => {
    cursorVisible.value = !cursorVisible.value;
  }, 500);

  if (codeContainer.value) {
    const progress = { index: 0 };
    typeTween = gsap.to(progress, {
      index: fullCode.length,
      duration: 10,
      ease: "none",
      repeat: -1,
      repeatDelay: 2,
      onUpdate: () => {
        displayedCode.value = fullCode.substring(0, Math.ceil(progress.index));
        if (scrollArea.value) {
          scrollArea.value.scrollTop = scrollArea.value.scrollHeight;
        }
      },
      paused: true,
    });

    observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.intersectionRatio >= 0.8) {
            typeTween?.play();
          } else if (entry.intersectionRatio < 0.5) {
            typeTween?.pause();
          }
        });
      },
      {
        threshold: [0.5, 0.8],
      }
    );

    observer.observe(codeContainer.value);
  }
});

onUnmounted(() => {
  if (cursorInterval) clearInterval(cursorInterval);
  if (typeTween) typeTween.kill();
  if (observer) observer.disconnect();
});
</script>

<template>
  <section
    class="scroll-section grid grid-cols-1 md:grid-cols-2 gap-12 items-center"
  >
    <!-- Left Content -->
    <div>
      <h2 class="text-4xl font-bold text-foreground mb-6">
        {{ $t("home.about.title") }}
      </h2>
      <div class="w-20 h-1 bg-primary mb-8"></div>
      <p class="text-lg text-muted-foreground leading-relaxed">
        {{ $t("home.about.description") }}
      </p>
    </div>

    <!-- Right Code Window -->
    <div
      ref="codeContainer"
      class="relative h-80 md:h-96 rounded-2xl overflow-hidden flex flex-col font-mono text-sm md:text-base bg-card border border-border shadow-sm"
    >
      <!-- Mac-style window controls (Header Bar) -->
      <div
        class="h-12 flex items-center px-4 shrink-0 border-b border-border bg-muted/30 z-20"
      >
        <div class="flex gap-2">
          <div class="w-3 h-3 rounded-full bg-red-500"></div>
          <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
          <div class="w-3 h-3 rounded-full bg-green-500"></div>
        </div>
      </div>

      <!-- Scrollable Code Area -->
      <div
        ref="scrollArea"
        class="w-full flex-1 overflow-y-auto p-6"
        style="scrollbar-width: none; -ms-overflow-style: none"
      >
        <pre><code class="text-primary">{{ displayedCode }}<span v-show="cursorVisible" class="text-foreground">|</span></code></pre>
      </div>
    </div>
  </section>
</template>
