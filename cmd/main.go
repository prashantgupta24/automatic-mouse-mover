package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/getlantern/systray"
	"github.com/prashantgupta24/automatic-mouse-mover/assets/icon"
	"github.com/prashantgupta24/automatic-mouse-mover/pkg/mousemover"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	go func() {
		systray.SetIcon(icon.Data)
		ammStart := systray.AddMenuItem("Start", "start the app")
		ammStop := systray.AddMenuItem("Stop", "stop the app")
		ammStop.Disable()
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
		// Sets the icon of a menu item. Only available on Mac.
		//mQuit.SetIcon(icon.Data)
		mouseMover := mousemover.GetInstance()
		for {
			select {
			case <-ammStart.ClickedCh:
				log.Infof("starting the app")
				mouseMover.Start()
				ammStart.Disable()
				ammStop.Enable()
				//notify.SendMessage("starting the app")

			case <-ammStop.ClickedCh:
				log.Infof("stopping the app")
				ammStart.Enable()
				ammStop.Disable()
				mouseMover.Quit()

			case <-mQuit.ClickedCh:
				log.Infof("Requesting quit")
				mouseMover.Quit()
				systray.Quit()
				return
			}
		}

	}()
}

func onExit() {
	// clean up here
	log.Infof("Finished quitting")
}
