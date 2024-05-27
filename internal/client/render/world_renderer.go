package render

import (
	_ "embed"
	"fmt"
	"image"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hexagon-0/voxel-game/internal/common/world"
)

//go:embed "chunk.vert"
var ChunkVsSource string

//go:embed "chunk.frag"
var ChunkFsSource string

type BlockRepo map[world.BlockId]image.Rectangle

type ChunkMesh struct {
	ModelMatrix  mgl32.Mat4
	VboSize      int
	EboSize      int
	Vbo          uint32
	Ebo          uint32
	Vao          uint32
	ElementCount int32
}

// Calculate number of voxels in the worst case scenario
func checkerboard(w, h, d uint) uint {
	loW, loH, loD := (w >> 1), (h >> 1), (d >> 1)
	hiW, hiH, hiD := loW+(w&1), loH+(h&1), loD+(d&1)

	return hiH*(hiW*hiD+loW*loD) + loH*(loW*hiD+hiW*loD)
}

func NewChunkMesh(width, height, depth uint) ChunkMesh {
	// number of cubes in the worst case (checkerboard chunk)
	worstCase := checkerboard(width, height, depth)

	// number of components in each vertex attribute
	var aPositionSize int32 = 3
	var aTexCoordsSize int32 = 2
	totalComponents := aPositionSize + aTexCoordsSize

	// worstCase * 6 sides * 4 vertices * 5 components (xyzst) * 4 bytes
	vboSize := int(worstCase * 6 * 4 * 5 * 4)
	// worstCase * 6 sides * 2 triangles * 3 vertex indices * 4 bytes
	eboSize := int(worstCase * 6 * 2 * 3 * 4)

	var vbo, ebo, vao uint32

	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)
	gl.GenVertexArrays(1, &vao)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, vboSize, nil, gl.DYNAMIC_DRAW)

	gl.EnableVertexAttribArray(0) // position
	gl.VertexAttribPointer(0, aPositionSize, gl.FLOAT, false, totalComponents*4, nil)

	gl.EnableVertexAttribArray(1) // tex coords
	gl.VertexAttribPointerWithOffset(1, aTexCoordsSize, gl.FLOAT, false, totalComponents*4, 3*4)

	// gl.EnableVertexAttribArray(2) // normal
	// gl.VertexAttribPointerWithOffset(2, 3, gl.FLOAT, false, 8, 5)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, eboSize, nil, gl.DYNAMIC_DRAW)
	gl.BindVertexArray(0)

	var elementCount int32 = 0 // int32(len(indices))

	return ChunkMesh{mgl32.Ident4(), vboSize, eboSize, vbo, ebo, vao, elementCount}
}

type WorldRenderer struct {
	Shader      ShaderProgram
	ChunkMeshes []ChunkMesh
}

func (self *WorldRenderer) CompileShaders() error {
	vs, err := NewShader(ChunkVsSource, gl.VERTEX_SHADER)
	if err != nil {
		return err
	}
	defer gl.DeleteShader(vs)

	fs, err := NewShader(ChunkFsSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}
	defer gl.DeleteShader(fs)

	self.Shader, err = NewShaderProgram(vs, fs)
	if err != nil {
		return err
	}

	return nil
}

func BuildChunkMesh(m *ChunkMesh, x, y, z int, wrld *world.World, blockRepo *BlockRepo) error {
	chunk := wrld.Chunks[[3]int{x, y, z}]
	if chunk == nil {
		return fmt.Errorf("Chunk not loaded: %d %d %d", x, y, z)
	}
	// ox, oy, oz := float32(x<<world.CHUNK_SHIFT), float32(y<<world.CHUNK_SHIFT), float32(z<<world.CHUNK_SHIFT)
	ox, oy, oz := float32(0.0), float32(0.0), float32(0.0)
	w, h, d := int(chunk.Width), int(chunk.Height), int(chunk.Depth)

	dir := [3][2][3]float32{} // perpendicular vectors for each axis (for vertex building)
	for i := 0; i < 3; i++ {
		dir[i][0][(i+1)%3] = 1.0
		dir[i][1][(i+2)%3] = 1.0
	}

	vertices := make([]float32, 0, m.VboSize)
	indices := make([]uint32, 0, m.EboSize)

	cb := [3]bool{} // for each axis, is the current block within chunk bounds
	nb := [3]bool{} // for each axis, is the next block on that axis within chunk bounds

	cb[2] = false
	nb[2] = true
	for k := -1; k < d; k++ {
		nb[2] = k < d-1

		cb[1] = false
		nb[1] = true
		for j := -1; j < h; j++ {
			nb[1] = j < h-1

			cb[0] = false
			nb[0] = true
			for i := -1; i < w; i++ {
				nb[0] = i < w-1

				var c world.BlockId // current block
				if cb[0] && cb[1] && cb[2] {
					c = chunk.BlockAt(uint(i), uint(j), uint(k))
				}

				n := [3]world.BlockId{} // neighbours for each axis
				if nb[0] && cb[1] && cb[2] {
					n[0] = chunk.BlockAt(uint(i)+1, uint(j), uint(k))
				}
				if cb[0] && nb[1] && cb[2] {
					n[1] = chunk.BlockAt(uint(i), uint(j)+1, uint(k))
				}
				if cb[0] && cb[1] && nb[2] {
					n[2] = chunk.BlockAt(uint(i), uint(j), uint(k)+1)
				}

				for di := 0; di < 3; di++ {
					if (c != 0) != (n[di] != 0) {
						var s int // winding order
						if c == 0 {
							s = 1
						}

						t := [3]float32{float32(i), float32(j), float32(k)}
						u := dir[di][s]
						v := dir[di][s^1]
						t[di]++

						ui := uint32(len(vertices)/5)
						vertices = append(vertices,
							ox+t[0], oy+t[1], oz+t[2], 0.0, 0.0,
							ox+t[0]+u[0], oy+t[1]+u[1], oz+t[2]+u[2], 1.0, 0.0,
							ox+t[0]+u[0]+v[0], oy+t[1]+u[1]+v[1], oz+t[2]+u[2]+v[2], 1.0, 1.0,
							ox+t[0]+v[0], oy+t[1]+v[1], oz+t[2]+v[2], 0.0, 1.0,
						)

						indices = append(indices,
							ui+0, ui+1, ui+2,
							ui+0, ui+2, ui+3,
						)
					}
				}

				cb[0] = true
			}

			cb[1] = true
		}

		cb[2] = true
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, m.Vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.Ebo)
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, len(indices)*4, gl.Ptr(indices))
	m.ElementCount = int32(len(indices))
	println("BuildChunkMesh: ElementCount =", m.ElementCount)

	return nil
}

func (self *WorldRenderer) BuildChunkMeshes(w *world.World, blockRepo *BlockRepo) {
	self.ChunkMeshes = make([]ChunkMesh, 0, len(w.Chunks))
	for key := range w.Chunks {
		println("building chunk", key[0], key[1], key[2])
		mesh := NewChunkMesh(world.CHUNK_SIZE, world.CHUNK_SIZE, world.CHUNK_SIZE)
		BuildChunkMesh(&mesh, key[0], key[1], key[2], w, blockRepo)
		self.ChunkMeshes = append(self.ChunkMeshes, mesh)
	}
}

func (self *WorldRenderer) Render(projection, view mgl32.Mat4) {
	self.Shader.UseProgram()
	self.Shader.SetUniformMatrix4fv("uProjection", projection)
	self.Shader.SetUniformMatrix4fv("uView", view)
	for _, mesh := range self.ChunkMeshes {
		gl.BindVertexArray(mesh.Vao)
		gl.DrawElements(gl.TRIANGLES, mesh.ElementCount, gl.UNSIGNED_INT, nil)
	}
}
