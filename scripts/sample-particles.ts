import * as THREE from "three";
import { OBJLoader } from "three/examples/jsm/loaders/OBJLoader.js";
import { MeshSurfaceSampler } from "three/examples/jsm/math/MeshSurfaceSampler.js";
import * as BufferGeometryUtils from "three/examples/jsm/utils/BufferGeometryUtils.js";
import { writeFileSync, readFileSync } from "node:fs";
import { resolve } from "node:path";

const PARTICLE_COUNT = 35000;
const SCALE = 0.35;
const OBJ_PATH = resolve("src/assets/IronMan3D.obj");
const OUTPUT = resolve("src/assets/particle-positions.json");

async function main() {
	const objText = readFileSync(OBJ_PATH, "utf-8");
	const group = new OBJLoader().parse(objText);

	const geometries: THREE.BufferGeometry[] = [];
	group.traverse((child) => {
		if ((child as THREE.Mesh).isMesh) {
			const mesh = child as THREE.Mesh;
			const geometry = mesh.geometry.clone();
			mesh.updateMatrix();
			geometry.applyMatrix4(mesh.matrix);
			geometries.push(geometry);
		}
	});

	if (geometries.length === 0) {
		console.error("No geometries found");
		process.exit(1);
	}

	const merged = BufferGeometryUtils.mergeGeometries(geometries);
	const mesh = new THREE.Mesh(merged, new THREE.MeshBasicMaterial());

	merged.computeBoundingBox();
	if (merged.boundingBox) {
		merged.center();
	}

	const sampler = new MeshSurfaceSampler(mesh).build();
	const positions: number[] = [];
	const temp = new THREE.Vector3();

	for (let i = 0; i < PARTICLE_COUNT; i++) {
		sampler.sample(temp);
		positions.push(
			Number((temp.x * SCALE).toFixed(4)),
			Number((temp.y * SCALE).toFixed(4)),
			Number((temp.z * SCALE).toFixed(4)),
		);
	}

	writeFileSync(OUTPUT, JSON.stringify(positions));
	console.log(`Generated ${PARTICLE_COUNT} positions → ${OUTPUT}`);
}

main().catch((err) => {
	console.error(err);
	process.exit(1);
});
