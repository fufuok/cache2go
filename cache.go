/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2012, Radu Ioan Fericean
 *                   2013-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package cache2go

import (
	"sync"
)

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex
)

// Cache returns the existing cache table with given name or creates a new one
// if the table does not exist yet.
func Cache(table string) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		mutex.Lock()
		t, ok = cache[table]
		// Double check whether the table exists or not.
		if !ok {
			t = &CacheTable{
				name:  table,
				items: make(map[interface{}]*CacheItem),
			}
			cache[table] = t
		}
		mutex.Unlock()
	}

	return t
}

// Count returns how many cache table.
func Count() int {
	mutex.RLock()
	defer mutex.RUnlock()
	return len(cache)
}

// Delete an cache table.
func Delete(table string) {
	mutex.Lock()
	defer mutex.Unlock()

	if t, ok := cache[table]; ok {
		t.Lock()
		defer t.Unlock()

		t.items = nil
		t.cleanupInterval = 0
		if t.cleanupTimer != nil {
			t.cleanupTimer.Stop()
		}

		delete(cache, table)
	}
}
