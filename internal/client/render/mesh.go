package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	VERTEX_ELEMENTS = 6

	MESH_VS_SOURCE = `
	#version 410 core
	uniform mat4 model;
	uniform mat4 view;
	uniform mat4 projection;
	layout (location = 0) in vec3 vertexPosition;
	layout (location = 1) in vec3 vertexColor;
	out vec3 color;
	void main() {
		gl_Position = projection * view * model * vec4(vertexPosition, 1.0);
		color = vertexColor;
	}
	` + "\x00"

	MESH_FS_SOURCE = `
	#version 410 core
	in vec3 color;
	out vec4 fragColor;
	void main() {
		//fragColor = vec4(0.5, 0.7, 0.5, 1.0);
		fragColor = vec4(color, 1.0);
	}
	` + "\x00"
)

type Mesh struct {
	vertices []float32
	vbo uint32
	vao uint32
	shader ShaderProgram
}

func NewMesh(vertices []float32, shader ShaderProgram) Mesh {
	vbo := makeVbo(vertices)
	vao := makeVao(vbo)
	return Mesh{vertices, vbo, vao, shader}
}

func makeVbo(vertices []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	return vbo
}

func makeVao(vbo uint32) uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	// vertex position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, VERTEX_ELEMENTS*4, nil)

	// vertex color
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, VERTEX_ELEMENTS*4, 3*4)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return vao
}

