package main

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/prashantgupta24/automatic-mouse-mover/pkg/mousemover"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	go func() {
		systray.SetIcon(icon.Data)
		systray.SetTitle("AMM")
		ammStart := systray.AddMenuItem("Start", "start the app")
		ammStop := systray.AddMenuItem("Stop", "stop the app")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
		// Sets the icon of a menu item. Only available on Mac.
		//mQuit.SetIcon(icon.Data)
		mouseMover := mousemover.GetInstance()
		for {
			select {
			case <-ammStart.ClickedCh:
				fmt.Println("starting the app")
				mouseMover.Start()
				//notify.SendMessage("starting the app")

			case <-ammStop.ClickedCh:
				fmt.Println("stopping the app")
				mouseMover.Quit()

			case <-mQuit.ClickedCh:
				fmt.Println("Requesting quit")
				mouseMover.Quit()
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
