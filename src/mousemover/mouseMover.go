package mousemover

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/src/activity"
	"github.com/prashantgupta24/activity-tracker/src/mouse"
)

func Start() (quit chan struct{}) {
	quit = make(chan struct{})

	activityTracker := &activity.ActivityTracker{
		TimeToCheck: 5,
	}

	heartbeatCh, quitActivityTracker := activityTracker.Start()

	go func() {
		movePixel := 10
		for {
			select {
			case heartbeat := <-heartbeatCh:
				if !heartbeat.IsActivity {
					currentMousePos := mouse.GetPosition()
					fmt.Println("moving mouse at : ", time.Now())
					nextMouseMov := &mouse.Position{
						MouseX: currentMousePos.MouseX + movePixel,
						MouseY: currentMousePos.MouseY + movePixel,
					}
					robotgo.MoveMouseSmooth(nextMouseMov.MouseX, nextMouseMov.MouseY)
					movePixel *= -1
				}
			case <-quit:
				fmt.Println("stopping mouse mover")
				quitActivityTracker <- struct{}{}
				return
			}
		}
	}()
	return quit
}
