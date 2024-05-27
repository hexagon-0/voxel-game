package world

type BlockId uint8

const (
	CHUNK_SHIFT = 5
	CHUNK_SIZE = 1<<CHUNK_SHIFT
)

type Chunk struct {
	Width, Height, Depth uint
	Blocks   []BlockId
}

func (self *Chunk) BlockAt(x, y, z uint) BlockId {
	return self.Blocks[x + z*self.Width + y*self.Width*self.Depth]
}

