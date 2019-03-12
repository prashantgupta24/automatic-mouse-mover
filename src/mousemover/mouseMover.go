package mousemover

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/src/activity"
	"github.com/prashantgupta24/activity-tracker/src/mouse"
)

func Start() (quit chan struct{}) {
	quit = make(chan struct{})
	activityCh, quitActivityDetector := activity.StartTracker(5)

	go func() {
		movePixel := 10
		for {
			select {
			case isActivity := <-activityCh:
				if !isActivity {
					currentMousePos := mouse.GetPosition()
					fmt.Println("moving mouse")
					nextMouseMov := &mouse.Position{
						MouseX: currentMousePos.MouseX + movePixel,
						MouseY: currentMousePos.MouseY + movePixel,
					}
					robotgo.Move(nextMouseMov.MouseX, nextMouseMov.MouseY)
					movePixel *= -1
				}
			case <-quit:
				fmt.Println("stopped main app")
				quitActivityDetector <- struct{}{}
				return
			}
		}
	}()
	return quit
}
