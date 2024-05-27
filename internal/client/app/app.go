package app

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"math"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hexagon-0/voxel-game/internal/client/render"
	"github.com/hexagon-0/voxel-game/internal/common/world"
)

const (
	WIDTH  = 1600
	HEIGHT = 900

	FPS = 60

	LOOK_SENSITIVITY = 0.025

	WORLD_HEIGHT = 2
	WORLD_SIZE   = 4
)

func resizeCallback(window *glfw.Window, w, h int) {
	gl.Viewport(0, 0, int32(w), int32(h))
}

type App struct {
	window        *glfw.Window
	world         world.World
	worldRenderer render.WorldRenderer
	raycast       world.VoxelRaycast
	flag          bool
}

func (self *App) Init() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	self.window, err = glfw.CreateWindow(WIDTH, HEIGHT, "Hello OpenGL", nil, nil)
	if err != nil {
		panic(err)
	}

	videoMode := glfw.GetPrimaryMonitor().GetVideoMode()
	self.window.SetPos((videoMode.Width-WIDTH)/2, (videoMode.Height-HEIGHT)/2)

	self.window.MakeContextCurrent()

	self.window.SetFramebufferSizeCallback(resizeCallback)
	self.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
}

func (self *App) Deinit() {
	glfw.Terminate()
}

func LoadTexture(path string) (uint32, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("Unsupported image stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// flip vertically
	tmp := image.NewRGBA(rgba.Bounds())
	height := rgba.Rect.Size().Y
	s := rgba.Stride
	for i := 0; i < height; i++ {
		j := height - i
		copy(tmp.Pix[(j-1)*s:j*s], rgba.Pix[i*s:(i+1)*s])
	}
	rgba = tmp

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)

	return texture, nil
}

func (self *App) Run() {
	self.Init()
	defer self.Deinit()

	blockRepo := render.BlockRepo(make(map[world.BlockId]image.Rectangle))
	blockAtlasTexture, err := LoadTexture("assets/textures/block_atlas.png")
	if err != nil {
		panic(err)
	}
	_ = blockAtlasTexture
	blockRepo[1] = image.Rectangle{image.Point{0, 0}, image.Point{16, 16}}
	blockRepo[2] = image.Rectangle{image.Point{16, 0}, image.Point{32, 16}}
	blockRepo[3] = image.Rectangle{image.Point{32, 0}, image.Point{48, 16}}
	blockRepo[4] = image.Rectangle{image.Point{48, 0}, image.Point{64, 16}}

	// commonDirtTexture, err := LoadTexture("assets/textures/common_dirt.png")
	commonDirtTexture, err := LoadTexture("assets/textures/uv_test.png")
	if err != nil {
		panic(err)
	}
	gl.BindTexture(gl.TEXTURE_2D, commonDirtTexture)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// err := self.worldRenderer.CompileShaders()
	// if err != nil {
	// 	panic(err)
	// }

	// Initialize world
	self.world.Chunks = make(map[[3]int]*world.Chunk)
	tcCoords := [3]int{0, 0, 0} // test chunk coordinates
	self.world.LoadChunk(tcCoords[0], tcCoords[1], tcCoords[2])
	// self.worldRenderer.BuildChunkMeshes(self.world)
	chunkMesh := render.NewChunkMesh(world.CHUNK_SIZE, world.CHUNK_SIZE, world.CHUNK_SIZE)
	err = render.BuildChunkMesh(&chunkMesh, tcCoords[0], tcCoords[1], tcCoords[2], &self.world, &blockRepo)
	if err != nil {
		panic(err)
	}

	// test

	// chunk shader
	vs, err := render.NewShader(render.ChunkVsSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fs, err := render.NewShader(render.ChunkFsSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	shader, err := render.NewShaderProgram(vs, fs)
	if err != nil {
		panic(err)
	}

	// selection shader
	vs, err = render.NewShader(render.SelectionVsSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fs, err = render.NewShader(render.SelectionFsSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	selectShader, err := render.NewShaderProgram(vs, fs)
	if err != nil {
		panic(err)
	}

	vertices := []uint32{
		// +X
		1, 0, 1, 0,
		1, 0, 0, 0,
		1, 1, 1, 0,
		1, 1, 0, 0,

		// -X
		0, 0, 0, 1,
		0, 0, 1, 1,
		0, 1, 0, 1,
		0, 1, 1, 1,

		// +Y
		0, 1, 1, 2,
		1, 1, 1, 2,
		0, 1, 0, 2,
		1, 1, 0, 2,

		// -Y
		0, 0, 0, 3,
		1, 0, 0, 3,
		0, 0, 1, 3,
		1, 0, 1, 3,

		// +Z
		0, 0, 1, 4,
		1, 0, 1, 4,
		0, 1, 1, 4,
		1, 1, 1, 4,

		// -Z
		1, 0, 0, 5,
		0, 0, 0, 5,
		1, 1, 0, 5,
		0, 1, 0, 5,
	}

	indices := []uint32{
		// +X
		0, 1, 2,
		1, 3, 2,

		// -X
		4, 5, 6,
		5, 7, 6,

		// +Y
		8, 9, 10,
		9, 11, 10,

		// -Y
		12, 13, 14,
		13, 15, 14,

		// +Z
		16, 17, 18,
		17, 19, 18,

		// -Z
		20, 21, 22,
		21, 23, 22,
	}

	var vbo, ebo, vao uint32
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)
	gl.GenVertexArrays(1, &vao)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribIPointerWithOffset(0, 3, gl.UNSIGNED_INT, 4*4, 0)

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribIPointerWithOffset(1, 1, gl.UNSIGNED_INT, 4*4, 3*4)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.BindVertexArray(0)

	// end: test

	// UI

	vs, err = render.NewShader(render.UiVsSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fs, err = render.NewShader(render.UiFsSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	uiShader, err := render.NewShaderProgram(vs, fs)
	if err != nil {
		panic(err)
	}

	orthoMatrix := mgl32.Ortho2D(0, WIDTH, 0, HEIGHT)
	err = uiShader.SetUniformMatrix4fv("uProjection", orthoMatrix)
	if err != nil {
		panic(err)
	}
	// err = uiShader.SetUniformMatrix4fv("uView", mgl32.Ident4())
	// if err != nil {
	// 	panic(err)
	// }
	err = uiShader.SetUniformMatrix4fv("uModel", mgl32.Translate3D(float32(WIDTH)/2, float32(HEIGHT)/2, 0.0))
	// err = uiShader.SetUniformMatrix4fv("uModel", mgl32.Ident4())
	if err != nil {
		panic(err)
	}
	err = uiShader.SetUniform3fv("uColor", mgl32.Vec3{0.4, 0.4, 0.8})
	if err != nil {
		panic(err)
	}

	var retVbo, retVao uint32
	var retExt float32 = 4.0
	retVertices := []float32{
		-retExt, -retExt,
		 retExt, -retExt,
		-retExt,  retExt,
		 retExt, -retExt,
		 retExt,  retExt,
		-retExt,  retExt,
	}
	gl.GenBuffers(1, &retVbo)
	gl.GenVertexArrays(1, &retVao)

	gl.BindBuffer(gl.ARRAY_BUFFER, retVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(retVertices)*4, gl.Ptr(retVertices), gl.STATIC_DRAW)

	gl.BindVertexArray(retVao)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	gl.BindVertexArray(0)

	// Main loop timing
	timePerTick := 1.0 / FPS
	delta := 0.0
	lastTime := glfw.GetTime()

	// Camera
	aspectRatio := float64(WIDTH) / HEIGHT
	fovx := math.Pi / 2.0
	fovy := math.Atan(math.Tan(fovx/2) * aspectRatio)
	projectionMatrix := mgl32.Perspective(float32(fovy), float32(aspectRatio), 0.001, 1000.0)

	cameraPos := mgl32.Vec3{0.0, 0.0, -5.0}
	cameraFront := mgl32.Vec3{0.0, 0.0, -1.0}
	cameraSpeed := 2.5

	px, py := self.window.GetCursorPos()
	dx, dy := 0.0, 0.0
	yaw, pitch := float64(mgl32.DegToRad(90.0)), 0.0
	maxPitch := float64(mgl32.DegToRad(89.0))

	for !self.window.ShouldClose() {
		now := glfw.GetTime()
		delta += now - lastTime
		lastTime = now

		for delta >= timePerTick {
			delta -= timePerTick
			if self.window.GetKey(glfw.KeyEscape) == glfw.Press {
				self.window.SetShouldClose(true)
			}

			mx, my := self.window.GetCursorPos()
			dx, dy = mx-px, py-my
			px, py = mx, my

			yaw += dx * LOOK_SENSITIVITY * timePerTick
			pitch += dy * LOOK_SENSITIVITY * timePerTick
			if pitch > maxPitch {
				pitch = maxPitch
			} else if pitch < -maxPitch {
				pitch = -maxPitch
			}

			cameraFront = mgl32.Vec3{
				float32(math.Cos(yaw) * math.Cos(pitch)),
				float32(math.Sin(pitch)),
				float32(math.Sin(yaw) * math.Cos(pitch)),
			}.Normalize()
			cameraRight := cameraFront.Cross(render.WorldUp).Normalize()
			cameraUp := cameraRight.Cross(cameraFront)

			cameraTranslation := mgl32.Vec3{}
			if self.window.GetKey(glfw.KeyW) == glfw.Press {
				cameraTranslation = cameraTranslation.Add(cameraFront)
			}
			if self.window.GetKey(glfw.KeyS) == glfw.Press {
				cameraTranslation = cameraTranslation.Sub(cameraFront)
			}
			if self.window.GetKey(glfw.KeyD) == glfw.Press {
				cameraTranslation = cameraTranslation.Add(cameraRight)
			}
			if self.window.GetKey(glfw.KeyA) == glfw.Press {
				cameraTranslation = cameraTranslation.Sub(cameraRight)
			}

			cameraTranslation[1] = 0.0
			if cameraTranslation.LenSqr() != 0.0 {
				cameraTranslation = cameraTranslation.Normalize()
			}

			if self.window.GetKey(glfw.KeySpace) == glfw.Press {
				cameraTranslation = cameraTranslation.Add(render.WorldUp)
			}
			if self.window.GetKey(glfw.KeyLeftShift) == glfw.Press {
				cameraTranslation = cameraTranslation.Sub(render.WorldUp)
			}

			speedMult := 1.0
			if self.window.GetKey(glfw.KeyLeftControl) == glfw.Press {
				speedMult = 4.0
			}
			cameraPos = cameraPos.Add(cameraTranslation.Mul(float32(cameraSpeed * speedMult * timePerTick)))

			// viewMatrix := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), render.WorldUp)
			viewMatrix := mgl32.Mat4{
				cameraRight[0], cameraUp[0], -cameraFront[0], 0,
				cameraRight[1], cameraUp[1], -cameraFront[1], 0,
				cameraRight[2], cameraUp[2], -cameraFront[2], 0,
				-cameraRight.Dot(cameraPos), -cameraUp.Dot(cameraPos), cameraFront.Dot(cameraPos), 1,
			}

			// Raycast
			dst := cameraPos.Add(cameraFront.Mul(6))
			self.raycast.Init(
				[3]float64{ float64(cameraPos[0]), float64(cameraPos[1]), float64(cameraPos[2]) },
				[3]float64{ float64(dst[0]),       float64(dst[1]),       float64(dst[2]) },
			)

			// DEBUG: Raycast
			printDebugInfo := false
			if self.window.GetKey(glfw.KeyH) == glfw.Press {
				if !self.flag {
					printDebugInfo = true
				}
				self.flag = true
			} else {
				self.flag = false
			}

			if printDebugInfo {
				fmt.Println("DEBUG: Raycast")
				fmt.Println("From:")
				fmt.Printf("%.2f %.2f %.2f\n",
					float64(cameraPos[0]), float64(cameraPos[1]), float64(cameraPos[2]),
				)
				fmt.Println("To:")
				fmt.Printf("%.2f %.2f %.2f\n",
					float64(dst[0]), float64(dst[1]), float64(dst[2]),
				)
			}

			raycastHit := false
			var selModelMatrix mgl32.Mat4
			var blockId world.BlockId

			// Render

			gl.ClearColor(0.068, 0.068, 0.098, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			// object.Render(projectionMatrix, viewMatrix)
			// self.worldRenderer.Render(projectionMatrix, viewMatrix)
			shader.UseProgram()
			shader.SetUniformMatrix4fv("uProjection", projectionMatrix)
			shader.SetUniformMatrix4fv("uView", viewMatrix)
			shader.SetUniformMatrix4fv("uModel", chunkMesh.ModelMatrix)
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, commonDirtTexture)
			gl.BindVertexArray(chunkMesh.Vao)
			gl.DrawElements(gl.TRIANGLES, chunkMesh.ElementCount, gl.UNSIGNED_INT, nil)

			for i := 0; i < 6; i++ {
				blockId, err = self.world.BlockAt(int(self.raycast.X), int(self.raycast.Y), int(self.raycast.Z))

				if printDebugInfo {
					fmt.Println("Step", i)
					fmt.Printf("Position: %d %d %d\n", self.raycast.X, self.raycast.Y, self.raycast.Z)
					// fmt.Printf("Position: %.2f %.2f %.2f\n", self.raycast.X, self.raycast.Y, self.raycast.Z)
					fmt.Printf("BlockId: %d Loaded: %v\n", blockId, err)
				}

				// TODO: Probably enable this once we're loading chunks around the player
				// if err != nil {
				// 	break
				// }

				// selectShader.UseProgram()
				// selectShader.SetUniformMatrix4fv("uProjection", projectionMatrix)
				// selectShader.SetUniformMatrix4fv("uView", viewMatrix)
				// selModelMatrix = mgl32.Translate3D(float32(self.raycast.X), float32(self.raycast.Y), float32(self.raycast.Z))
				// selectShader.SetUniformMatrix4fv("uModel", selModelMatrix)
				// gl.Enable(gl.POLYGON_OFFSET_FILL)
				// gl.PolygonOffset(-1.0, -1.0)
				// gl.BindVertexArray(vao)
				// gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)
				// gl.Disable(gl.POLYGON_OFFSET_FILL)

				if blockId != 0 {
					raycastHit = true
					selModelMatrix = mgl32.Translate3D(float32(self.raycast.X), float32(self.raycast.Y), float32(self.raycast.Z))
					break
				}

				self.raycast.Step()
			}
			// _ = raycastHit
			if raycastHit {
				selectShader.UseProgram()
				selectShader.SetUniformMatrix4fv("uProjection", projectionMatrix)
				selectShader.SetUniformMatrix4fv("uView", viewMatrix)
				selectShader.SetUniformMatrix4fv("uModel", selModelMatrix)
				gl.Enable(gl.POLYGON_OFFSET_FILL)
				gl.PolygonOffset(-1.0, -1.0)
				gl.BindVertexArray(vao)
				gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)
				gl.Disable(gl.POLYGON_OFFSET_FILL)
			}

			uiShader.UseProgram()
			gl.BindVertexArray(retVao)
			gl.DrawArrays(gl.TRIANGLES, 0, 6)

			gl.BindVertexArray(0)

			self.window.SwapBuffers()
			glfw.PollEvents()
		}
	}
}
