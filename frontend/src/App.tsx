import { useEffect, useMemo, useState } from "react";

import { ConjunctionAlerts } from "./components/ConjunctionAlerts";
import { GlobeVisualizer } from "./components/GlobeVisualizer";
import { ObjectInspector } from "./components/ObjectInspector";
import type { CatalogObject, ConjunctionAlert } from "./types";

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? "http://127.0.0.1:8000";

export default function App() {
  const [objects, setObjects] = useState<CatalogObject[]>([]);
  const [alerts, setAlerts] = useState<ConjunctionAlert[]>([]);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [status, setStatus] = useState<"connecting" | "live" | "fallback">("connecting");

  useEffect(() => {
    let active = true;

    const load = async () => {
      try {
        const [catalogResponse, alertsResponse] = await Promise.all([
          fetch(`${API_BASE}/catalog`),
          fetch(`${API_BASE}/conjunctions`),
        ]);
        const [catalog, conjunctions] = await Promise.all([catalogResponse.json(), alertsResponse.json()]);
        if (!active) {
          return;
        }
        setObjects(catalog);
        setAlerts(conjunctions);
        setSelectedId((current) => current ?? catalog[0]?.object_id ?? null);
        setStatus("live");
      } catch {
        if (!active) {
          return;
        }
        const fallback = buildFallbackObjects();
        setObjects(fallback.objects);
        setAlerts(fallback.alerts);
        setSelectedId((current) => current ?? fallback.objects[0]?.object_id ?? null);
        setStatus("fallback");
      }
    };

    void load();
    const intervalId = window.setInterval(() => void load(), 3500);
    return () => {
      active = false;
      window.clearInterval(intervalId);
    };
  }, []);

  const selectedObject = useMemo(
    () => objects.find((object) => object.object_id === selectedId) ?? null,
    [objects, selectedId]
  );

  return (
    <main className="app-shell">
      <section className="hero-panel">
        <div>
          <p className="eyebrow">Orbital Command Center</p>
          <h1>Debris Cloud Surveillance</h1>
          <p className="hero-copy">
            Track catalog objects, surface close approaches, and inspect predicted orbital states
            from a single debris-monitoring console.
          </p>
        </div>
        <div className="status-grid">
          <article className="status-card"><span>Feed</span><strong>{status.toUpperCase()}</strong></article>
          <article className="status-card"><span>Catalog</span><strong>{objects.length}</strong></article>
          <article className="status-card"><span>Alerts</span><strong>{alerts.length}</strong></article>
        </div>
      </section>

      <section className="content-grid">
        <div className="visual-column">
          <section className="panel visual-panel">
            <div className="panel-header">
              <p className="eyebrow">3D Globe Visualizer</p>
              <h2>Orbital Geometry</h2>
            </div>
            <GlobeVisualizer objects={objects} alerts={alerts} />
          </section>
        </div>

        <div className="side-column">
          <ConjunctionAlerts alerts={alerts} onSelect={setSelectedId} />
          <ObjectInspector object={selectedObject} />
        </div>
      </section>
    </main>
  );
}

function buildFallbackObjects(): { objects: CatalogObject[]; alerts: ConjunctionAlert[] } {
  const objects = [
    { object_id: "25544", name: "ISS (ZARYA)", x_km: 6780, y_km: 1200, z_km: 820, speed_km_s: 7.67 },
    { object_id: "29716", name: "FENGYUN 1C DEB", x_km: -3220, y_km: 5300, z_km: 4210, speed_km_s: 7.13 },
    { object_id: "54237", name: "COSMOS 1408 DEB", x_km: -3600, y_km: 5180, z_km: 4180, speed_km_s: 7.38 },
  ];
  const alerts = [
    {
      primary_id: "29716",
      secondary_id: "54237",
      miss_distance_km: 28.4,
      probability_collision: 0.00057,
      tca: new Date().toISOString(),
    },
  ];
  return { objects, alerts };
}
