package middleware

import (
	"fmt"
	"github.com/gofiber/storage/memory"
	"strings"
	"sync"
	"time"
)

// ✅ Exported cache variable
var Cache = memory.New()
var cacheKeys = make(map[string]bool) // ✅ Track all keys manually
var lock sync.Mutex                   // ✅ Protect access to cacheKeys slice

// CacheSet stores data with expiration
func CacheSet(key string, value string, expiration time.Duration) {
	lock.Lock()
	defer lock.Unlock()

	// ✅ Try storing value in cache
	err := Cache.Set(key, []byte(value), expiration)
	if err != nil {
		fmt.Println("❌ CacheSet Failed: Unable to store value")
	}

	// ✅ Verify storage immediately
	storedValue, err := Cache.Get(key)
	if err != nil {
		fmt.Println("❌ CacheSet Verification Failed: Value not stored")
	} else {
		// ✅ Add key to tracked keys
		cacheKeys[key] = true
		fmt.Println("✅ CacheSet Success: Value stored ->", string(storedValue))
	}
}

// CacheExists checks if a key exists
func CacheExists(key string) bool {
	_, err := Cache.Get(key)
	return err == nil
}

// ✅ Fix CacheListKeys to return tracked keys
func CacheListKeys(prefix string) []string {
	lock.Lock()
	defer lock.Unlock()

	var keys []string
	for key := range cacheKeys {
		if strings.HasPrefix(key, prefix) {
			keys = append(keys, key)
		}
	}
	return keys
}

// ✅ CacheGet retrieves data from cache
func CacheGet(key string) (string, bool) {
	lock.Lock()
	defer lock.Unlock()

	val, err := Cache.Get(key)
	if err != nil {
		return "", false
	}

	return string(val), true
}

// ✅ Fix CacheDelete to remove key from tracking list
func CacheDelete(key string) {
	lock.Lock()
	defer lock.Unlock()

	delete(cacheKeys, key)
	Cache.Delete(key)
}
