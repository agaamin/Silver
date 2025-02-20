package syncmanager

import (
    "log"
    "time"
    "context"
    "fingertip/internal/config"
    "silverbullet/internal/resolvers/spvproc"
    "silverbullet/internal/cache"
    "os"
)

const (
    SyncRetryInterval = 10 * time.Minute  // Retry every 10 minutes on failure
    SyncTimeout       = 10 * time.Minute  // Stop sync attempts after 10 minutes if no progress
    MaxRetries        = 3                 // Maximum retries before prompting user
    WarningThreshold  = 72 * time.Hour    // Notify user after 3 days of failure
)

type SyncManager struct {
    proc        *spvproc.SPVProc  // Replacing HNSD with SPV-based resolver
    lastSync    time.Time
    cache       *cache.CacheManager
    retryCount  int
    syncInterval time.Duration
}

func NewSyncManager(proc *spvproc.SPVProc, cache *cache.CacheManager, syncInterval time.Duration) *SyncManager {
    return &SyncManager{
        proc:  proc,
        cache: cache,
        lastSync: cache.GetLastSyncTime(),
        retryCount: 0,
        syncInterval: syncInterval,
    }
}

func (s *SyncManager) StartSync(ctx context.Context) {
    log.Println("[SyncManager] Starting SPV sync...")
    syncStart := time.Now()

    if !s.cache.ValidateCache(WarningThreshold) {
        log.Println("[SyncManager] WARNING: Cache is too old. Proceeding with full sync.")
    } else {
        log.Println("[SyncManager] Cache is fresh. Skipping sync.")
        return
    }

    if time.Since(s.lastSync) < s.syncInterval {
        log.Println("[SyncManager] Last sync was recent. Performing lightweight health check.")
        if s.performHealthCheck() {
            log.Println("[SyncManager] No blockchain updates detected. Skipping full sync.")
            return
        }
    }

    log.Println("[SyncManager] Checking if SPV resolver API is running...")

    // Try to get the latest block height to confirm API is active
    _, err := s.proc.GetLatestBlock()
    if err != nil {
        log.Println("[SyncManager] ERROR: SPV resolver is not running or unreachable. Skipping sync.")
        return // Do not proceed with sync if resolver is down
    }

    log.Println("[SyncManager] SPV resolver is active. Initiating SPV sync process...")
    done := make(chan struct{})
    go func() {
        defer close(done)
        err := s.proc.Sync()
        if err != nil {
            log.Println("[SyncManager] Sync error:", err)
        }
    }()

    select {
    case <-done:
        log.Println("[SyncManager] SPV Sync complete in", time.Since(syncStart))
        s.lastSync = time.Now()
        s.cache.UpdateLastSyncTime(s.lastSync)
        s.retryCount = 0 // Reset retry count after successful sync
    case <-time.After(SyncTimeout):
        log.Println("[SyncManager] Sync timeout reached. Using cached data.")
    }
}

func (s *SyncManager) RetryFailedSync(ctx context.Context) {
    for s.retryCount < MaxRetries {
        time.Sleep(SyncRetryInterval)
        if err := s.proc.Sync(); err == nil {
            log.Println("[SyncManager] Retried SPV sync successfully.")
            s.lastSync = time.Now()
            s.cache.UpdateLastSyncTime(s.lastSync)
            s.retryCount = 0 // Reset retry count
            return
        }
        s.retryCount++
        log.Println("[SyncManager] Retry failed (attempt", s.retryCount, "), will attempt again in", SyncRetryInterval)
    }
    log.Println("[SyncManager] No internet available for Auto Sync. Please sync manually when internet is available.")
}

func (s *SyncManager) performHealthCheck() bool {
    log.Println("[SyncManager] Checking SPV blockchain state before full sync...")
    latestBlock, err := s.proc.GetLatestBlock()
    if err != nil {
        log.Println("[SyncManager] Failed to fetch latest block height. Proceeding with full sync.")
        return false
    }

    if latestBlock == s.cache.GetLastSyncTime().Unix() {
        log.Println("[SyncManager] Blockchain state unchanged. Skipping full sync.")
        return true
    }

    log.Println("[SyncManager] Blockchain state changed. Proceeding with full sync.")
    return false
}
