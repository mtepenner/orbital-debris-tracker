export const RISK_COLOR_VERTEX = `
varying vec3 vWorldNormal;
void main() {
  vWorldNormal = normalize(normalMatrix * normal);
  gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
}`;

export const RISK_COLOR_FRAGMENT = `
varying vec3 vWorldNormal;
void main() {
  float intensity = pow(0.65 - dot(vWorldNormal, vec3(0.0, 0.0, 1.0)), 2.4);
  vec3 glow = mix(vec3(0.03, 0.08, 0.20), vec3(0.15, 0.50, 1.0), intensity);
  gl_FragColor = vec4(glow, 1.0);
}`;
