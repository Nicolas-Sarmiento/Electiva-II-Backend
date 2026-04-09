package cache

import (
	"time"

	goCache "github.com/patrickmn/go-cache"
)

var (
	// AppCache is the global memory cache for the DB read operations
	AppCache *goCache.Cache
)

func InitCache() {
	// Initialize cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	AppCache = goCache.New(5*time.Minute, 10*time.Minute)
}
