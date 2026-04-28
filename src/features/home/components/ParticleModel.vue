<script setup lang="ts">
import * as THREE from "three";
import { oklchToHex } from "@/utils/color";
import modelData from "@/assets/particle-positions.json";

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
let modelPositions: Float32Array;

const init = () => {
	if (!containerRef.value) return;

	scene = new THREE.Scene();

	const width = containerRef.value.clientWidth;
	const height = containerRef.value.clientHeight;
	camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 1000);
	camera.position.z = 40;

	renderer = new THREE.WebGLRenderer({ alpha: true, antialias: true });
	renderer.setSize(width, height);
	renderer.setPixelRatio(window.devicePixelRatio);
	containerRef.value.appendChild(renderer.domElement);

	const particleCount = modelData.length / 3;
	const positions = new Float32Array(modelData);

	modelPositions = new Float32Array(modelData);
	originalPositions = new Float32Array(modelData);
	dispersedPositions = new Float32Array(particleCount * 3);
	velocities = new Float32Array(particleCount * 3);

	for (let i = 0; i < particleCount; i++) {
		const r = 10 + Math.random() * 20;
		const theta = Math.random() * Math.PI * 2;
		const phi = Math.acos(2 * Math.random() - 1);
		dispersedPositions[i * 3] = r * Math.sin(phi) * Math.cos(theta);
		dispersedPositions[i * 3 + 1] = r * Math.sin(phi) * Math.sin(theta);
		dispersedPositions[i * 3 + 2] = r * Math.cos(phi);
	}

	const geometry = new THREE.BufferGeometry();
	geometry.setAttribute("position", new THREE.BufferAttribute(positions, 3));

	const material = new THREE.PointsMaterial({
		color: getThemeColor(),
		size: 0.1,
		transparent: true,
		opacity: 0.8,
		sizeAttenuation: true,
	});

	particles = new THREE.Points(geometry, material);
	scene.add(particles);

	if (props.active) {
		originalPositions.set(modelPositions);
	} else {
		originalPositions.set(dispersedPositions);
	}

	watch(
		() => props.active,
		(isActive) => {
			if (isActive) {
				originalPositions.set(modelPositions);
			} else {
				originalPositions.set(dispersedPositions);
			}
		}
	);

	watch(
		() => [themeStore.mode, themeStore.themeName],
		async () => {
			await nextTick();
			setTimeout(() => {
				if (particles?.material) {
					(particles.material as THREE.PointsMaterial).color = getThemeColor();
				}
			}, 100);
		}
	);

	mouse = new THREE.Vector2(-1000, -1000);
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

	particles.rotation.y += 0.002;

	const vector = new THREE.Vector3(mouse.x, mouse.y, 0.5);
	vector.unproject(camera);
	const dir = vector.sub(camera.position).normalize();
	const distance = -camera.position.z / dir.z;
	const mousePos = camera.position.clone().add(dir.multiplyScalar(distance));

	const positionsAttribute = particles.geometry.attributes.position;
	if (!positionsAttribute) return;
	const p = positionsAttribute.array as Float32Array;
	const v = velocities as unknown as Float32Array;
	const op = originalPositions as unknown as Float32Array;

	const localMouse = mousePos
		.clone()
		.applyMatrix4(particles.matrixWorld.invert());

	const repulsionRadius = 8;
	const repulsionStrength = 2;
	const returnSpeed = 0.05;
	const friction = 0.9;
	const boilingStrength = 0.0025;
	const count = p.length / 3;

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

			v[i3] += (dx / dist) * force * repulsionStrength;
			v[i3 + 1] += (dy / dist) * force * repulsionStrength;
			v[i3 + 2] += (dz / dist) * force * repulsionStrength;
		}

		v[i3] += (Math.random() - 0.5) * boilingStrength;
		v[i3 + 1] += (Math.random() - 0.5) * boilingStrength;
		v[i3 + 2] += (Math.random() - 0.5) * boilingStrength;

		v[i3] += (op[i3] - px) * returnSpeed;
		v[i3 + 1] += (op[i3 + 1] - py) * returnSpeed;
		v[i3 + 2] += (op[i3 + 2] - pz) * returnSpeed;

		v[i3] *= friction;
		v[i3 + 1] *= friction;
		v[i3 + 2] *= friction;

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
