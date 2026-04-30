<script setup lang="ts">
import type { Project } from "@/config/site";

defineProps<{
	works: Project[];
	onNavigate: (link: string) => void;
}>();

const trackRef = ref<HTMLElement | null>(null);

let rafId: number | null = null;
let offset = 0;
const speed = 0.5;
let idleTimer: ReturnType<typeof setTimeout> | null = null;

const animate = () => {
	const track = trackRef.value;
	if (!track) return;

	offset -= speed;
	const halfWidth = track.scrollWidth / 2;
	if (offset <= -halfWidth) offset += halfWidth;
	if (offset >= halfWidth) offset -= halfWidth;

	track.style.transform = `translate3d(${offset}px, 0, 0)`;
	rafId = requestAnimationFrame(animate);
};

const stopAuto = () => {
	if (rafId !== null) {
		cancelAnimationFrame(rafId);
		rafId = null;
	}
	if (idleTimer !== null) {
		clearTimeout(idleTimer);
		idleTimer = null;
	}
};

const startAuto = () => {
	if (rafId === null) animate();
};

const scheduleResume = () => {
	if (idleTimer !== null) clearTimeout(idleTimer);
	idleTimer = setTimeout(() => {
		idleTimer = null;
		startAuto();
	}, 1000);
};

// Drag state
let isDragging = false;
let dragStartX = 0;
let dragStartOffset = 0;
let hasMoved = false;

const onPointerDown = (e: PointerEvent) => {
	stopAuto();
	isDragging = true;
	hasMoved = false;
	dragStartX = e.clientX;
	dragStartOffset = offset;
	(e.target as HTMLElement).setPointerCapture(e.pointerId);
};

const onPointerMove = (e: PointerEvent) => {
	if (!isDragging) return;
	const dx = e.clientX - dragStartX;
	if (Math.abs(dx) > 3) hasMoved = true;
	offset = dragStartOffset + dx;

	const track = trackRef.value;
	if (track) {
		const halfWidth = track.scrollWidth / 2;
		if (offset <= -halfWidth) offset += halfWidth;
		if (offset >= halfWidth) offset -= halfWidth;
	}

	if (trackRef.value) {
		trackRef.value.style.transform = `translate3d(${offset}px, 0, 0)`;
	}
};

const onPointerUp = (e: PointerEvent) => {
	if (!isDragging) return;
	isDragging = false;
	(e.target as HTMLElement).releasePointerCapture(e.pointerId);

	if (hasMoved) {
		e.preventDefault();
		e.stopPropagation();
	}

	scheduleResume();
};

const onPointerLeave = () => {
	if (isDragging) {
		isDragging = false;
	}
	scheduleResume();
};

const onWrapperClick = (e: MouseEvent) => {
	if (hasMoved) {
		e.stopPropagation();
		e.preventDefault();
	}
};

onMounted(() => {
	startAuto();
});

onUnmounted(() => {
	stopAuto();
});
</script>

<template>
	<section class="scroll-section py-20">
		<div class="text-center mb-16">
			<h2 class="text-4xl md:text-5xl font-bold text-foreground mb-4">
				{{ $t("home.works.title") }}
			</h2>
			<div class="w-24 h-1.5 bg-primary mx-auto rounded-full"></div>
		</div>

		<div
			class="overflow-hidden w-full relative touch-none select-none cursor-grab active:cursor-grabbing before:absolute before:left-0 before:top-0 before:h-full before:w-20 before:bg-gradient-to-r before:from-background before:to-transparent before:z-10 before:pointer-events-none after:absolute after:right-0 after:top-0 after:h-full after:w-20 after:bg-gradient-to-l after:from-background after:to-transparent after:z-10 after:pointer-events-none"
			@pointerenter="stopAuto"
			@pointerleave="onPointerLeave"
			@pointerdown="onPointerDown"
			@pointermove="onPointerMove"
			@pointerup="onPointerUp"
			@click.capture="onWrapperClick"
		>
			<div
				ref="trackRef"
				class="flex gap-8 will-change-transform px-6"
			>
				<template v-for="dup in 2" :key="dup">
					<div
						v-for="work in works"
						:key="`${dup}-${work.name}`"
						class="bg-card border border-border shadow-sm p-8 rounded-2xl group cursor-pointer relative overflow-hidden transition-all duration-300 hover:border-ring hover:shadow-lg shrink-0 w-[min(85vw,420px)]"
						@click="onNavigate(work.link)"
					>
						<div
							class="absolute inset-0 opacity-0 group-hover:opacity-10 transition-opacity duration-500 bg-linear-to-br"
							:class="work.color"
						></div>

						<div class="relative z-10 flex flex-col h-full">
							<div
								class="mb-6 p-4 rounded-xl bg-primary/10 w-fit group-hover:scale-110 transition-transform duration-300"
							>
								<component :is="work.icon" class="w-8 h-8 text-primary" />
							</div>

							<h3
								class="text-2xl font-bold mb-3 text-foreground group-hover:text-primary transition-colors"
							>
								{{ $t("home.works." + work.name + ".title") }}
							</h3>

							<p class="text-muted-foreground mb-8 grow leading-relaxed">
								{{ $t("home.works." + work.name + ".description") }}
							</p>

							<div class="flex items-center text-primary font-medium group/link">
								<span class="mr-2 group-hover/link:mr-3 transition-all">{{
									$t("home.works.explore")
								}}</span>
								<IconLucideExternalLink
									v-if="work.external"
									class="w-4 h-4"
								/>
								<IconLucideArrowRight v-else class="w-4 h-4" />
							</div>
						</div>
					</div>
				</template>
			</div>
		</div>
	</section>
</template>
