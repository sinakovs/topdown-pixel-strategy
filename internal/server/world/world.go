package world

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	chunkrepo "github.com/sinakovs/topdown-pixel-strategy/internal/chunkRepo"
)

const (
	ChunkCols         = uint32(32)
	ChunkRows         = uint32(32)
	WorldSizeInChunks = 16
)

func New(
	repo chunkrepo.ChunkRepository,
	chunkCols, chunkRows int,
	tickRate uint,
) (*world, error) {
	if tickRate <= 0 {
		return nil, fmt.Errorf("invalid tickrate %d", tickRate)
	}

	return &world{
		repo: repo,
		Visibility: VisibilityPolicy{
			RadiusTiles:      1,
			IncludeSelfChunk: true,
			Mode:             NeighbourMoore,
		},
		tickRate:         tickRate,
		subscribersMutex: new(sync.RWMutex),
		subscribers:      make(map[PlayerID]chan struct{}),
		players:          make(map[PlayerID]*player),
		units:            make(map[UnitID]*unit),
		chunkCache:       make(map[chunkrepo.ChunkKey]*chunkrepo.Chunk),
	}, nil
}

type world struct {
	repo chunkrepo.ChunkRepository

	// Policy controlling which chunks a player may render around each unit.
	Visibility VisibilityPolicy

	// Config
	tickRate uint

	subscribersMutex *sync.RWMutex
	subscribers      map[PlayerID]chan struct{}

	// State
	mu      sync.RWMutex
	players map[PlayerID]*player
	units   map[UnitID]*unit

	// Optional RAM cache to avoid decoding the same chunk repeatedly.
	chunkMu    sync.RWMutex
	chunkCache map[chunkrepo.ChunkKey]*chunkrepo.Chunk
}

func (w *world) Start() {
	ticker := time.NewTicker(time.Second / time.Duration(w.tickRate))
	go func() {
		for range ticker.C {
			w.update()
			w.notify()
		}
	}()

}

// TODO:
func (w *world) Stop() error {
	return nil
}

func (w *world) update() error {
	return nil
}

func (w *world) notify() {
	w.subscribersMutex.RLock()
	defer w.subscribersMutex.RUnlock()

	for _, subChan := range w.subscribers {
		select {
		case subChan <- struct{}{}:
		default:
		}
	}
}

func (w *world) Subscribe(playerId PlayerID) <-chan struct{} {
	updateCh := make(chan struct{}, 1)

	w.subscribersMutex.Lock()
	defer w.subscribersMutex.Unlock()

	w.subscribers[playerId] = updateCh

	return updateCh
}

func (w *world) Unsubscribe(playerId PlayerID) {
	w.subscribersMutex.Lock()
	defer w.subscribersMutex.Unlock()

	updateCh := w.subscribers[playerId]

	delete(w.subscribers, playerId)

	close(updateCh)
}

// ChunksForPlayer returns the authoritative list of chunks the player is allowed
// to render given the current positions of their units and the visibility policy.
// It loads any missing chunks via the repository (and caches them).
func (w *world) ChunksForPlayer(ctx context.Context, playerID PlayerID) ([]chunkrepo.RenderChunk, error) {
	// 1) Snapshot player units (avoid holding the lock while loading chunks)
	unitPositions := w.unitsOwnedBy(playerID)

	// 2) Compute the set of chunk keys to render (dedup)
	keys := make(map[chunkrepo.ChunkKey]struct{}, 16)
	for _, pos := range unitPositions {
		unitChunkKey, _, _ := w.worldTileToChunkKey(pos.WorldCol, pos.WorldRow)
		for _, key := range w.Visibility.neighbours(unitChunkKey.X, unitChunkKey.Y, WorldSizeInChunks-1, WorldSizeInChunks-1) {
			keys[key] = struct{}{}
		}
		// If IncludeSelfChunk was false but you still want at least the unit's current
		// chunk, uncomment:
		// if !w.Visibility.IncludeSelfChunk {
		//     keys[unitChunkKey] = struct{}{}
		// }
	}

	// 3) Resolve keys -> chunks (load or fetch from cache)
	result := make([]chunkrepo.RenderChunk, 0, len(keys))
	for key := range keys {
		ch, err := w.ensureChunk(ctx, key)
		if err != nil {
			return nil, err
		}
		result = append(
			result,
			chunkrepo.RenderChunk{
				Key:   key,
				Chunk: ch,
			},
		)
	}
	return result, nil
}

func (w *world) unitsOwnedBy(playerID PlayerID) []*unit {
	w.mu.RLock()
	defer w.mu.RUnlock()
	plr, ok := w.players[playerID]
	if !ok {
		return nil
	}
	out := make([]*unit, 0, len(plr.UnitsIDs))
	for uid := range plr.UnitsIDs {
		if u, ok := w.units[uid]; ok {
			out = append(out, &unit{
				ID:       u.ID,
				OwnerID:  u.OwnerID,
				WorldCol: u.WorldCol,
				WorldRow: u.WorldRow,
			})
		}
	}
	return out
}

func (w *world) ensureChunk(ctx context.Context, key chunkrepo.ChunkKey) (*chunkrepo.Chunk, error) {
	// RAM cache fast path
	w.chunkMu.RLock()
	if ch, ok := w.chunkCache[key]; ok {
		w.chunkMu.RUnlock()
		return ch, nil
	}
	w.chunkMu.RUnlock()

	// Load from repository
	ch, err := w.repo.Load(ctx, key)
	if err == nil {
		w.chunkMu.Lock()
		w.chunkCache[key] = ch
		w.chunkMu.Unlock()
		return ch, nil
	}
	if err != nil && errors.Is(err, chunkrepo.ErrChunkNotFound) {
		return nil, err // needs to be modified
	}

	w.chunkMu.Lock()
	w.chunkCache[key] = ch
	w.chunkMu.Unlock()
	return ch, nil
}

func (w *world) worldTileToChunkKey(worldCol, worldRow uint32) (key chunkrepo.ChunkKey, localCol, localRow uint32) {
	cx := uint32(math.Floor(float64(worldCol) / float64(ChunkCols)))
	cy := uint32(math.Floor(float64(worldRow) / float64(ChunkRows)))
	lc := worldCol - cx*ChunkCols
	lr := worldRow - cy*ChunkRows

	return chunkrepo.ChunkKey{X: cx, Y: cy}, lc, lr
}
