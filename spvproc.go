package spvproc

import (
    "fmt"
    "log"
    "net/http"
    "bytes"
    "encoding/json"
    "sync/atomic"
    "time"
)

const syncWaitDuration = 5 * time.Second

type SPVProc struct {
    apiURL     string
    isRunning  int32
    lastHeight int64
}

func NewSPVProc(apiURL string) *SPVProc {
    return &SPVProc{
        apiURL: apiURL,
    }
}

func (s *SPVProc) Start() error {
    if atomic.LoadInt32(&s.isRunning) == 1 {
        log.Println("[SPVProc] SPV resolver already running.")
        return nil
    }
    log.Println("[SPVProc] Starting SPV resolver...")
    atomic.StoreInt32(&s.isRunning, 1)
    return nil
}

func (s *SPVProc) Stop() error {
    if atomic.LoadInt32(&s.isRunning) == 0 {
        return nil
    }
    log.Println("[SPVProc] Stopping SPV resolver...")
    atomic.StoreInt32(&s.isRunning, 0)
    return nil
}

func (s *SPVProc) Sync() error {
    log.Println("[SPVProc] Triggering SPV sync...")
    reqBody, err := json.Marshal(map[string]string{})
    if err != nil {
        return fmt.Errorf("failed to marshal request body: %v", err)
    }
    resp, err := http.Post(s.apiURL+"/sync", "application/json", bytes.NewBuffer(reqBody))
    if err != nil {
        return fmt.Errorf("HTTP post failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
    }

    time.Sleep(syncWaitDuration)
    return nil
}

func (s *SPVProc) GetLatestBlock() (int64, error) {
    log.Println("[SPVProc] Fetching latest block height...")
    resp, err := http.Get(s.apiURL + "/status")
    if err != nil {
        return 0, fmt.Errorf("failed to get status: %v", err)
    }
    defer resp.Body.Close()
    
    var data struct {
        Synced bool  `json:"synced"`
        Height int64 `json:"height"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return 0, fmt.Errorf("failed to decode response: %v", err)
    }
    
    s.lastHeight = data.Height
    return data.Height, nil
}

func (s *SPVProc) ResolveDomain(name string) (map[string]interface{}, error) {
    log.Println("[SPVProc] Resolving domain:", name)
    reqBody, err := json.Marshal(map[string]string{"name": name})
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request body: %v", err)
    }
    resp, err := http.Post(s.apiURL+"/resolve", "application/json", bytes.NewBuffer(reqBody))
    if err != nil {
        return nil, fmt.Errorf("HTTP post failed: %v", err)
    }
    defer resp.Body.Close()
    
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }
    
    return result, nil
}