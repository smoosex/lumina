<script setup lang="ts">
import * as THREE from "three";
import ModelUrl from "@/assets/IronMan3D.obj?url";
import { oklchToHex } from "@/utils/color";

const containerRef = ref<HTMLElement | null>(null);

let scene: THREE.Scene;
let camera: THREE.PerspectiveCamera;
let renderer: THREE.WebGLRenderer;
let particles: THREE.Points;

let mouse: THREE.Vector2;
let animationId: number;

const themeStore = useThemeStore();

const getThemeColor = () => {
  if (typeof document === "undefined") return new THREE.Color("#5227ff");

  const style = getComputedStyle(document.documentElement);
  const colorVar = style.getPropertyValue("--primary").trim();

  if (colorVar.startsWith("oklch")) {
    return new THREE.Color(oklchToHex(colorVar));
  }

  return new THREE.Color(colorVar || "#5227ff");
};

const props = withDefaults(
  defineProps<{
    active?: boolean;
  }>(),
  {
    active: true,
  }
);

let originalPositions: Float32Array;
let velocities: Float32Array;
let dispersedPositions: Float32Array;
let modelPositions: Float32Array | null = null;

const init = () => {
  if (!containerRef.value) return;

  // Scene setup
  scene = new THREE.Scene();

  // Camera setup
  const width = containerRef.value.clientWidth;
  const height = containerRef.value.clientHeight;
  camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 1000);
  camera.position.z = 40;

  // Renderer setup
  renderer = new THREE.WebGLRenderer({ alpha: true, antialias: true });
  renderer.setSize(width, height);
  renderer.setPixelRatio(window.devicePixelRatio);
  containerRef.value.appendChild(renderer.domElement);

  // Initialize with random particles immediately
  const particleCount = 35000;
  const positions = new Float32Array(particleCount * 3);
  originalPositions = new Float32Array(particleCount * 3);
  dispersedPositions = new Float32Array(particleCount * 3);
  velocities = new Float32Array(particleCount * 3);

  const geometry = new THREE.BufferGeometry();

  // Create a random cloud initially
  for (let i = 0; i < particleCount; i++) {
    const r = 10 + Math.random() * 20;
    const theta = Math.random() * Math.PI * 2;
    const phi = Math.acos(2 * Math.random() - 1);

    const x = r * Math.sin(phi) * Math.cos(theta);
    const y = r * Math.sin(phi) * Math.sin(theta);
    const z = r * Math.cos(phi);

    positions[i * 3] = x;
    positions[i * 3 + 1] = y;
    positions[i * 3 + 2] = z;

    // Store as dispersed positions
    dispersedPositions[i * 3] = x;
    dispersedPositions[i * 3 + 1] = y;
    dispersedPositions[i * 3 + 2] = z;

    // Initial target is dispersed
    originalPositions[i * 3] = x;
    originalPositions[i * 3 + 1] = y;
    originalPositions[i * 3 + 2] = z;
  }

  geometry.setAttribute("position", new THREE.BufferAttribute(positions, 3));

  // Material
  const material = new THREE.PointsMaterial({
    color: getThemeColor(),
    size: 0.1,
    transparent: true,
    opacity: 0.8,
    sizeAttenuation: true,
  });

  particles = new THREE.Points(geometry, material);
  scene.add(particles);

  // Use Web Worker for generation
  const worker = new Worker(
    new URL("../workers/particle.worker.ts", import.meta.url),
    { type: "module" }
  );

  worker.onmessage = (e) => {
    if (e.data.error) {
      console.error("Worker error:", e.data.error);
      return;
    }

    const newPositions = e.data.positions;

    // Store the model positions
    if (newPositions.length === originalPositions.length) {
      modelPositions = newPositions;

      // If currently active, switch target to model immediately
      if (props.active) {
        originalPositions.set(modelPositions!);
      }
    } else {
      console.warn(
        "Particle count mismatch",
        newPositions.length,
        originalPositions.length
      );
    }

    worker.terminate();
  };

  worker.postMessage({
    modelUrl: ModelUrl,
    particleCount: particleCount,
    scale: 0.35,
  });

  // Watch for active prop changes to toggle dispersion
  watch(
    () => props.active,
    (isActive) => {
      if (!modelPositions) return;

      if (isActive) {
        // Form the model
        originalPositions.set(modelPositions);
      } else {
        // Disperse
        originalPositions.set(dispersedPositions);
      }
    }
  );

  // Watch for theme changes
  watch(
    () => [themeStore.mode, themeStore.themeName],
    async () => {
      await nextTick();
      // Small delay to ensure CSS variables are updated
      setTimeout(() => {
        if (particles && particles.material) {
          const newColor = getThemeColor();
          (particles.material as THREE.PointsMaterial).color = newColor;
        }
      }, 100);
    }
  );

  // Raycaster & Mouse

  mouse = new THREE.Vector2(-1000, -1000); // Start off-screen

  // Event Listeners
  window.addEventListener("resize", onWindowResize);
  containerRef.value.addEventListener("mousemove", onMouseMove);
  containerRef.value.addEventListener("mouseleave", onMouseLeave);

  animate();
};

const onWindowResize = () => {
  if (!containerRef.value) return;
  const width = containerRef.value.clientWidth;
  const height = containerRef.value.clientHeight;
  camera.aspect = width / height;
  camera.updateProjectionMatrix();
  renderer.setSize(width, height);
};

const onMouseMove = (event: MouseEvent) => {
  if (!containerRef.value) return;
  const rect = containerRef.value.getBoundingClientRect();
  mouse.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
  mouse.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;
};

const onMouseLeave = () => {
  mouse.x = -1000;
  mouse.y = -1000;
};

const animate = () => {
  animationId = requestAnimationFrame(animate);

  if (!particles || !originalPositions || !velocities) {
    renderer.render(scene, camera);
    return;
  }

  // Rotate the whole model slowly
  particles.rotation.y += 0.002;

  // Simplified interaction:
  // Calculate mouse position in world space at z=0 (approximate)
  const vector = new THREE.Vector3(mouse.x, mouse.y, 0.5);
  vector.unproject(camera);
  const dir = vector.sub(camera.position).normalize();
  const distance = -camera.position.z / dir.z;
  const mousePos = camera.position.clone().add(dir.multiplyScalar(distance));

  const positionsAttribute = particles.geometry.attributes.position;
  if (!positionsAttribute) return;
  const positions = positionsAttribute.array as Float32Array;

  // Physics params
  const repulsionRadius = 8;
  const repulsionStrength = 2;
  const returnSpeed = 0.05;
  const friction = 0.9;
  const boilingStrength = 0.0025;

  // Transform mouse to local space of the particles object.
  const localMouse = mousePos
    .clone()
    .applyMatrix4(particles.matrixWorld.invert());

  const count = positions.length / 3;

  // Cast to any to avoid strict TS undefined checks in tight loop
  const p = positions as any;
  const v = velocities as any;
  const op = originalPositions as any;

  for (let i = 0; i < count; i++) {
    const i3 = i * 3;
    const px = p[i3];
    const py = p[i3 + 1];
    const pz = p[i3 + 2];

    const dx = px - localMouse.x;
    const dy = py - localMouse.y;
    const dz = pz - localMouse.z;
    const distSq = dx * dx + dy * dy + dz * dz;

    if (distSq < repulsionRadius * repulsionRadius) {
      const dist = Math.sqrt(distSq);
      const force = (repulsionRadius - dist) / repulsionRadius;

      const angleX = dx / dist;
      const angleY = dy / dist;
      const angleZ = dz / dist;

      v[i3] += angleX * force * repulsionStrength;
      v[i3 + 1] += angleY * force * repulsionStrength;
      v[i3 + 2] += angleZ * force * repulsionStrength;
    }

    // Boiling effect (idle animation)
    v[i3] += (Math.random() - 0.5) * boilingStrength;
    v[i3 + 1] += (Math.random() - 0.5) * boilingStrength;
    v[i3 + 2] += (Math.random() - 0.5) * boilingStrength;

    // Return to original position
    v[i3] += (op[i3] - px) * returnSpeed;
    v[i3 + 1] += (op[i3 + 1] - py) * returnSpeed;
    v[i3 + 2] += (op[i3 + 2] - pz) * returnSpeed;

    // Friction
    v[i3] *= friction;
    v[i3 + 1] *= friction;
    v[i3 + 2] *= friction;

    // Update position
    p[i3] = px + v[i3];
    p[i3 + 1] = py + v[i3 + 1];
    p[i3 + 2] = pz + v[i3 + 2];
  }

  positionsAttribute.needsUpdate = true;
  renderer.render(scene, camera);
};

onMounted(() => {
  requestAnimationFrame(() => {
    init();
  });
});

onUnmounted(() => {
  cancelAnimationFrame(animationId);
  window.removeEventListener("resize", onWindowResize);
  if (containerRef.value) {
    containerRef.value.removeEventListener("mousemove", onMouseMove);
    containerRef.value.removeEventListener("mouseleave", onMouseLeave);
    if (renderer) {
      containerRef.value.removeChild(renderer.domElement);
      renderer.dispose();
    }
  }
});
</script>

<template>
  <div ref="containerRef" class="w-full h-full min-h-[400px]"></div>
</template>
