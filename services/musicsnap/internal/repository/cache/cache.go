package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/service/ports"
	"sync"
	"time"
	"unsafe"
)

var _ ports.ProfileCache = &bannerMapCache{}

type bannerMapCache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	storage           map[Key]domain.CachedBanner
}

type Key struct {
	Tag     int
	Feature int
}

func New(config *Config) ports.ProfileCache {

	if config.InitialSize == 0 {
		// Для 1000 тэгов и 1000 фичей
		config.InitialSize = 1_000_000
	}

	items := make(map[Key]domain.CachedBanner, config.InitialSize)

	cache := bannerMapCache{
		storage:           items,
		defaultExpiration: config.Expiration,
	}

	return &cache
}
func (c *bannerMapCache) Set(tag, feature int, banner domain.CachedBanner, duration time.Duration) {

	if duration == 0 {
		duration = c.defaultExpiration
	}

	// Устанавливаем время истечения баннера-кеша
	if duration > 0 {
		banner.Expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()
	defer c.Unlock()

	c.storage[Key{Tag: tag, Feature: feature}] = banner

}

func (c *bannerMapCache) Get(tag, feature int) (*domain.CachedBanner, bool) {

	c.RLock()

	defer c.RUnlock()

	banner, found := c.storage[Key{
		Tag:     tag,
		Feature: feature,
	}]

	if !found {
		return nil, false
	}

	if banner.Expiration > 0 {
		if time.Now().UnixNano() > banner.Expiration {
			return nil, false
		}
	}

	return &banner, true
}

func (c *bannerMapCache) DeleteWithStatus(tag, feature int) error {

	c.Lock()

	defer c.Unlock()

	key := Key{
		Tag:     tag,
		Feature: feature}

	if _, found := c.storage[key]; !found {
		return errors.New(fmt.Sprintf("key {%d %d} not found", tag, feature))
	}

	delete(c.storage, key)

	return nil
}

func (c *bannerMapCache) Delete(tag, feature int) {

	c.Lock()

	defer c.Unlock()

	delete(c.storage, Key{
		Tag:     tag,
		Feature: feature})

	return
}

// expiredKeys returns a slice of keys of expired banners
func (c *bannerMapCache) expiredKeys() (keys []Key) {

	c.RLock()

	defer c.RUnlock()

	for k, b := range c.storage {
		if time.Now().UnixNano() > b.Expiration && b.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

func (c *bannerMapCache) clearItems(keys []Key) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.storage, k)
	}
}

func (c *bannerMapCache) Clean() {
	if keys := c.expiredKeys(); len(keys) != 0 {
		c.clearItems(keys)
	}
}

// TODO make counter for memory size instead of this
func (c *bannerMapCache) memorySizeBytes() int {

	l := len(c.storage)
	val := unsafe.Sizeof(domain.CachedBanner{
		ID:       0,
		Content:  json.RawMessage{},
		IsActive: false,
	}) + 2*unsafe.Sizeof("temp")
	key := unsafe.Sizeof(Key{1, 1})

	return 8*l + 8*l*int(key) + 8*l*int(val)
}
