<!-- Transform from reactbits -->

<script setup lang="ts">
import { gsap } from "gsap";
import { InertiaPlugin } from "gsap/InertiaPlugin";

if (typeof window !== "undefined") {
  gsap.registerPlugin(InertiaPlugin);
}

interface Dot {
  cx: number;
  cy: number;
  xOffset: number;
  yOffset: number;
  _inertiaApplied: boolean;
}

interface Props {
  dotSize?: number;
  gap?: number;
  baseColor?: string;
  activeColor?: string;
  proximity?: number;
  speedTrigger?: number;
  shockRadius?: number;
  shockStrength?: number;
  maxSpeed?: number;
  resistance?: number;
  returnDuration?: number;
  className?: string;
  style?: Record<string, any>;
}

const props = withDefaults(defineProps<Props>(), {
  dotSize: 16,
  gap: 32,
  baseColor: "#5227FF",
  activeColor: "#5227FF",
  proximity: 150,
  speedTrigger: 100,
  shockRadius: 250,
  shockStrength: 5,
  maxSpeed: 5000,
  resistance: 750,
  returnDuration: 1.5,
  className: "",
});

const throttle = (func: (...args: any[]) => void, limit: number) => {
  let lastCall = 0;
  return function (this: any, ...args: any[]) {
    const now = performance.now();
    if (now - lastCall >= limit) {
      lastCall = now;
      func.apply(this, args);
    }
  };
};

const wrapperRef = ref<HTMLDivElement | null>(null);
const canvasRef = ref<HTMLCanvasElement | null>(null);

let dots: Dot[] = [];

const pointer = {
  x: 0,
  y: 0,
  vx: 0,
  vy: 0,
  speed: 0,
  lastTime: 0,
  lastX: 0,
  lastY: 0,
};

let rafId: number;
let ro: ResizeObserver | null = null;
let circlePath: Path2D | null = null;

const colorInterpolator = computed(() =>
  gsap.utils.interpolate(props.baseColor, props.activeColor)
);

const initPath = () => {
  if (typeof window === "undefined" || !window.Path2D) return;
  const p = new Path2D();
  p.arc(0, 0, props.dotSize / 2, 0, Math.PI * 2);
  circlePath = p;
};

const buildGrid = () => {
  const wrap = wrapperRef.value;
  const canvas = canvasRef.value;
  if (!wrap || !canvas) return;

  const { width, height } = wrap.getBoundingClientRect();
  const dpr = window.devicePixelRatio || 1;

  canvas.width = width * dpr;
  canvas.height = height * dpr;
  canvas.style.width = `${width}px`;
  canvas.style.height = `${height}px`;

  const ctx = canvas.getContext("2d");
  if (ctx) ctx.scale(dpr, dpr);

  const cols = Math.floor((width + props.gap) / (props.dotSize + props.gap));
  const rows = Math.floor((height + props.gap) / (props.dotSize + props.gap));
  const cell = props.dotSize + props.gap;

  const gridW = cell * cols - props.gap;
  const gridH = cell * rows - props.gap;

  const extraX = width - gridW;
  const extraY = height - gridH;

  const startX = extraX / 2 + props.dotSize / 2;
  const startY = extraY / 2 + props.dotSize / 2;

  dots = []; // 重置数组
  for (let y = 0; y < rows; y++) {
    for (let x = 0; x < cols; x++) {
      const cx = startX + x * cell;
      const cy = startY + y * cell;
      dots.push({ cx, cy, xOffset: 0, yOffset: 0, _inertiaApplied: false });
    }
  }
};

const draw = () => {
  const canvas = canvasRef.value;
  if (!canvas || !circlePath) return;
  const ctx = canvas.getContext("2d");
  if (!ctx) return;

  ctx.clearRect(0, 0, canvas.width, canvas.height);

  const { x: px, y: py } = pointer;
  const proxSq = props.proximity * props.proximity;
  const interpolator = colorInterpolator.value;

  for (let i = 0; i < dots.length; i++) {
    const dot = dots[i]!;
    const ox = dot.cx + dot.xOffset;
    const oy = dot.cy + dot.yOffset;
    const dx = dot.cx - px;
    const dy = dot.cy - py;
    const dsq = dx * dx + dy * dy;

    let fillStyle = props.baseColor;

    if (dsq <= proxSq) {
      const dist = Math.sqrt(dsq);
      const t = 1 - dist / props.proximity;
      fillStyle = interpolator(t);
    }

    ctx.save();
    ctx.translate(ox, oy);
    ctx.fillStyle = fillStyle;
    ctx.fill(circlePath);
    ctx.restore();
  }

  rafId = requestAnimationFrame(draw);
};

const onMove = (e: MouseEvent) => {
  const now = performance.now();
  const dt = pointer.lastTime ? now - pointer.lastTime : 16;
  const dx = e.clientX - pointer.lastX;
  const dy = e.clientY - pointer.lastY;

  let vx = (dx / dt) * 1000;
  let vy = (dy / dt) * 1000;
  let speed = Math.hypot(vx, vy);

  if (speed > props.maxSpeed) {
    const scale = props.maxSpeed / speed;
    vx *= scale;
    vy *= scale;
    speed = props.maxSpeed;
  }

  pointer.lastTime = now;
  pointer.lastX = e.clientX;
  pointer.lastY = e.clientY;
  pointer.vx = vx;
  pointer.vy = vy;
  pointer.speed = speed;

  if (canvasRef.value) {
    const rect = canvasRef.value.getBoundingClientRect();
    pointer.x = e.clientX - rect.left;
    pointer.y = e.clientY - rect.top;
  }

  for (let i = 0; i < dots.length; i++) {
    const dot = dots[i]!;
    const dist = Math.hypot(dot.cx - pointer.x, dot.cy - pointer.y);

    if (
      speed > props.speedTrigger &&
      dist < props.proximity &&
      !dot._inertiaApplied
    ) {
      dot._inertiaApplied = true;
      gsap.killTweensOf(dot);

      const pushX = dot.cx - pointer.x + vx * 0.005;
      const pushY = dot.cy - pointer.y + vy * 0.005;

      gsap.to(dot, {
        inertia: {
          xOffset: pushX,
          yOffset: pushY,
          resistance: props.resistance,
        },
        onComplete: () => {
          gsap.to(dot, {
            xOffset: 0,
            yOffset: 0,
            duration: props.returnDuration,
            ease: "elastic.out(1,0.75)",
            onComplete: () => {
              dot._inertiaApplied = false;
            },
          });
        },
      });
    }
  }
};

const onClick = (e: MouseEvent) => {
  if (!canvasRef.value) return;
  const rect = canvasRef.value.getBoundingClientRect();
  const cx = e.clientX - rect.left;
  const cy = e.clientY - rect.top;

  for (let i = 0; i < dots.length; i++) {
    const dot = dots[i]!;
    const dist = Math.hypot(dot.cx - cx, dot.cy - cy);

    if (dist < props.shockRadius && !dot._inertiaApplied) {
      dot._inertiaApplied = true;
      gsap.killTweensOf(dot);

      const falloff = Math.max(0, 1 - dist / props.shockRadius);
      const pushX = (dot.cx - cx) * props.shockStrength * falloff;
      const pushY = (dot.cy - cy) * props.shockStrength * falloff;

      gsap.to(dot, {
        inertia: {
          xOffset: pushX,
          yOffset: pushY,
          resistance: props.resistance,
        },
        onComplete: () => {
          gsap.to(dot, {
            xOffset: 0,
            yOffset: 0,
            duration: props.returnDuration,
            ease: "elastic.out(1,0.75)",
            onComplete: () => {
              dot._inertiaApplied = false;
            },
          });
        },
      });
    }
  }
};

watch(
  () => props.dotSize,
  () => {
    initPath();
    buildGrid();
  }
);

watch(() => props.gap, buildGrid);

onMounted(() => {
  initPath();
  buildGrid();
  draw();

  const throttledMove = throttle(onMove, 50);

  const onTouchMove = (e: TouchEvent) => {
    if (e.touches.length > 0) {
      const touch = e.touches[0];
      if (!touch) return;
      // Create a fake mouse event from the touch
      const mouseEvent = {
        clientX: touch.clientX,
        clientY: touch.clientY,
      } as MouseEvent;
      onMove(mouseEvent);
    }
  };

  const throttledTouch = throttle(onTouchMove, 50);

  window.addEventListener("mousemove", throttledMove, { passive: true });
  window.addEventListener("touchmove", throttledTouch, { passive: true });
  window.addEventListener("click", onClick);

  if ("ResizeObserver" in window && wrapperRef.value) {
    ro = new ResizeObserver(buildGrid);
    ro.observe(wrapperRef.value);
  } else {
    window.addEventListener("resize", buildGrid);
  }

  onUnmounted(() => {
    cancelAnimationFrame(rafId);
    window.removeEventListener("mousemove", throttledMove);
    window.removeEventListener("touchmove", throttledTouch);
    window.removeEventListener("click", onClick);
    if (ro) ro.disconnect();
    else window.removeEventListener("resize", buildGrid);

    dots.forEach((dot) => gsap.killTweensOf(dot));
  });
});
</script>

<template>
  <section
    ref="wrapperRef"
    :class="`p-4 flex items-center justify-center h-full w-full relative ${props.className}`"
    :style="props.style"
  >
    <div className="w-full h-full relative">
      <canvas
        ref="canvasRef"
        className="absolute inset-0 w-full h-full pointer-events-none"
      />
    </div>
  </section>
</template>
