package main

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/prashantgupta24/automatic-mouse-mover/src/mousemover"
)

func main() {
	systray.Run(onReady, onExit)
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
				quit = mousemover.Start()
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
