package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/getlantern/systray"
	"github.com/go-vgo/robotgo"
	"github.com/kirsle/configdir"
	"github.com/prashantgupta24/automatic-mouse-mover/assets/icon"
	"github.com/prashantgupta24/automatic-mouse-mover/pkg/mousemover"
	log "github.com/sirupsen/logrus"
)

type AppSettings struct {
	Icon string `json:"icon"`
}

var configPath = configdir.LocalConfig("amm")
var configFile = filepath.Join(configPath, "settings.json")

func main() {
	systray.Run(onReady, onExit)
}

func setIcon(iconName string, configFile string) {
	switch {
	case iconName == "mouse":
		systray.SetIcon(icon.Data)
	case iconName == "cloud":
		systray.SetIcon(icon.CloudIcon)
	case iconName == "man":
		systray.SetIcon(icon.ManIcon)
	case iconName == "geometric":
		systray.SetIcon(icon.GeometricIcon)
	default:
		systray.SetIcon(icon.Data)
	}
	if configFile != "" {
		var settings AppSettings
		settings = AppSettings{iconName}
		fh, _ := os.Create(configFile)
		defer fh.Close()

		encoder := json.NewEncoder(fh)
		encoder.Encode(settings)
	}
}

func onReady() {
	go func() {
		err := configdir.MakePath(configPath)
		if err != nil {
			panic(err)
		}
		var settings AppSettings
		settings = AppSettings{"mouse"}

		if _, err = os.Stat(configFile); os.IsNotExist(err) {
			fh, err := os.Create(configFile)
			if err != nil {
				panic(err)
			}
			defer fh.Close()
			encoder := json.NewEncoder(fh)
			encoder.Encode(settings)

		} else {
			fh, err := os.Open(configFile)
			if err != nil {
				panic(err)
			}
			defer fh.Close()

			decoder := json.NewDecoder(fh)
			decoder.Decode(&settings)
		}
		setIcon(settings.Icon, "")

		about := systray.AddMenuItem("About AMM", "Information about the app")
		systray.AddSeparator()
		ammStart := systray.AddMenuItem("Start", "start the app")
		ammStop := systray.AddMenuItem("Stop", "stop the app")

		icons := systray.AddMenuItem("Icons", "icon of the app")
		mouse := icons.AddSubMenuItem("Mouse", "Mouse icon")
		mouse.SetIcon(icon.Data)
		cloud := icons.AddSubMenuItem("Cloud", "Cloud icon")
		cloud.SetIcon(icon.CloudIcon)
		man := icons.AddSubMenuItem("Man", "Man icon")
		man.SetIcon(icon.ManIcon)
		geometric := icons.AddSubMenuItem("Geometric", "Geometric")
		geometric.SetIcon(icon.GeometricIcon)

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
			case <-mouse.ClickedCh:
				setIcon("mouse", configFile)
			case <-cloud.ClickedCh:
				setIcon("cloud", configFile)
			case <-man.ClickedCh:
				setIcon("man", configFile)
			case <-geometric.ClickedCh:
				setIcon("geometric", configFile)
			case <-about.ClickedCh:
				log.Infof("Requesting about")
				robotgo.Alert("Automatic-mouse-mover app v1.2.0", "Developed by Prashant Gupta. \n\nMore info at: https://github.com/prashantgupta24/automatic-mouse-mover", "OK", "")
			}
		}

	}()
}

func onExit() {
	// clean up here
	log.Infof("Finished quitting")
}
