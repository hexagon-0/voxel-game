#version 410 core

uniform mat4 uModel;
uniform mat4 uView;
uniform mat4 uProjection;

layout (location = 0) in uvec3 aPosition;
layout (location = 1) in uint aFaceIndex;

out vec3 vPos;
flat out uint vFaceIndex;

void main() {
	vec3 pos = vec3(aPosition);
	vPos = pos;
	vFaceIndex = aFaceIndex;
	gl_Position = uProjection * uView * uModel * vec4(pos, 1.0);
}
