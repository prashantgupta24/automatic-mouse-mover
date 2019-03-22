package mousemover

import (
	"fmt"
	"log"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/pkg/tracker"
)

var instance *mouseMover

type mouseMover struct {
	quit      chan struct{}
	isRunning bool
}

const (
	timeout = 100 //ms
)

func (m *mouseMover) Start() {
	m.quit = make(chan struct{})

	frequency := 5 //value always in seconds
	activityTracker := &tracker.Instance{
		Frequency: frequency,
		//LogLevel:  "debug", //if we want verbose logging
	}

	heartbeatCh := activityTracker.Start()

	go func(m *mouseMover) {
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

func (m *mouseMover) Quit() {
	//making it idempotent
	if m != nil && m.isRunning {
		m.quit <- struct{}{}
	}
}

//GetInstance gets the singleton instance for mouse mover app
func GetInstance() *mouseMover {
	if instance == nil {
		instance = &mouseMover{}
	}
	return instance
}
