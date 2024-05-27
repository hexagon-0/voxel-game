package world

import "math"

func sign(n float64) float64 {
	if n > 0.0 {
		return 1.0
	}

	if n < 0.0 {
		return -1.0
	}

	return 0.0
}

func frac0(n float64) float64 {
	return n - math.Floor(n)
}

func frac1(n float64) float64 {
	return 1 - n + math.Floor(n)
}

type VoxelRaycast struct {
	X, Y, Z int
	step    [3]float64
	tDelta  [3]float64
	tMax    [3]float64
}

func (self *VoxelRaycast) Init(from, to [3]float64) {
	posInf := 1e7
	d := [3]float64{}

	for i := 0; i < 3; i++ {
		d[i] = to[i] - from[i]
		self.step[i] = sign(d[i])
		if self.step[i] != 0 {
			self.tDelta[i] = math.Min(self.step[i] / d[i], posInf)
		} else {
			self.tDelta[i] = posInf
		}
		if self.step[i] > 0 {
			self.tMax[i] = self.tDelta[i] * frac1(from[i])
		} else {
			self.tMax[i] = self.tDelta[i] * frac0(from[i])
		}
	}
	self.X = int(math.Floor(from[0]))
	self.Y = int(math.Floor(from[1]))
	self.Z = int(math.Floor(from[2]))
}

func (self *VoxelRaycast) Step() {
	// if self.tMax[0] > 1 && self.tMax[1] > 1 && self.tMax[2] > 1 {
	// 	return
	// }

	if self.tMax[0] < self.tMax[1] {
		if self.tMax[0] < self.tMax[2] {
			self.X += int(self.step[0])
			self.tMax[0] += self.tDelta[0]
		} else {
			self.Z += int(self.step[2])
			self.tMax[2] += self.tDelta[2]
		}
	} else {
		if self.tMax[1] < self.tMax[2] {
			self.Y += int(self.step[1])
			self.tMax[1] += self.tDelta[1]
		} else {
			self.Z += int(self.step[2])
			self.tMax[2] += self.tDelta[2]
		}
	}
}
