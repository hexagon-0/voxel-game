package render

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader = uint32

func NewShader(source string, shaderType Shader) (Shader, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile shader: %v", log)
	}

	return shader, nil
}

type ShaderProgram uint32

func NewShaderProgram(vs, fs Shader) (ShaderProgram, error) {
	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return ShaderProgram(0), fmt.Errorf("failed to link program: %v", log)
	}

	return ShaderProgram(program), nil
}

func (self ShaderProgram) UseProgram() {
	gl.UseProgram(uint32(self))
}

func (self ShaderProgram) SetUniformMatrix4fv(name string, value mgl32.Mat4) error {
	location := gl.GetUniformLocation(uint32(self), gl.Str(name + "\x00"))
	if location == -1 {
		return fmt.Errorf("glGetUniformLocation returned -1 for name `%s`. This could mean there is no uniform with this name.", name)
	}

	self.UseProgram()
	gl.UniformMatrix4fv(location, 1, false, &value[0])

	return nil
}

func (self ShaderProgram) SetUniform3fv(name string, value mgl32.Vec3) error {
	location := gl.GetUniformLocation(uint32(self), gl.Str(name + "\x00"))
	if location == -1 {
		return fmt.Errorf("glGetUniformLocation returned -1 for name `%s`. This could mean there is no uniform with this name.", name)
	}

	self.UseProgram()
	gl.Uniform3fv(location, 1, &value[0])

	return nil
}

