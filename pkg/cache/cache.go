package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/UArt-project/UArt-proxy/domain/marketdomain"
)

var errUserNotInCache = errors.New("the user isn't in cache")

type cachedMarketPage struct {
	items             []marketdomain.MarketItem
	expireAtTimestamp int64
}

type LocalCache struct {
	stop        chan struct{}
	wg          *sync.WaitGroup
	mu          *sync.RWMutex
	marketItems map[int]cachedMarketPage
}

func NewLocalCache(cleanupInterval time.Duration) *LocalCache {
	localCache := &LocalCache{
		stop:        make(chan struct{}),
		wg:          new(sync.WaitGroup),
		mu:          new(sync.RWMutex),
		marketItems: make(map[int]cachedMarketPage),
	}

	localCache.wg.Add(1)

	go func(cleanupInterval time.Duration) {
		defer localCache.wg.Done()
		localCache.cleanupLoop(cleanupInterval)
	}(cleanupInterval)

	return localCache
}

func (lc *LocalCache) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for {
		select {
		case <-lc.stop:
			return
		case <-ticker.C:
			lc.mu.Lock()

			for uid, cu := range lc.marketItems {
				if cu.expireAtTimestamp <= time.Now().Unix() {
					delete(lc.marketItems, uid)
				}
			}

			lc.mu.Unlock()
		}
	}
}

func (lc *LocalCache) stopCleanup() { //nolint:unused
	close(lc.stop)

	lc.wg.Wait()
}

func (lc *LocalCache) Update(pageID int, mItems []marketdomain.MarketItem, expireAtTimestamp int64) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.marketItems[pageID] = cachedMarketPage{
		items:             mItems,
		expireAtTimestamp: expireAtTimestamp,
	}
}

func (lc *LocalCache) Read(id int) ([]marketdomain.MarketItem, error) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	cu, ok := lc.marketItems[id]
	if !ok {
		return nil, errUserNotInCache
	}

	return cu.items, nil
}

func (lc *LocalCache) delete(id int) { //nolint:unused
	lc.mu.Lock()
	defer lc.mu.Unlock()

	delete(lc.marketItems, id)
}
