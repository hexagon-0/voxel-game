package render

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

var WorldUp = mgl32.Vec3{0, 1, 0}

type Camera struct {
	Position mgl32.Vec3
	Pitch float64
	Yaw float64
	Projection mgl32.Mat4
}

func NewCamera(fovy, aspectRatio float32) Camera {
	position := mgl32.Vec3{0, 0, 0}

	return Camera{
		position,
		0,
		-math.Pi/2,
		mgl32.Perspective(fovy, aspectRatio, 0.001, 1000.0),
	}
}

func (self Camera) Direction() mgl32.Vec3 {
	return mgl32.Vec3{
		float32(math.Cos(self.Pitch) * math.Cos(self.Yaw)),
		float32(math.Sin(self.Yaw)),
		float32(math.Sin(self.Pitch) * math.Cos(self.Yaw)),
	}.Normalize()
}

func (self Camera) RotateX(value float64) {
	self.Pitch += value
}

func (self Camera) RotateXClamped(value, cap float64) {
	self.RotateX(value)
	if self.Pitch > cap {
		self.Pitch = cap
	} else if self.Pitch < -cap {
		self.Pitch = cap
	}
}

func (self Camera) RotateY(value float64) {
	self.Yaw += value
}

func (self Camera) ViewMatrix() mgl32.Mat4 {
	z := self.Direction()
	x := z.Cross(WorldUp).Normalize()
	y := x.Cross(z)
	z = mgl32.Vec3{}.Sub(z)
	return mgl32.Mat4FromRows(
		x.Vec4(-x.Dot(self.Position)),
		y.Vec4(-y.Dot(self.Position)),
		z.Vec4(-z.Dot(self.Position)),
		mgl32.Vec4{0, 0, 0, 1},
	)
}

