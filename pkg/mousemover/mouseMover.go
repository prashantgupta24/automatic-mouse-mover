package mousemover

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	"github.com/prashantgupta24/activity-tracker/pkg/tracker"
)

var instance *MouseMover

//MouseMover is the main struct for the app
type MouseMover struct {
	quit          chan struct{}
	mutex         sync.RWMutex
	runningStatus bool
	logFile       *os.File
}

const (
	timeout     = 100 //ms
	logDir      = "log"
	logFileName = "logFile-amm-2"
)

//Start the main app
func (m *MouseMover) Start() {
	if m.isRunning() {
		return
	}
	m.quit = make(chan struct{})

	heartbeatInterval := 60 //value always in seconds
	workerInterval := 10

	activityTracker := &tracker.Instance{
		HeartbeatInterval: heartbeatInterval,
		WorkerInterval:    workerInterval,
		//LogLevel:  "debug", //if we want verbose logging
	}

	heartbeatCh := activityTracker.Start()
	m.run(heartbeatCh, activityTracker)
}

func (m *MouseMover) run(heartbeatCh chan *tracker.Heartbeat, activityTracker *tracker.Instance) {
	go func() {
		if m.isRunning() {
			return
		}
		logger := getLogger(m, false) //set writeToFile=true only for debugging
		m.updateRunningStatus(true)
		movePixel := 10
		var lastMoved time.Time
		isSystemSleeping := false
		didNotMoveTimes := 0
		for {
			select {
			case heartbeat := <-heartbeatCh:
				if !heartbeat.WasAnyActivity {
					if isSystemSleeping {
						logger.Infof("system sleeping")
						continue
					}
					mouseMoveSuccessCh := make(chan bool)
					go moveAndCheck(movePixel, mouseMoveSuccessCh)
					select {
					case wasMouseMoveSuccess := <-mouseMoveSuccessCh:
						if wasMouseMoveSuccess {
							lastMoved = time.Now()
							logger.Infof("moved mouse at : %v\n\n", lastMoved)
							movePixel *= -1
							didNotMoveTimes = 0
						} else {
							didNotMoveTimes++
							msg := fmt.Sprintf("Mouse pointer cannot be moved at %v. Last moved at %v. Happened %v times. See README for details.",
								time.Now(), lastMoved, didNotMoveTimes)
							logger.Errorf(msg)
							if didNotMoveTimes >= 3 {
								go func() {
									robotgo.ShowAlert("Error with Automatic Mouse Mover", msg)
								}()
							}
						}
					case <-time.After(timeout * time.Millisecond):
						//timeout, do nothing
						logger.Errorf("timeout happened after %vms while trying to move mouse", timeout)
					}
				} else {
					logger.Infof("activity detected in the last %v seconds.", int(activityTracker.HeartbeatInterval))
					logger.Infof("Activity type:\n")
					for activityType, times := range heartbeat.ActivityMap {
						logger.Infof("activityType : %v times: %v\n", activityType, len(times))
						if activityType == activity.MachineSleep {
							isSystemSleeping = true
						} else if activityType == activity.MachineWake {
							isSystemSleeping = false
						}
					}
					logger.Infof("\n\n\n")
				}
			case <-m.quit:
				logger.Infof("stopping mouse mover")
				m.updateRunningStatus(false)
				activityTracker.Quit()
				return
			}
		}
	}()
}

//Quit the app
func (m *MouseMover) Quit() {
	//making it idempotent
	if m != nil && m.isRunning() {
		m.quit <- struct{}{}
	}
	if m.logFile != nil {
		m.logFile.Close()
	}
}

//GetInstance gets the singleton instance for mouse mover app
func GetInstance() *MouseMover {
	if instance == nil {
		instance = &MouseMover{}
	}
	return instance
}
