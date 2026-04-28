import type { CatalogObject } from "../types";

type Props = {
  object: CatalogObject | null;
};

export function ObjectInspector({ object }: Props) {
  return (
    <section className="panel inspector-panel">
      <div className="panel-header">
        <p className="eyebrow">Object Inspector</p>
        <h2>{object?.name ?? "Select an object"}</h2>
      </div>
      {object ? (
        <dl className="inspector-grid">
          <div><dt>NORAD</dt><dd>{object.object_id}</dd></div>
          <div><dt>X</dt><dd>{object.x_km.toFixed(1)} km</dd></div>
          <div><dt>Y</dt><dd>{object.y_km.toFixed(1)} km</dd></div>
          <div><dt>Z</dt><dd>{object.z_km.toFixed(1)} km</dd></div>
          <div><dt>Speed</dt><dd>{object.speed_km_s.toFixed(2)} km/s</dd></div>
        </dl>
      ) : (
        <p className="muted-text">Choose a conjunction alert to inspect the object state.</p>
      )}
    </section>
  );
}
