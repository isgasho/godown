package server

import (
	"log"
	"time"

	"github.com/namreg/godown-v2/internal/pkg/storage"
	"github.com/namreg/godown-v2/pkg/clock"
)

//gc is the garbage collector that collects expired values
type gc struct {
	logger   *log.Logger
	clck     clock.Clock
	strg     storage.Storage
	interval time.Duration
	ticker   *time.Ticker
}

func newGc(strg storage.Storage, logger *log.Logger, clck clock.Clock, interval time.Duration) *gc {
	return &gc{
		logger:   logger,
		clck:     clck,
		strg:     strg,
		interval: interval,
	}
}

func (g *gc) start() {
	g.ticker = time.NewTicker(g.interval)

	for range g.ticker.C {
		g.deleteExpired()
	}
}

func (g *gc) stop() {
	if g.ticker != nil {
		g.ticker.Stop()
	}
}

func (g *gc) deleteExpired() {
	g.strg.RLock()
	items, err := g.strg.AllWithTTL()
	g.strg.RUnlock()
	if err != nil {
		g.logger.Printf("[WARN] gc: could not retrieve values: %v", err)
	}
	for k, v := range items {
		if v.IsExpired(g.clck.Now()) {
			g.strg.Lock()
			if err := g.strg.Del(k); err != nil {
				g.logger.Printf("[WANR] gc: could not delete item: %v", err)
			}
			g.strg.Unlock()
		}
	}
}
