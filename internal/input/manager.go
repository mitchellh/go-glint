package input

import (
	"io"
	"sync"
	"sync/atomic"
)

// Func is the callback type for any input.
type Func func([]byte, error)

// Manager is responsible for managing the input loop for a single
// input reader.
type Manager struct {
	mu     *sync.Mutex
	cond   *sync.Cond
	cbs    map[interface{}]Func
	closed bool
	count  int32
	r      io.Reader
}

// For will return the input manager for the given key. If it doesn't exist
// it will be created. Close should be called for the key when done to clean
// up any resources.
func For(key interface{}, r io.Reader) *Manager {
	managerLock.Lock()
	defer managerLock.Unlock()

	// If we have one return it.
	if m, ok := managers[key]; ok {
		atomic.AddInt32(&m.count, 1)
		return m
	}

	// If we don't have one then create it.
	var mu sync.Mutex
	m := &Manager{
		mu:     &mu,
		cond:   sync.NewCond(&mu),
		cbs:    map[interface{}]Func{},
		closed: true,
		count:  1,
		r:      r,
	}

	// Start the watcher if we have a reader. If the reader is nil we
	// always send an EOF.
	if r != nil {
		// We aren't closed yet
		m.closed = false

		// Make this reader raw if we can
		makeRaw(r)

		// Start our listener
		go m.run(r)
	}

	// Store it
	managers[key] = m

	return m
}

func Close(key interface{}) {
	managerLock.Lock()
	defer managerLock.Unlock()

	// If we don't track this then we do nothing.
	m, ok := managers[key]
	if !ok {
		return
	}

	// Decrement the count. If we aren't the last one, do nothing.
	new := atomic.AddInt32(&m.count, -1)
	if new > 0 {
		return
	}

	// If we have no reader, do nothing
	if m.r == nil {
		return
	}

	// Clean up the reader
	unRaw(m.r)
}

// AddCallback adds the callback to get called for any input. If a callback
// already exists for key, nothing will be done.
func (m *Manager) AddCallback(key interface{}, cb Func) {
	if cb == nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// If we don't have this callback, add it and broadcast.
	if _, ok := m.cbs[key]; !ok {
		m.cbs[key] = cb
		m.cond.Broadcast()

		// If we're already closed, send an EOF
		if m.closed {
			cb(nil, io.EOF)
		}
	}
}

// DeleteCallback removes the callback for the given key. This will do
// nothing if the key doesn't already have a callback associated.
func (m *Manager) DeleteCallback(key interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.cbs, key)
}

func (m *Manager) run(r io.Reader) {
	var buf [128]byte
	for {
		n, err := r.Read(buf[:])

		// Wrap the logic below in a func so we can use a defer.
		exit := false
		func() {
			// Call the callbacks. If we have no callbacks then we wait until
			// a callback is added and block reading the input until then.
			m.mu.Lock()
			defer m.mu.Unlock()

			for len(m.cbs) == 0 {
				m.cond.Wait()
			}

			for _, f := range m.cbs {
				f(buf[:n], err)
			}

			// If we got an EOF, then mark that we're done
			if err == io.EOF {
				m.closed = true
			}

			exit = m.closed
		}()

		// If we got an EOF then we're done
		if exit {
			return
		}
	}
}

var managerLock sync.Mutex
var managers = map[interface{}]*Manager{}
