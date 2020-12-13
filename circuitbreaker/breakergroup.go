package circuitbreaker

import (
	"sync"
)

type BreakerGroup struct {
	BreakerMap   map[string]*Breaker
	groupRWMutex sync.RWMutex
}

func NewBreakerGroup() *BreakerGroup {
	return &BreakerGroup{
		BreakerMap: make(map[string]*Breaker),
	}
}

func (bg *BreakerGroup) Add(name string, b *Breaker) {
	bg.groupRWMutex.Lock()
	bg.BreakerMap[name] = b
	bg.groupRWMutex.Unlock()
}

func (bg *BreakerGroup) Get(name string) *Breaker {
	bg.groupRWMutex.RLock()
	b, ok := bg.BreakerMap[name]
	bg.groupRWMutex.RUnlock()
	if !ok {
		return nil
	}

	return b
}
