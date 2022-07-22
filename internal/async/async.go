// Package async contains helpers for asynchronous operations.
package async

import "context"

// Mutex is a context aware mutex implementation.
type Mutex struct {
	ch chan struct{}
}

// Lock acquires a lock for a protected resource or code path.
func (mu *Mutex) Lock(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	case mu.ch <- struct{}{}:
		return true
	}
}

// Unlock releases a lock for a protected resource or code path.
func (mu *Mutex) Unlock() {
	<-mu.ch
}
