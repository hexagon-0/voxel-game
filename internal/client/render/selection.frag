#version 410 core

in vec3 vPos;
flat in uint vFaceIndex;

out vec4 FragColor;

vec3 faceNormals[6] = vec3[](
	vec3( 1.0,  0.0,  0.0),
	vec3(-1.0,  0.0,  0.0),
	vec3( 0.0,  1.0,  0.0),
	vec3( 0.0, -1.0,  0.0),
	vec3( 0.0,  0.0,  1.0),
	vec3( 0.0,  0.0, -1.0)
);

float gridFactor (vec2 parameter, float width, float feather) {
  float w1 = width - feather * 0.5;
  vec2 d = fwidth(parameter);
  vec2 looped = 0.5 - abs(mod(parameter, 1.0) - 0.5);
  vec2 a2 = smoothstep(d * w1, d * (w1 + feather), looped);
  return min(a2.x, a2.y);
}

float gridFactor (vec2 parameter, float width) {
  vec2 d = fwidth(parameter);
  vec2 looped = 0.5 - abs(mod(parameter, 1.0) - 0.5);
  vec2 a2 = smoothstep(d * (width - 0.5), d * (width + 0.5), looped);
  return min(a2.x, a2.y);
}

void main() {
	vec3 color = vec3(0.8);
	uint i = vFaceIndex;
	vec3 normal = abs(faceNormals[i]);
	vec2 uv = vec2(dot(normal.zxy, vPos), dot(normal.yzx, vPos));
	float f = 1.0-gridFactor(uv, 1.0);
	if (f == 0.0) discard;
	color = vec3(0.6);
	FragColor = vec4(color, 1.0);
}
