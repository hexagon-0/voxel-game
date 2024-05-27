package space

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Basis mgl32.Mat3

func (self Basis) X() mgl32.Vec3 {
	return mgl32.Mat3(self).Col(0)
}

func (self Basis) Y() mgl32.Vec3 {
	return mgl32.Mat3(self).Col(1)
}

func (self Basis) Z() mgl32.Vec3 {
	return mgl32.Mat3(self).Col(2)
}

type Transform struct {
	Origin mgl32.Vec3
	Basis Basis
}

func (self Transform) Mat4() mgl32.Mat4 {
	return mgl32.Mat4FromCols(
		self.Basis.X().Vec4(0),
		self.Basis.Y().Vec4(0),
		self.Basis.Z().Vec4(0),
		self.Origin.Vec4(1),
	)
}

func (self Transform) ViewMatrix() mgl32.Mat4 {
	return mgl32.Mat4FromRows(
		self.Basis.X().Vec4(-self.Basis.X().Dot(self.Origin)),
		self.Basis.Y().Vec4(-self.Basis.Y().Dot(self.Origin)),
		self.Basis.Z().Vec4(-self.Basis.Z().Dot(self.Origin)),
		mgl32.Vec4{0, 0, 0, 1},
	)
}

