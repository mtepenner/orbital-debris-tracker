import { useEffect, useRef } from "react";
import * as THREE from "three";

import { createObjectCloud } from "../rendering/instanced_mesh";
import { RISK_COLOR_FRAGMENT, RISK_COLOR_VERTEX } from "../rendering/shaders/riskColor";
import type { CatalogObject, ConjunctionAlert } from "../types";

type Props = {
  objects: CatalogObject[];
  alerts: ConjunctionAlert[];
};

export function GlobeVisualizer({ objects, alerts }: Props) {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const cloudRef = useRef<ReturnType<typeof createObjectCloud> | null>(null);
  const sceneStateRef = useRef<{
    renderer: THREE.WebGLRenderer;
    camera: THREE.PerspectiveCamera;
    scene: THREE.Scene;
    frameId: number;
    resize: () => void;
  } | null>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) {
      return;
    }

    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(50, 1, 0.1, 400);
    camera.position.set(0, 0, 38);

    const renderer = new THREE.WebGLRenderer({ canvas, antialias: true, powerPreference: "high-performance" });
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

    const ambient = new THREE.AmbientLight("#dbeafe", 0.8);
    const key = new THREE.DirectionalLight("#ffffff", 1.2);
    key.position.set(12, 18, 16);
    scene.add(ambient, key);

    const earth = new THREE.Mesh(
      new THREE.SphereGeometry(7, 48, 48),
      new THREE.ShaderMaterial({ vertexShader: RISK_COLOR_VERTEX, fragmentShader: RISK_COLOR_FRAGMENT })
    );
    scene.add(earth);

    const cloud = createObjectCloud(40000);
    cloudRef.current = cloud;
    scene.add(cloud.mesh);

    const resize = () => {
      const width = canvas.clientWidth;
      const height = canvas.clientHeight;
      if (width === 0 || height === 0) {
        return;
      }
      camera.aspect = width / height;
      camera.updateProjectionMatrix();
      renderer.setSize(width, height, false);
    };

    const animate = () => {
      earth.rotation.y += 0.002;
      cloud.mesh.rotation.y += 0.0008;
      renderer.render(scene, camera);
      sceneStateRef.current!.frameId = window.requestAnimationFrame(animate);
    };

    sceneStateRef.current = { renderer, camera, scene, frameId: window.requestAnimationFrame(animate), resize };
    window.addEventListener("resize", resize);
    resize();

    return () => {
      const state = sceneStateRef.current;
      if (!state) {
        return;
      }
      window.cancelAnimationFrame(state.frameId);
      window.removeEventListener("resize", state.resize);
      state.renderer.dispose();
      sceneStateRef.current = null;
    };
  }, []);

  useEffect(() => {
    const cloud = cloudRef.current;
    if (!cloud) {
      return;
    }
    const flagged = new Set<string>();
    alerts.forEach((alert) => {
      flagged.add(alert.primary_id);
      flagged.add(alert.secondary_id);
    });
    cloud.update(objects, flagged);
  }, [alerts, objects]);

  return <canvas ref={canvasRef} className="globe-canvas" />;
}
