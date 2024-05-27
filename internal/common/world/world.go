package world

import (
	"fmt"
	"math"
)

func generate(x, y, z int) BlockId {
	// if x&1 == y&1 && y&1 == z&1 {
	if y < int(math.Abs(math.Sin(float64(z)/16.0*math.Pi)) * 16) {
		return 1
	}

	return 0
}

func GenerateChunk(x, y, z int, width, height, depth uint, genFn func(w, h, d int) BlockId) *Chunk {
	chunk := Chunk{width, height, depth, make([]BlockId, width*height*depth)}

	for k := uint(0); k < depth; k++ {
		for j := uint(0); j < height; j++ {
			for i := uint(0); i < width; i++ {
				index := i + j*width + k*width*height
				chunk.Blocks[index] = genFn(x+int(i), y+int(k), z+int(j))
			}
		}
	}

	return &chunk
}

type World struct {
	Width, Height, Depth uint
	Chunks               map[[3]int]*Chunk
}

func (self *World) LoadChunk(x, y, z int) {
	x *= CHUNK_SIZE
	y *= CHUNK_SIZE
	z *= CHUNK_SIZE
	self.Chunks[[3]int{x, y, z}] = GenerateChunk(x, y, z, CHUNK_SIZE, CHUNK_SIZE, CHUNK_SIZE, generate)
}

type ChunkNotLoadedError [3]int

func (self *ChunkNotLoadedError) Error() string {
	return fmt.Sprintf("Chunk not loaded: %d %d %d", self[0], self[1], self[2])
}

func (self *World) BlockAt(x, y, z int) (BlockId, *ChunkNotLoadedError) {
	cx, cy, cz := x>>CHUNK_SHIFT, y>>CHUNK_SHIFT, z>>CHUNK_SHIFT
	chunk, ok := self.Chunks[[3]int{cx, cy, cz}]

	if ok {
		return chunk.BlockAt(uint(x-cx), uint(y-cy), uint(z-cz)), nil
	}

	return 0, &ChunkNotLoadedError{x, y, z}
}

