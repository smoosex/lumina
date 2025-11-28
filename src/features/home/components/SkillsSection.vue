<script setup lang="ts">
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";

gsap.registerPlugin(ScrollTrigger);

interface Skill {
  name: string;
  icon: Component;
}

defineProps<{
  skills: Skill[];
}>();

const containerRef = ref<HTMLElement | null>(null);
const titleRef = ref<HTMLElement | null>(null);
const dividerRef = ref<HTMLElement | null>(null);
const skillRefs = ref<HTMLElement[]>([]);

let ctx: gsap.Context;

onMounted(() => {
  ctx = gsap.context(() => {
    const tl = gsap.timeline({
      scrollTrigger: {
        trigger: containerRef.value,
        start: "top 80%",
      },
    });

    tl.from([titleRef.value, dividerRef.value], {
      y: 50,
      opacity: 0,
      duration: 0.8,
      stagger: 0.2,
      ease: "power3.out",
    }).from(
      skillRefs.value,
      {
        y: 100,
        opacity: 0,
        duration: 1,
        stagger: 0.05,
        ease: "elastic.out(1, 0.5)",
        onComplete: startFloating,
      },
      "-=0.4"
    );
  }, containerRef.value as Element);
});

const startFloating = () => {
  skillRefs.value.forEach((el) => {
    gsap.to(el, {
      y: "random(-8, 8)",
      rotation: "random(-2, 2)",
      duration: "random(2, 4)",
      repeat: -1,
      yoyo: true,
      ease: "sine.inOut",
      delay: "random(0, 2)",
    });
  });
};

const onEnter = (el: HTMLElement) => {
  gsap.to(el, {
    scale: 1.1,
    boxShadow: "0 0 20px rgba(var(--primary), 0.4)",
    borderColor: "var(--primary)",
    duration: 0.3,
    ease: "power2.out",
    zIndex: 10,
  });
};

const onLeave = (el: HTMLElement) => {
  gsap.to(el, {
    scale: 1,
    boxShadow: "var(--shadow-sm)",
    borderColor: "var(--border)",
    duration: 0.3,
    ease: "power2.out",
    zIndex: 1,
  });
};

onUnmounted(() => {
  ctx?.revert();
});
</script>

<template>
  <section ref="containerRef" class="py-20">
    <div class="text-center mb-16">
      <h2
        ref="titleRef"
        class="text-4xl md:text-5xl font-bold text-foreground mb-4"
      >
        {{ $t("home.skills.title") }}
      </h2>
      <div
        ref="dividerRef"
        class="w-24 h-1.5 bg-primary mx-auto rounded-full"
      ></div>
    </div>

    <div class="flex flex-wrap justify-center gap-6 md:gap-8 max-w-5xl mx-auto">
      <div
        v-for="(skill, index) in skills"
        :key="skill.name"
        :ref="(el) => (skillRefs[index] = el as HTMLElement)"
        class="bg-card border border-border shadow-sm px-6 py-4 rounded-2xl flex items-center gap-4 cursor-pointer relative overflow-hidden group will-change-transform"
        @mouseenter="skillRefs[index] && onEnter(skillRefs[index])"
        @mouseleave="skillRefs[index] && onLeave(skillRefs[index])"
      >
        <!-- Hover Glow Background -->
        <div
          class="absolute inset-0 bg-primary/5 opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none"
        ></div>

        <component
          :is="skill.icon"
          class="w-8 h-8 md:w-10 md:h-10 transition-transform duration-300 group-hover:scale-110 text-foreground"
        />
        <span
          class="font-bold text-lg md:text-xl text-muted-foreground group-hover:text-foreground transition-colors"
        >
          {{ skill.name }}
        </span>
      </div>
    </div>
  </section>
</template>
