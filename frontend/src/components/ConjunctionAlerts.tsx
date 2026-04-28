import type { ConjunctionAlert } from "../types";

type Props = {
  alerts: ConjunctionAlert[];
  onSelect: (objectId: string) => void;
};

export function ConjunctionAlerts({ alerts, onSelect }: Props) {
  return (
    <section className="panel alert-panel">
      <div className="panel-header">
        <p className="eyebrow">Conjunction Feed</p>
        <h2>Upcoming Close Approaches</h2>
      </div>
      <div className="alert-list">
        {alerts.length === 0 ? <p className="muted-text">No elevated conjunction risks in the current prediction window.</p> : null}
        {alerts.map((alert) => (
          <button key={`${alert.primary_id}-${alert.secondary_id}`} type="button" className="alert-card" onClick={() => onSelect(alert.primary_id)}>
            <strong>{alert.primary_id} ↔ {alert.secondary_id}</strong>
            <span>Miss distance {alert.miss_distance_km.toFixed(1)} km</span>
            <span>Pc {(alert.probability_collision * 100).toFixed(3)}%</span>
            <span>TCA {new Date(alert.tca).toLocaleString()}</span>
          </button>
        ))}
      </div>
    </section>
  );
}
