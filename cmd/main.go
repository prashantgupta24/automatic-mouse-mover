package main

import (
	"github.com/go-vgo/robotgo"
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
		about := systray.AddMenuItem("About AMM", "Information about the app")
		systray.AddSeparator()
		ammStart := systray.AddMenuItem("Start", "start the app")
		ammStop := systray.AddMenuItem("Stop", "stop the app")
		ammStop.Disable()
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
		// Sets the icon of a menu item. Only available on Mac.
		//mQuit.SetIcon(icon.Data)
		mouseMover := mousemover.GetInstance()
		mouseMover.Start()
		ammStart.Disable()
		ammStop.Enable()
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

			case <-about.ClickedCh:
				log.Infof("Requesting about")
				robotgo.ShowAlert("Automatic-mouse-mover app v1.1.0", "Developed by Prashant Gupta. \n\nMore info at: https://github.com/prashantgupta24/automatic-mouse-mover")
			}
		}

	}()
}

func onExit() {
	// clean up here
	log.Infof("Finished quitting")
}
