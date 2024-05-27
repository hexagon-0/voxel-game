package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type MeshInstance struct {
	mesh      *Mesh
	transform mgl32.Mat4
}

func NewMeshInstance(mesh *Mesh, transform mgl32.Mat4) MeshInstance {
	return MeshInstance{mesh, transform}
}

func (self *MeshInstance) Render(projection, view mgl32.Mat4) {
	self.mesh.shader.UseProgram()
	self.mesh.shader.SetUniformMatrix4fv("model", self.transform)
	self.mesh.shader.SetUniformMatrix4fv("view", view)
	self.mesh.shader.SetUniformMatrix4fv("projection", projection)
	gl.BindVertexArray(self.mesh.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(self.mesh.vertices)/VERTEX_ELEMENTS))
}

