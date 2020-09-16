package input

import (
	"io"
	"os"
	"sync"

	"golang.org/x/crypto/ssh/terminal"
)

// makeRaw puts the reader r in raw mode if it is a TTY and it isn't
// already in raw mode. This is safe to call multiple times.
func makeRaw(r io.Reader) {
	rawLock.Lock()
	defer rawLock.Unlock()

	f, ok := r.(*os.File)
	if !ok {
		return
	}

	fd := f.Fd()
	if !terminal.IsTerminal(int(fd)) {
		return
	}

	st, err := terminal.MakeRaw(int(fd))
	if err != nil {
		return
	}

	rawTracker[fd] = st
}

func unRaw(r io.Reader) {
	rawLock.Lock()
	defer rawLock.Unlock()

	f, ok := r.(*os.File)
	if !ok {
		return
	}

	fd := f.Fd()
	st, ok := rawTracker[fd]
	if !ok {
		return
	}

	terminal.Restore(int(fd), st)
	delete(rawTracker, fd)
}

var rawLock sync.Mutex
var rawTracker = map[uintptr]*terminal.State{}
