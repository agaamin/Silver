package cache

import (
    "log"
    "os"
    "time"
    "encoding/json"
)

const CacheFile = "sync_cache.json"

type CacheManager struct {
    lastSync time.Time
    valid    bool
}

func NewCacheManager() *CacheManager {
    cm := &CacheManager{}
    cm.loadCache()
    return cm
}

func (c *CacheManager) loadCache() {
    file, err := os.Open(CacheFile)
    if err != nil {
        log.Println("[CacheManager] No existing cache found. Initializing new cache.")
        c.valid = false
        return
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&c.lastSync)
    if err != nil {
        log.Println("[CacheManager] Failed to read cache. Initializing new cache.")
        c.valid = false
        return
    }
    c.valid = true
}

func (c *CacheManager) GetLastSyncTime() time.Time {
    return c.lastSync
}

func (c *CacheManager) UpdateLastSyncTime(syncTime time.Time) {
    c.lastSync = syncTime
    c.valid = true
    c.saveCache()
}

func (c *CacheManager) saveCache() {
    file, err := os.Create(CacheFile)
    if err != nil {
        log.Println("[CacheManager] Failed to save cache:", err)
        return
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    err = encoder.Encode(c.lastSync)
    if err != nil {
        log.Println("[CacheManager] Failed to write cache:", err)
        c.valid = false
    } else {
        c.valid = true
    }
}

func (c *CacheManager) ValidateCache(maxAge time.Duration) bool {
    if !c.valid || time.Since(c.lastSync) > maxAge {
        log.Println("[CacheManager] Cache is too old or invalid. Recommend manual sync.")
        return false
    }
    return true
}

func (c *CacheManager) ClearCache() {
    err := os.Remove(CacheFile)
    if err != nil {
        log.Println("[CacheManager] Failed to clear cache:", err)
    } else {
        log.Println("[CacheManager] Cache cleared successfully.")
        c.valid = false
    }
}
