<script setup lang="ts">
import DotGrid from "@/components/DotGrid.vue";
import { gsap } from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import { ScrollToPlugin } from "gsap/ScrollToPlugin";

import { works, skills } from "@/config/site";

// Feature Components
import {
  HeroSection,
  AboutSection,
  SkillsSection,
  WorksSection,
  ContactSection,
} from "@/features/home";

gsap.registerPlugin(ScrollTrigger, ScrollToPlugin);

const themeStore = useThemeStore();

const navigateTo = (link: string) => {
  window.open(link, "_blank", "noopener,noreferrer");
};

const scrollToAbout = () => {
  const aboutSection = document.getElementById("about-section");
  if (aboutSection) {
    // Calculate position to center the section
    const elementHeight = aboutSection.offsetHeight;
    const windowHeight = window.innerHeight;
    const offset = (windowHeight - elementHeight) / 2;

    // Use GSAP for smooth scrolling with control over duration and easing
    gsap.to(window, {
      duration: 1.5,
      scrollTo: { y: aboutSection, offsetY: offset },
      ease: "power3.inOut",
    });
  }
};

const scrollToWorks = () => {
  const worksSection = document.getElementById("works-section");
  if (worksSection) {
    const elementHeight = worksSection.offsetHeight;
    const windowHeight = window.innerHeight;
    const offset = (windowHeight - elementHeight) / 2;

    gsap.to(window, {
      duration: 1.5,
      scrollTo: { y: worksSection, offsetY: offset },
      ease: "power3.inOut",
    });
  }
};

// Theme colors for DotGrid
const dotBaseColor = ref("rgba(128, 128, 128, 0.2)");
const dotActiveColor = ref("#000000");

const updateThemeColors = () => {
  const style = getComputedStyle(document.documentElement);
  const border = style.getPropertyValue("--border").trim();
  const primary = style.getPropertyValue("--primary").trim();

  if (border) {
    dotBaseColor.value = border;
  }

  if (primary) {
    dotActiveColor.value = primary;
  }
};

watch(
  [() => themeStore.themeName, () => themeStore.mode],
  () => {
    setTimeout(updateThemeColors, 50);
  },
  { immediate: true }
);

onMounted(async () => {
  updateThemeColors();
  await nextTick();

  // Hero Animation
  gsap.from(".hero-content > *", {
    y: 50,
    opacity: 0,
    duration: 1,
    stagger: 0.2,
    ease: "power3.out",
  });

  // Scroll Animations
  const sections = document.querySelectorAll(".scroll-section");
  sections.forEach((section) => {
    gsap.from(section.children, {
      scrollTrigger: {
        trigger: section,
        start: "top 80%",
        toggleActions: "play none none reverse",
      },
      y: 50,
      opacity: 0,
      duration: 0.8,
      stagger: 0.2,
      ease: "power3.out",
    });
  });
});
</script>

<template>
  <div class="relative w-full min-h-screen overflow-hidden">
    <!-- Background DotGrid -->
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

    <!-- Main Content -->
    <div
      class="relative z-10 w-full max-w-7xl mx-auto px-6 py-20 flex flex-col gap-32"
    >
      <HeroSection :onNavigate="scrollToAbout" :onWorksClick="scrollToWorks" />

      <AboutSection id="about-section" />

      <SkillsSection :skills="skills" />

      <WorksSection
        id="works-section"
        :works="works"
        :onNavigate="navigateTo"
      />

      <ContactSection />
    </div>
  </div>
</template>
