package world

import chunkrepo "github.com/sinakovs/topdown-pixel-strategy/internal/chunkRepo"

type NeighborMode uint8

const (
	// VonNeumann = 4-neighbourhood (N, S, E, W)
	NeighbourVonNeumann NeighborMode = iota
	// Moore = 8-neighbourhood (N, NE, E, SE, S, SW, W, NW) including diagonals
	NeighbourMoore
)

type VisibilityPolicy struct {
	RadiusTiles      uint32
	IncludeSelfChunk bool
	Mode             NeighborMode
}

// neighbours returns neighbouring chunk keys around (cx,cy),
// clamped to [0..maxCX] × [0..maxCY]. Works with non-negative worlds.
func (p VisibilityPolicy) neighbours(cx, cy, maxCX, maxCY uint32) []chunkrepo.ChunkKey {
	r := int(p.RadiusTiles) // use signed for loops
	includeSelf := p.IncludeSelfChunk
	mode := p.Mode

	// Worst-case capacity hint.
	keys := make([]chunkrepo.ChunkKey, 0, (2*r+1)*(2*r+1))

	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			// Skip center if requested
			if dx == 0 && dy == 0 && !includeSelf {
				continue
			}
			// Neighbourhood shape
			if mode == NeighbourVonNeumann && absInt(dx)+absInt(dy) > r {
				continue
			}

			// Signed math, then clamp to [0..max]
			nx := int64(int64(cx) + int64(dx))
			ny := int64(int64(cy) + int64(dy))
			if nx < 0 || ny < 0 {
				continue
			}
			if nx > int64(maxCX) || ny > int64(maxCY) {
				continue
			}

			keys = append(keys, chunkrepo.ChunkKey{
				X: uint32(nx),
				Y: uint32(ny),
			})
		}
	}
	return keys
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
