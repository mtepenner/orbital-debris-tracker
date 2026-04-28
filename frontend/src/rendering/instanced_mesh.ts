import * as THREE from "three";

import type { CatalogObject } from "../types";

export function createObjectCloud(capacity: number) {
  const geometry = new THREE.SphereGeometry(0.11, 6, 6);
  const material = new THREE.MeshBasicMaterial({ vertexColors: true });
  const mesh = new THREE.InstancedMesh(geometry, material, capacity);
  mesh.instanceMatrix.setUsage(THREE.DynamicDrawUsage);
  mesh.count = 0;
  const helper = new THREE.Object3D();

  const update = (objects: CatalogObject[], flagged: Set<string>) => {
    mesh.count = Math.min(objects.length, capacity);
    for (let index = 0; index < mesh.count; index += 1) {
      const object = objects[index];
      helper.position.set(object.x_km / 200, object.y_km / 200, object.z_km / 200);
      const scale = THREE.MathUtils.mapLinear(object.speed_km_s, 6.8, 8.0, 0.8, 1.3);
      helper.scale.setScalar(scale);
      helper.updateMatrix();
      mesh.setMatrixAt(index, helper.matrix);
      mesh.setColorAt(index, flagged.has(object.object_id) ? new THREE.Color("#f97316") : new THREE.Color("#34d399"));
    }
    mesh.instanceMatrix.needsUpdate = true;
    if (mesh.instanceColor) {
      mesh.instanceColor.needsUpdate = true;
    }
  };

  return { mesh, update };
}
