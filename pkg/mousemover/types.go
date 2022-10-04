package mousemover

import (
	"os"
	"sync"
	"time"
)

//MouseMover is the main struct for the app
type MouseMover struct {
	quit    chan struct{}
	logFile *os.File
	state   *state
}

//state manages the internal working of the app
type state struct {
	mutex              sync.RWMutex
	isAppRunning       bool
	isSysSleeping      bool
	lastMouseMovedTime time.Time
	lastErrorTime      time.Time
	didNotMoveCount    int
	override           *override
}

//only needed for tests
type override struct {
	valueToReturn bool
}
