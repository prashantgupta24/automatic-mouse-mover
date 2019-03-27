package mousemover

import (
	"fmt"
	"log"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/pkg/tracker"
)

var instance *MouseMover

//MouseMover is the main struct for the app
type MouseMover struct {
	quit      chan struct{}
	isRunning bool
}

const (
	timeout = 100 //ms
)

//Start the main app
func (m *MouseMover) Start() {
	m.quit = make(chan struct{})

	frequency := 5 //value always in seconds
	activityTracker := &tracker.Instance{
		Frequency: frequency,
		//LogLevel:  "debug", //if we want verbose logging
	}

	heartbeatCh := activityTracker.Start()

	go func(m *MouseMover) {
		m.isRunning = true
		movePixel := 10
		for {
			select {
			case heartbeat := <-heartbeatCh:
				if !heartbeat.IsActivity {
					commCh := make(chan bool)
					go moveMouse(movePixel, commCh)
					select {
					case wasMouseMoveSuccess := <-commCh:
						if wasMouseMoveSuccess {
							fmt.Printf("\nmoving mouse at : %v\n\n", time.Now())
							movePixel *= -1
						}
					case <-time.After(timeout * time.Millisecond):
						//timeout, do nothing
						log.Printf("timeout happened after %vms while trying to move mouse", timeout)
					}

				}
			case <-m.quit:
				fmt.Println("stopping mouse mover")
				m.isRunning = false
				activityTracker.Quit()
				return
			}
		}
	}(m)
}

func moveMouse(movePixel int, commCh chan bool) {
	currentX, currentY := robotgo.GetMousePos()
	moveToX := currentX + movePixel
	moveToY := currentY + movePixel
	robotgo.MoveMouse(moveToX, moveToY)
	commCh <- true
}

//Quit the app
func (m *MouseMover) Quit() {
	//making it idempotent
	if m != nil && m.isRunning {
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
