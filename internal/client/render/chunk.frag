#version 410 core

uniform sampler2D tTexture;

in vec2 fragTexCoord;

out vec4 FragColor;

void main() {
	FragColor = texture(tTexture, fragTexCoord);
}

