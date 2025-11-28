<!-- Transform from reactbits -->

<script setup lang="ts">
import { type StyleValue } from "vue";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import { SplitText } from "gsap/SplitText";

gsap.registerPlugin(ScrollTrigger, SplitText);

export interface ShuffleProps {
  text: string;
  class?: string;
  style?: StyleValue;
  shuffleDirection?: "left" | "right";
  duration?: number;
  maxDelay?: number;
  ease?: string | ((t: number) => number);
  threshold?: number;
  rootMargin?: string;
  tag?: string;
  textAlign?: "left" | "right" | "center" | "justify" | "initial" | "inherit";
  shuffleTimes?: number;
  animationMode?: "random" | "evenodd";
  loop?: boolean;
  loopDelay?: number;
  stagger?: number;
  scrambleCharset?: string;
  colorFrom?: string;
  colorTo?: string;
  triggerOnce?: boolean;
  respectReducedMotion?: boolean;
  triggerOnHover?: boolean;
}

const props = withDefaults(defineProps<ShuffleProps>(), {
  class: "",
  style: () => ({}),
  shuffleDirection: "right",
  duration: 0.35,
  maxDelay: 0,
  ease: "power3.out",
  threshold: 0.1,
  rootMargin: "-100px",
  tag: "p",
  textAlign: "center",
  shuffleTimes: 1,
  animationMode: "evenodd",
  loop: false,
  loopDelay: 0,
  stagger: 0.03,
  scrambleCharset: "",
  colorFrom: undefined,
  colorTo: undefined,
  triggerOnce: true,
  respectReducedMotion: true,
  triggerOnHover: true,
});

const emit = defineEmits<{
  (e: "shuffleComplete"): void;
}>();

const rootRef = ref<HTMLElement | null>(null);
const fontsLoaded = ref(false);
const isReady = ref(false);

const splitInstance = ref<SplitText | null>(null);
const wrappersRef = ref<HTMLElement[]>([]);
const tlRef = ref<gsap.core.Timeline | null>(null);
const playingRef = ref(false);
const hoverHandlerRef = ref<((e: Event) => void) | null>(null);
const ctx = ref<gsap.Context | null>(null);
const scrollTriggerInstance = ref<ScrollTrigger | null>(null);

onMounted(() => {
  if ("fonts" in document) {
    if (document.fonts.status === "loaded") {
      fontsLoaded.value = true;
    } else {
      document.fonts.ready.then(() => {
        fontsLoaded.value = true;
      });
    }
  } else {
    fontsLoaded.value = true;
  }
});

const computedClasses = computed(() => {
  const baseTw =
    "inline-block whitespace-normal break-words will-change-transform uppercase text-2xl leading-none";
  return `${baseTw} ${isReady.value ? "visible" : "invisible"} ${props.class}`.trim();
});

const computedStyles = computed((): StyleValue => {
  const userHasFont = props.class && /font[-[]/i.test(props.class);
  const fallbackFont = userHasFont
    ? {}
    : { fontFamily: `'Press Start 2P', sans-serif` };

  return {
    textAlign: props.textAlign,
    ...fallbackFont,
    ...(props.style as object),
  };
});

const removeHover = () => {
  if (hoverHandlerRef.value && rootRef.value) {
    rootRef.value.removeEventListener("mouseenter", hoverHandlerRef.value);
    hoverHandlerRef.value = null;
  }
};

const teardown = () => {
  if (tlRef.value) {
    tlRef.value.kill();
    tlRef.value = null;
  }

  if (wrappersRef.value.length) {
    wrappersRef.value.forEach((wrap) => {
      const inner = wrap.firstElementChild as HTMLElement | null;
      const orig = inner?.querySelector(
        '[data-orig="1"]',
      ) as HTMLElement | null;
      if (orig && wrap.parentNode) {
        wrap.parentNode.replaceChild(orig, wrap);
      }
    });
    wrappersRef.value = [];
  }

  try {
    splitInstance.value?.revert();
  } catch (e) {}

  splitInstance.value = null;
  playingRef.value = false;
};

const randomizeScrambles = () => {
  const charset = props.scrambleCharset;
  if (!charset) return;

  wrappersRef.value.forEach((w) => {
    const strip = w.firstElementChild as HTMLElement;
    if (!strip) return;
    const kids = Array.from(strip.children) as HTMLElement[];
    for (let i = 1; i < kids.length - 1; i++) {
      kids[i]!.textContent = charset.charAt(
        Math.floor(Math.random() * charset.length),
      );
    }
  });
};

const cleanupToStill = () => {
  wrappersRef.value.forEach((w) => {
    const strip = w.firstElementChild as HTMLElement;
    if (!strip) return;
    const real = strip.querySelector('[data-orig="1"]') as HTMLElement | null;
    if (!real) return;
    strip.replaceChildren(real);
    strip.style.transform = "none";
    strip.style.willChange = "auto";
  });
};

const build = () => {
  teardown();
  const el = rootRef.value;
  if (!el) return;

  const computedFont = getComputedStyle(el).fontFamily;

  splitInstance.value = new SplitText(el, {
    type: "chars",
    charsClass: "shuffle-char",
    wordsClass: "shuffle-word",
    linesClass: "shuffle-line",
  } as any);

  const chars = (splitInstance.value.chars || []) as HTMLElement[];
  wrappersRef.value = [];

  const rolls = Math.max(1, Math.floor(props.shuffleTimes));
  const rand = (set: string) =>
    set.charAt(Math.floor(Math.random() * set.length)) || "";
  const charset = props.scrambleCharset;

  chars.forEach((ch) => {
    const parent = ch.parentElement;
    if (!parent) return;

    const w = ch.getBoundingClientRect().width;
    if (!w) return;

    const wrap = document.createElement("span");
    wrap.className = "inline-block overflow-hidden align-baseline text-left";
    Object.assign(wrap.style, { width: w + "px" });

    const inner = document.createElement("span");
    inner.className =
      "inline-block whitespace-nowrap will-change-transform origin-left transform-gpu";

    parent.insertBefore(wrap, ch);
    wrap.appendChild(inner);

    const firstOrig = ch.cloneNode(true) as HTMLElement;
    firstOrig.className = "inline-block text-left";
    Object.assign(firstOrig.style, {
      width: w + "px",
      fontFamily: computedFont,
    });

    ch.setAttribute("data-orig", "1");
    ch.className = "inline-block text-left";
    Object.assign(ch.style, { width: w + "px", fontFamily: computedFont });

    inner.appendChild(firstOrig);

    for (let k = 0; k < rolls; k++) {
      const c = ch.cloneNode(true) as HTMLElement;
      if (charset) {
        c.textContent = rand(charset);
      }
      c.className = "inline-block text-left";
      Object.assign(c.style, { width: w + "px", fontFamily: computedFont });
      inner.appendChild(c);
    }

    inner.appendChild(ch);

    const steps = rolls + 1;
    let startX = 0;
    let finalX = -steps * w;

    if (props.shuffleDirection === "right") {
      const firstCopy = inner.firstElementChild as HTMLElement | null;
      const real = inner.lastElementChild as HTMLElement | null;
      if (real) inner.insertBefore(real, inner.firstChild);
      if (firstCopy) inner.appendChild(firstCopy);
      startX = -steps * w;
      finalX = 0;
    }

    gsap.set(inner, { x: startX, force3D: true });
    if (props.colorFrom) {
      (inner.style as any).color = props.colorFrom;
    }

    inner.setAttribute("data-final-x", String(finalX));
    inner.setAttribute("data-start-x", String(startX));

    wrappersRef.value.push(wrap);
  });
};

const play = () => {
  const strips = wrappersRef.value.map(
    (w) => w.firstElementChild as HTMLElement,
  );
  if (!strips.length) return;

  playingRef.value = true;

  const tl = gsap.timeline({
    smoothChildTiming: true,
    repeat: props.loop ? -1 : 0,
    repeatDelay: props.loop ? props.loopDelay : 0,
    onRepeat: () => {
      if (props.scrambleCharset) randomizeScrambles();
      gsap.set(strips, {
        x: (_i, t: HTMLElement) =>
          parseFloat(t.getAttribute("data-start-x") || "0"),
      });
      emit("shuffleComplete");
    },
    onComplete: () => {
      playingRef.value = false;
      if (!props.loop) {
        cleanupToStill();
        if (props.colorTo) gsap.set(strips, { color: props.colorTo });
        emit("shuffleComplete");
        armHover();
      }
    },
  });

  const addTween = (targets: HTMLElement[], at: number | string) => {
    tl.to(
      targets,
      {
        x: (_i, t: HTMLElement) =>
          parseFloat(t.getAttribute("data-final-x") || "0"),
        duration: props.duration,
        ease: props.ease,
        force3D: true,
        stagger: props.animationMode === "evenodd" ? props.stagger : 0,
      },
      at,
    );
    if (props.colorFrom && props.colorTo) {
      tl.to(
        targets,
        { color: props.colorTo, duration: props.duration, ease: props.ease },
        at,
      );
    }
  };

  if (props.animationMode === "evenodd") {
    const odd = strips.filter((_, i) => i % 2 === 1);
    const even = strips.filter((_, i) => i % 2 === 0);
    const oddTotal =
      (props.duration || 0.35) +
      Math.max(0, odd.length - 1) * (props.stagger || 0.03);
    const evenStart = odd.length ? oddTotal * 0.7 : 0;

    if (odd.length) addTween(odd, 0);
    if (even.length) addTween(even, evenStart);
  } else {
    strips.forEach((strip) => {
      const d = Math.random() * (props.maxDelay || 0);
      tl.to(
        strip,
        {
          x: parseFloat(strip.getAttribute("data-final-x") || "0"),
          duration: props.duration,
          ease: props.ease,
          force3D: true,
        },
        d,
      );
      if (props.colorFrom && props.colorTo) {
        tl.fromTo(
          strip,
          { color: props.colorFrom },
          { color: props.colorTo, duration: props.duration, ease: props.ease },
          d,
        );
      }
    });
  }

  tlRef.value = tl;
};

const armHover = () => {
  if (!props.triggerOnHover || !rootRef.value) return;
  removeHover();

  const handler = () => {
    if (playingRef.value) return;
    build();
    if (props.scrambleCharset) randomizeScrambles();
    play();
  };

  hoverHandlerRef.value = handler;
  rootRef.value.addEventListener("mouseenter", handler);
};

const init = () => {
  if (!rootRef.value || !props.text || !fontsLoaded.value) return;

  if (
    props.respectReducedMotion &&
    window.matchMedia &&
    window.matchMedia("(prefers-reduced-motion: reduce)").matches
  ) {
    emit("shuffleComplete");
    isReady.value = true;
    return;
  }

  const startPct = (1 - (props.threshold || 0.1)) * 100;
  const mm = /^(-?\d+(?:\.\d+)?)(px|em|rem|%)?$/.exec(props.rootMargin || "");
  const mv = mm ? parseFloat(mm[1] as string) : 0;
  const mu = mm ? mm[2] || "px" : "px";
  const sign =
    mv === 0 ? "" : mv < 0 ? `-=${Math.abs(mv)}${mu}` : `+=${mv}${mu}`;
  const start = `top ${startPct}%${sign}`;

  const create = () => {
    build();
    if (props.scrambleCharset) randomizeScrambles();
    play();
    armHover();
    isReady.value = true;
  };

  scrollTriggerInstance.value = ScrollTrigger.create({
    trigger: rootRef.value,
    start,
    once: props.triggerOnce,
    onEnter: create,
  });
};

watch(
  [
    () => props.text,
    () => props.duration,
    () => props.shuffleDirection,
    () => props.shuffleTimes,
    () => props.animationMode,
    fontsLoaded,
  ],
  () => {
    if (scrollTriggerInstance.value) scrollTriggerInstance.value.kill();
    removeHover();
    teardown();
    isReady.value = false;

    nextTick(() => {
      if (ctx.value) ctx.value.revert();
      ctx.value = gsap.context(() => {
        init();
      }, rootRef.value as Element);
    });
  },
  { deep: true },
);

onMounted(() => {
  ctx.value = gsap.context(
    () => {
      init();
    },
    rootRef.value as Element | undefined,
  );
});

onBeforeUnmount(() => {
  if (scrollTriggerInstance.value) scrollTriggerInstance.value.kill();
  removeHover();
  teardown();
  if (ctx.value) ctx.value.revert();
});
</script>

<template>
  <component
    :is="tag"
    ref="rootRef"
    :class="computedClasses"
    :style="computedStyles"
  >
    {{ text }}
  </component>
</template>
