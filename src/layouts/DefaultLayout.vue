<script setup lang="ts">
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import { ScrollToPlugin } from "gsap/ScrollToPlugin";

gsap.registerPlugin(ScrollTrigger, ScrollToPlugin);

const themeStore = useThemeStore();
const isDark = computed(() => themeStore.mode === "dark");
const isScrolled = ref(false);
let scrollTrigger: ScrollTrigger | null = null;

const toggleThemeMode = () => {
  themeStore.setMode(isDark.value ? "light" : "dark");
};

const { locale } = useI18n();

const localeOptions = [
  { label: "简体中文", value: "zhHans" },
  { label: "English", value: "en" },
];

const themeOptions = [
  { label: "Everforest Hard", value: "everforest-hard" },
  { label: "Everforest Medium", value: "everforest-medium" },
  { label: "Everforest Soft", value: "everforest-soft" },
];

const currentLocale = computed({
  get: () => locale.value,
  set: (val: string) => {
    locale.value = val;
    localStorage.setItem("user-locale", val);
    activeMenu.value = null;
  },
});

const currentThemeName = computed({
  get: () => themeStore.themeName,
  set: (val: string) => {
    themeStore.setThemeName(val);
    activeMenu.value = null;
  },
});

const middleTextRef = ref<HTMLElement | null>(null);
const activeMenu = ref<"theme" | "lang" | null>(null);

const toggleMenu = (menu: "theme" | "lang", event: Event) => {
  event.stopPropagation();
  if (activeMenu.value === menu) {
    activeMenu.value = null;
  } else {
    activeMenu.value = menu;
  }
};

const closeMenus = () => {
  activeMenu.value = null;
};

const navItems = [
  { key: "about", i18n: "home.about.title" },
  { key: "skills", i18n: "home.skills.title" },
  { key: "works", i18n: "home.works.title" },
  { key: "contact", i18n: "home.contact.title" },
];

const scrollToSection = (key: string) => {
  const el = document.getElementById(`${key}-section`);
  if (!el) return;
  const elementHeight = el.offsetHeight;
  const windowHeight = window.innerHeight;
  const offset = (windowHeight - elementHeight) / 2;
  gsap.to(window, {
    duration: 1.5,
    scrollTo: { y: el, offsetY: offset },
    ease: "power3.inOut",
  });
};

onMounted(() => {
  themeStore.setMode(themeStore.mode);
  themeStore.setThemeName(themeStore.themeName);

  document.addEventListener("click", closeMenus);

  scrollTrigger = ScrollTrigger.create({
    trigger: document.body,
    start: "top top",
    end: "bottom bottom",
    onUpdate: (self) => {
      isScrolled.value = self.scroll() > 50;
    },
  });

  if (middleTextRef.value) {
    gsap.from(middleTextRef.value, {
      width: 0,
      opacity: 0,
      duration: 1,
      ease: "power3.out",
      delay: 0.5,
    });
  }
});

watch(isScrolled, (scrolled) => {
  if (scrolled) activeMenu.value = null;

  if (!middleTextRef.value) return;
  if (scrolled) {
    gsap.to(middleTextRef.value, {
      width: 0,
      opacity: 0,
      duration: 0.5,
      ease: "power3.inOut",
      overwrite: true,
    });
  } else {
    gsap.to(middleTextRef.value, {
      width: "auto",
      opacity: 1,
      duration: 0.5,
      ease: "power3.inOut",
      overwrite: true,
    });
  }
});

onUnmounted(() => {
  if (scrollTrigger) {
    scrollTrigger.kill();
  }
  document.removeEventListener("click", closeMenus);
});
</script>

<template>
  <div
    class="min-h-screen w-full relative overflow-hidden transition-colors duration-500 bg-background"
  >
    <!-- Header -->
    <header
      class="fixed top-0 left-0 right-0 z-50 px-6 py-2 flex items-center gap-3 transition-all duration-300"
      :class="[
        isScrolled
          ? 'bg-background border-border shadow-sm'
          : 'bg-transparent border-transparent',
      ]"
    >
      <!-- Logo -->
      <div
        class="text-xl leading-none font-bold tracking-wider drop-shadow-sm flex items-center whitespace-nowrap overflow-hidden text-foreground w-40 shrink-0"
      >
        <span>H</span>
        <span
          ref="middleTextRef"
          class="inline-block overflow-hidden whitespace-nowrap"
          >ELLO \ WORL</span
        >
        <span>D</span>
      </div>

      <!-- Navigation -->
      <nav class="hidden md:flex items-center gap-1 flex-1 justify-center">
        <button
          v-for="item in navItems"
          :key="item.key"
          @click="scrollToSection(item.key)"
          class="px-3 py-1.5 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors rounded-full cursor-pointer"
        >
          {{ $t(item.i18n) }}
        </button>
      </nav>

      <div class="flex items-center gap-2"> 
        <!-- Theme Selector -->
        <div class="relative">
          <button
            @click="(e) => toggleMenu('theme', e)"
            class="p-1.5 rounded-full transition-all cursor-pointer bg-card border border-border text-foreground shadow-sm hover:bg-accent hover:text-accent-foreground"
            :class="{ 'ring-2 ring-ring': activeMenu === 'theme' }"
          >
            <IconLucidePalette class="w-4 h-4" />
          </button>

          <div
            class="absolute right-0 mt-2 w-40 rounded-xl z-50 transition-all transform origin-top-right bg-card border border-border shadow-lg overflow-hidden"
            :class="[
              activeMenu === 'theme'
                ? 'opacity-100 visible scale-100'
                : 'opacity-0 invisible scale-95 pointer-events-none',
            ]"
            @click.stop
          >
            <button
              v-for="opt in themeOptions"
              :key="opt.value"
              @click="currentThemeName = opt.value"
              class="w-full text-left px-4 py-2 text-sm transition-colors flex items-center justify-between cursor-pointer"
              :class="{
                'bg-accent text-accent-foreground font-bold':
                  currentThemeName === opt.value,
                'text-foreground hover:bg-muted':
                  currentThemeName !== opt.value,
              }"
            >
              {{ opt.label }}
            </button>
          </div>
        </div>

        <!-- Language Switcher -->
        <div class="relative">
          <button
            @click="(e) => toggleMenu('lang', e)"
            class="flex items-center gap-1.5 px-2 py-1.5 rounded-full transition-all cursor-pointer bg-card border border-border text-foreground shadow-sm hover:bg-accent hover:text-accent-foreground"
            :class="{ 'ring-2 ring-ring': activeMenu === 'lang' }"
          >
            <IconLucideLanguages class="w-4 h-4" />
            <span class="text-xs font-medium uppercase hidden sm:inline">{{
              currentLocale
            }}</span>
          </button>
          <!-- Dropdown -->
          <div
            class="absolute right-0 mt-2 w-32 rounded-xl z-50 transition-all transform origin-top-right bg-card border border-border shadow-lg overflow-hidden"
            :class="[
              activeMenu === 'lang'
                ? 'opacity-100 visible scale-100'
                : 'opacity-0 invisible scale-95 pointer-events-none',
            ]"
            @click.stop
          >
            <button
              v-for="opt in localeOptions"
              :key="opt.value"
              @click="currentLocale = opt.value"
              class="w-full text-left px-4 py-2 text-sm transition-colors cursor-pointer"
              :class="{
                'bg-accent text-accent-foreground font-bold':
                  currentLocale === opt.value,
                'text-foreground hover:bg-muted': currentLocale !== opt.value,
              }"
            >
              {{ opt.label }}
            </button>
          </div>
        </div>

        <!-- Mode Switcher -->
        <button
          @click="toggleThemeMode"
          class="p-1.5 rounded-full transition-all hover:rotate-12 cursor-pointer bg-card border border-border text-foreground shadow-sm hover:bg-accent hover:text-accent-foreground"
          :title="isDark ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
        >
          <IconLucideMoon v-if="isDark" class="w-4 h-4" />
          <IconLucideSun v-else class="w-4 h-4" />
        </button>
      </div>
    </header>

    <main
      class="relative z-10 w-full h-full flex flex-col items-center justify-center min-h-screen p-4"
    >
      <slot />
    </main>
  </div>
</template>
