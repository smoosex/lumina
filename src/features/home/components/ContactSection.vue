<script setup lang="ts">
import { siteConfig } from "@/config/site";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";

gsap.registerPlugin(ScrollTrigger);

const containerRef = ref<HTMLElement | null>(null);
const cardsRef = ref<HTMLElement[]>([]);

let ctx: gsap.Context;

onMounted(() => {
  ctx = gsap.context(() => {
    gsap.from(cardsRef.value, {
      scrollTrigger: {
        trigger: containerRef.value,
        start: "top 80%",
      },
      y: 100,
      opacity: 0,
      duration: 1,
      stagger: 0.2,
      ease: "power3.out",
    });
  }, containerRef.value as Element);
});

const onMouseMove = (e: MouseEvent, index: number) => {
  const card = cardsRef.value[index];
  if (!card) return;

  const rect = card.getBoundingClientRect();
  const x = e.clientX - rect.left;
  const y = e.clientY - rect.top;

  const centerX = rect.width / 2;
  const centerY = rect.height / 2;

  const rotateX = ((y - centerY) / centerY) * -10;
  const rotateY = ((x - centerX) / centerX) * 10;

  gsap.to(card, {
    rotationX: rotateX,
    rotationY: rotateY,
    scale: 1.05,
    duration: 0.4,
    ease: "power2.out",
    transformPerspective: 1000,
    transformOrigin: "center",
  });
};

const onMouseLeave = (index: number) => {
  const card = cardsRef.value[index];
  if (!card) return;

  gsap.to(card, {
    rotationX: 0,
    rotationY: 0,
    scale: 1,
    duration: 0.7,
    ease: "elastic.out(1, 0.5)",
  });
};

const onTouchMove = (e: TouchEvent, index: number) => {
  const card = cardsRef.value[index];
  if (!card || e.touches.length === 0) return;

  const touch = e.touches[0];
  if (!touch) return;

  const rect = card.getBoundingClientRect();
  const x = touch.clientX - rect.left;
  const y = touch.clientY - rect.top;

  const centerX = rect.width / 2;
  const centerY = rect.height / 2;

  const rotateX = ((y - centerY) / centerY) * -10;
  const rotateY = ((x - centerX) / centerX) * 10;

  gsap.to(card, {
    rotationX: rotateX,
    rotationY: rotateY,
    scale: 1.05,
    duration: 0.4,
    ease: "power2.out",
    transformPerspective: 1000,
    transformOrigin: "center",
  });
};

onUnmounted(() => {
  ctx?.revert();
});
</script>

<template>
  <section ref="containerRef" class="scroll-section py-20">
    <div class="text-center mb-16">
      <h2 class="text-4xl md:text-5xl font-bold text-foreground mb-4">
        {{ $t("home.contact.title") }}
      </h2>
      <div class="w-24 h-1.5 bg-primary mx-auto rounded-full"></div>
    </div>

    <div
      class="flex flex-col md:flex-row justify-center gap-8 md:gap-12 max-w-4xl mx-auto px-4"
    >
      <!-- Email Card -->
      <a
        :href="`mailto:${siteConfig.profile.email}`"
        :ref="(el) => (cardsRef[0] = el as HTMLElement)"
        class="bg-card border border-border shadow-sm flex-1 p-10 rounded-3xl flex flex-col items-center gap-6 text-center group relative overflow-hidden hover:border-ring hover:shadow-xl transition-colors"
        @mousemove="(e) => onMouseMove(e, 0)"
        @mouseleave="onMouseLeave(0)"
        @touchmove="(e) => onTouchMove(e, 0)"
        @touchend="onMouseLeave(0)"
      >
        <div
          class="absolute inset-0 bg-linear-to-br from-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"
        ></div>

        <div
          class="w-20 h-20 rounded-2xl bg-primary/10 flex items-center justify-center group-hover:scale-110 transition-transform duration-500 shadow-inner"
        >
          <IconLucideMail class="w-10 h-10 text-primary" />
        </div>

        <div class="relative z-10">
          <h3 class="text-2xl font-bold mb-2 text-foreground">Email Me</h3>
          <p class="text-muted-foreground">{{ siteConfig.profile.email }}</p>
        </div>
      </a>

      <!-- GitHub Card -->
      <a
        :href="siteConfig.profile.github"
        target="_blank"
        :ref="(el) => (cardsRef[1] = el as HTMLElement)"
        class="bg-card border border-border shadow-sm flex-1 p-10 rounded-3xl flex flex-col items-center gap-6 text-center group relative overflow-hidden hover:border-ring hover:shadow-xl transition-colors"
        @mousemove="(e) => onMouseMove(e, 1)"
        @mouseleave="onMouseLeave(1)"
        @touchmove="(e) => onTouchMove(e, 1)"
        @touchend="onMouseLeave(1)"
      >
        <div
          class="absolute inset-0 bg-linear-to-br from-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"
        ></div>

        <div
          class="w-20 h-20 rounded-2xl bg-primary/10 flex items-center justify-center group-hover:scale-110 transition-transform duration-500 shadow-inner"
        >
          <IconLucideGithub class="w-10 h-10 text-primary" />
        </div>

        <div class="relative z-10">
          <h3 class="text-2xl font-bold mb-2 text-foreground">GitHub</h3>
          <p class="text-muted-foreground">Check out my code</p>
        </div>
      </a>
    </div>
  </section>
</template>
