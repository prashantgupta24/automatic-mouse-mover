package mousemover

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/pkg/tracker"
)

var instance *MouseMover

//MouseMover is the main struct for the app
type MouseMover struct {
	quit          chan struct{}
	mutex         sync.RWMutex
	runningStatus bool
}

const (
	timeout = 100 //ms
)

//Start the main app
func (m *MouseMover) Start() {
	if m.isRunning() {
		return
	}
	m.quit = make(chan struct{})

	frequency := 60 //value always in seconds
	activityTracker := &tracker.Instance{
		Frequency: frequency,
		//LogLevel:  "debug", //if we want verbose logging
	}

	heartbeatCh := activityTracker.Start()

	go func(m *MouseMover) {
		m.updateRunningStatus(true)
		movePixel := 10
		for {
			select {
			case heartbeat := <-heartbeatCh:
				if !heartbeat.WasAnyActivity {
					mouseMoveSuccessCh := make(chan bool)
					go moveAndCheck(movePixel, mouseMoveSuccessCh)
					select {
					case wasMouseMoveSuccess := <-mouseMoveSuccessCh:
						if wasMouseMoveSuccess {
							log.Infof("moved mouse at : %v\n\n", time.Now())
							movePixel *= -1
						} else {
							msg := "Mouse pointer cannot be moved. See README for details."
							log.Errorf(msg)
							robotgo.ShowAlert("Error with Automatic Mouse Mover", msg)
						}
					case <-time.After(timeout * time.Millisecond):
						//timeout, do nothing
						log.Errorf("timeout happened after %vms while trying to move mouse", timeout)
					}

				}
			case <-m.quit:
				log.Infof("stopping mouse mover")
				m.updateRunningStatus(false)
				activityTracker.Quit()
				return
			}
		}
	}(m)
}

func (m *MouseMover) isRunning() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.runningStatus
}
func (m *MouseMover) updateRunningStatus(isRunning bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.runningStatus = isRunning
}

func moveAndCheck(movePixel int, mouseMoveSuccessCh chan bool) {
	currentX, currentY := robotgo.GetMousePos()
	moveToX := currentX + movePixel
	moveToY := currentY + movePixel
	robotgo.MoveMouse(moveToX, moveToY)

	//check if mouse moved. Sometimes mac users need to give
	//extra permission for controlling the mouse
	movedX, movedY := robotgo.GetMousePos()
	if movedX == currentX && movedY == currentY {
		mouseMoveSuccessCh <- false
	} else {
		mouseMoveSuccessCh <- true
	}
}

//Quit the app
func (m *MouseMover) Quit() {
	//making it idempotent
	if m != nil && m.isRunning() {
		m.quit <- struct{}{}
	}
}

//GetInstance gets the singleton instance for mouse mover app
func GetInstance() *MouseMover {
	if instance == nil {
		instance = &MouseMover{}
	}
	return instance
}
