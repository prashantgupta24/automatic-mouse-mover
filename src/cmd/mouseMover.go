package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/go-vgo/robotgo"
)

const (
	timeToCheck = 5
)

type mousePos struct {
	mouseX int
	mouseY int
}

func moveMouse(comm, quit chan struct{}) {
	quitMouseClickHandler := make(chan struct{})
	go isMouseClicked1(comm, quitMouseClickHandler)
	ticker := time.NewTicker(time.Second * timeToCheck)
	isIdle := true
	movePixel := 10
	lastMousePos := getMousePos()
	for {
		select {
		case <-ticker.C:
			//fmt.Println("ticked : ", isIdle)
			currentMousePos := getMousePos()
			if isIdle && isPointerIdle(currentMousePos, lastMousePos) {
				fmt.Printf("no activity detected in the last %v seconds. moving mouse ...\n", timeToCheck)
				nextMouseMov := &mousePos{
					mouseX: currentMousePos.mouseX + movePixel,
					mouseY: currentMousePos.mouseY + movePixel,
				}
				robotgo.Move(nextMouseMov.mouseX, nextMouseMov.mouseY)
				lastMousePos = nextMouseMov
			} else {
				lastMousePos = currentMousePos
			}
			isIdle = true
			movePixel *= -1
		case <-comm:
			isIdle = false
			//fmt.Println("val received: ", isIdle)
		case <-quit:
			fmt.Println("stopped main app")
			quitMouseClickHandler <- struct{}{}
			robotgo.StopEvent()
			return
		}
	}
}

func isPointerIdle(currentMousePos, lastMousePos *mousePos) bool {
	// fmt.Println("current : ", currentMousePos)
	// fmt.Println("last : ", lastMousePos)

	if currentMousePos.mouseX == lastMousePos.mouseX &&
		currentMousePos.mouseY == lastMousePos.mouseY {
		return true
	}
	return false
}

func getMousePos() *mousePos {
	x, y := robotgo.GetMousePos()
	return &mousePos{
		mouseX: x,
		mouseY: y,
	}
}

// func isMouseClicked1(comm, quit chan struct{}) {
// 	var wg sync.WaitGroup

// 	go func() {}(&wg)

// 	go func() {}(&wg)
// }

// func shoudStop(quit chan struct{}) bool {
// 	fmt.Println("calling should stop")
// 	select {
// 	case <-quit:
// 		fmt.Println("should stop now")
// 		return true
// 	default:
// 		return false
// 	}
// }

func isMouseClicked1(comm, quit chan struct{}) {
	ticker := time.NewTicker(time.Second * timeToCheck)
	ch := make(chan struct{})
	go func() {
		isRunning := false
		for {
			select {
			case <-ticker.C:
				fmt.Println("ticked mouse click at : ", time.Now())

				if !isRunning {
					isRunning = true
					go func(ch chan struct{}) {
						fmt.Printf("adding reg\n\n\n")
						mleft := robotgo.AddEvent("mleft")
						if mleft == 0 {
							fmt.Println("mleft clicked")
							comm <- struct{}{}
							ch <- struct{}{}
							return
						}
					}(ch)
				}

				select {
				case _, ok := <-ch:
					if ok {
						fmt.Println("channel free")
						isRunning = false
					} else {
						fmt.Println("Channel closed!")
					}
				default:
					fmt.Println("function is busy")
					isRunning = true
				}

				// wg.Wait()
				// wg.Add(1)
				// go func(wg *sync.WaitGroup) {
				// 	fmt.Printf("adding reg\n\n\n")
				// 	defer wg.Done()
				// 	mleft := robotgo.AddEvent("mleft")
				// 	if mleft == 0 {
				// 		fmt.Println("mleft clicked")
				// 		comm <- struct{}{}
				// 		return
				// 	}
				// }(&wg)
			case <-quit:
				fmt.Println("should stop now")
				close(comm)
				return
			}
		}
	}()
}

func shoudStop(quit chan struct{}) bool {
	fmt.Println("calling should stop")
	select {
	case <-quit:
		fmt.Println("should stop now")
		return true
	default:
		return false
	}
}

func isMouseClicked(comm, quit chan struct{}) {
	//var wg sync.WaitGroup
	go func() {
		for {
			if !shoudStop(quit) {
				fmt.Println("adding reg")
				mleft := robotgo.AddEvent("mleft")
				if mleft == 0 {
					fmt.Println("mleft clicked")
					comm <- struct{}{}
				}
			} else {
				fmt.Println("closing mouse clicked function")
				close(comm)
				return
			}
		}
		// select {
		// case <-quit:
		// 	fmt.Println("asked to quit")
		// 	isRunning = false
		// 	close(comm)
		// default:

		// 	fmt.Println("adding reg")
		// 	mleft := robotgo.AddEvent("mleft")
		// 	if mleft == 0 {
		// 		if isRunning {
		// 			fmt.Println("mleft clicked")
		// 			comm <- struct{}{}
		// 		} else {
		// 			return
		// 		}

		// 		//time.Sleep(1000 * time.Millisecond)
		// 	}
		// 	//}
		// }

	}()

	//waiting for next robot go release in which
	// we can handle multiple event registrations
	// go func() {
	// 	for {
	// 		mright := robotgo.AddEvent("mright")
	// 		if mright == 0 {
	// 			fmt.Println("mright clicked")
	// 			time.Sleep(100 * time.Millisecond)
	// 			comm <- struct{}{}
	// 		}
	// 	}
	// }()
}

func main() {
	systray.Run(onReady, onExit)

}
func startApp() chan struct{} {
	comm := make(chan struct{})
	quit := make(chan struct{})
	go moveMouse(comm, quit)
	return quit
}

func onReady() {
	go func() {
		systray.SetIcon(icon.Data)
		systray.SetTitle("AMM")
		ammStart := systray.AddMenuItem("Start", "start the app")
		ammPause := systray.AddMenuItem("Pause", "pause the app")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
		// Sets the icon of a menu item. Only available on Mac.
		//mQuit.SetIcon(icon.Data)
		var quit chan struct{}
		for {
			select {
			case <-ammStart.ClickedCh:
				fmt.Println("starting the app")
				quit = startApp()
				//notify.SendMessage("starting the app")

			case <-ammPause.ClickedCh:
				fmt.Println("pausing the app")
				if quit != nil {
					quit <- struct{}{}
				} else {
					fmt.Println("app is not started")
				}

			case <-mQuit.ClickedCh:
				fmt.Println("Requesting quit")
				systray.Quit()
				fmt.Println("Finished quitting")
				return
			}
		}

	}()
}

func onExit() {
	// clean up here
	fmt.Println("exiting")

}
