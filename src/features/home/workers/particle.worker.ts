import * as THREE from "three";
import { OBJLoader } from "three/examples/jsm/loaders/OBJLoader.js";
import { MeshSurfaceSampler } from "three/examples/jsm/math/MeshSurfaceSampler.js";
import * as BufferGeometryUtils from "three/examples/jsm/utils/BufferGeometryUtils.js";

self.onmessage = async (e) => {
	const { modelUrl, particleCount, scale } = e.data;

	try {
		const loader = new OBJLoader();

		loader.load(
			modelUrl,
			(group) => {
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
					self.postMessage({ error: "No geometries found" });
					return;
				}

				const mergedGeometry =
					BufferGeometryUtils.mergeGeometries(geometries);
				const mesh = new THREE.Mesh(
					mergedGeometry,
					new THREE.MeshBasicMaterial(),
				);

				mergedGeometry.computeBoundingBox();
				if (mergedGeometry.boundingBox) {
					mergedGeometry.center();
				}

				const sampler = new MeshSurfaceSampler(mesh).build();
				const positions = new Float32Array(particleCount * 3);
				const tempPosition = new THREE.Vector3();

				for (let i = 0; i < particleCount; i++) {
					sampler.sample(tempPosition);
					positions[i * 3] = tempPosition.x * scale;
					positions[i * 3 + 1] = tempPosition.y * scale;
					positions[i * 3 + 2] = tempPosition.z * scale;
				}

				// Transfer the buffer to main thread
				// @ts-ignore - Worker postMessage signature mismatch in some TS configs
				self.postMessage({ positions }, [positions.buffer]);
			},
			undefined,
			(error: unknown) => {
				const msg =
					error instanceof Error ? error.message : String(error);
				self.postMessage({ error: msg });
			},
		);
	} catch (err: any) {
		self.postMessage({ error: err.message });
	}
};
