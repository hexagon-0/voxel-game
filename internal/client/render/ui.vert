#version 410 core

uniform mat4 uModel;
// uniform mat4 uView;
uniform mat4 uProjection;

layout (location = 0) in vec2 aPosition;

void main() {
	gl_Position = uProjection * uModel * vec4(aPosition, 0.0, 1.0);
}

