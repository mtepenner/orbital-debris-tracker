export type CatalogObject = {
  object_id: string;
  name: string;
  x_km: number;
  y_km: number;
  z_km: number;
  speed_km_s: number;
};

export type ConjunctionAlert = {
  primary_id: string;
  secondary_id: string;
  miss_distance_km: number;
  probability_collision: number;
  tca: string;
};
