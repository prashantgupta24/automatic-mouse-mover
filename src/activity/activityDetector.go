package activity

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/mousemover/src/mouse"
)

const (
	timeToCheck = 5
)

func StartDetector() (activity chan bool, quit chan struct{}) {

	comm, quitMouseClickHandler := isMouseClicked()

	activity = make(chan bool)
	quit = make(chan struct{})
	ticker := time.NewTicker(time.Second * timeToCheck)

	go func() {
		isIdle := true
		lastMousePos := mouse.GetPosition()
		for {
			select {
			case <-ticker.C:
				//fmt.Println("ticked : ", isIdle)
				currentMousePos := mouse.GetPosition()
				if isIdle && isPointerIdle(currentMousePos, lastMousePos) {
					fmt.Printf("no activity detected in the last %v seconds ...\n", timeToCheck)
					activity <- false
				} else {
					fmt.Printf("activity detected in the last %v seconds ...\n", timeToCheck)
					activity <- true
					lastMousePos = currentMousePos
				}
				isIdle = true
			case <-comm:
				isIdle = false
				//fmt.Println("val received: ", isIdle)
			case <-quit:
				fmt.Println("stopped activity tracker")
				quitMouseClickHandler <- struct{}{}
				robotgo.StopEvent()
				return
			}
		}
	}()

	return activity, quit
}

func isPointerIdle(currentMousePos, lastMousePos *mouse.Position) bool {
	// fmt.Println("current : ", currentMousePos)
	// fmt.Println("last : ", lastMousePos)
	if currentMousePos.MouseX == lastMousePos.MouseX &&
		currentMousePos.MouseY == lastMousePos.MouseY {
		return true
	}
	return false
}

func isMouseClicked() (clickComm, quit chan struct{}) {
	ticker := time.NewTicker(time.Second * timeToCheck)
	clickComm = make(chan struct{})
	quit = make(chan struct{})
	registrationFree := make(chan struct{})
	go func() {
		isRunning := false
		for {
			select {
			case <-ticker.C:
				fmt.Println("ticked mouse click at : ", time.Now())

				if !isRunning {
					isRunning = true
					go func(registrationFree chan struct{}) {
						fmt.Printf("adding reg\n\n\n")
						mleft := robotgo.AddEvent("mleft")
						if mleft == 0 {
							fmt.Println("mleft clicked")
							clickComm <- struct{}{}
							registrationFree <- struct{}{}
							return
						}
					}(registrationFree)
				}

				select {
				case _, ok := <-registrationFree:
					if ok {
						fmt.Println("registration free")
						isRunning = false
					} else {
						fmt.Println("Channel closed!")
					}
				default:
					fmt.Println("registration is busy")
					isRunning = true
				}

			case <-quit:
				fmt.Println("stopped click handler")
				close(clickComm)
				return
			}
		}
	}()
	return clickComm, quit
}
