#version 410 core

uniform mat4 uModel;
uniform mat4 uView;
uniform mat4 uProjection;

layout (location = 0) in vec3 aPosition;
layout (location = 1) in vec2 aTex;

out vec2 fragTexCoord;

void main() {
	gl_Position = uProjection * uView * uModel * vec4(aPosition, 1.0);
	fragTexCoord = aTex;
}

