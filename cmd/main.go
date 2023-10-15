package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/getlantern/systray"
	"github.com/go-vgo/robotgo"
	"github.com/kirsle/configdir"
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

func getIcon(iconName string, active bool) []byte {
	if iconName != "mouse" && iconName != "cloud" && iconName != "geometric" && iconName != "man" {
		iconName = "mouse"
	}
	var b []byte
	var err error
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)

	if _, err := os.Stat(exPath + "/../Resources/assets/icon"); os.IsNotExist(err) {
		b, err = os.ReadFile(exPath + "/../assets/icon/" + iconName + ".png")
	} else {
		b, err = os.ReadFile(exPath + "/../Resources/assets/icon/" + iconName + ".png")
	}
	if err != nil {
		panic(err)
	}
	if active {
		img, err := png.Decode(bytes.NewReader(b))
		if err != nil {
			log.Fatalln(err)
		}
		var dimg *image.RGBA = image.NewRGBA(img.Bounds())
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				if a != 0 {
					dimg.Set(x, y, color.RGBA{30, 144, 255, 255})
				} else {
					dimg.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
				}
			}
		}
		var c bytes.Buffer
		png.Encode(&c, dimg)
		return c.Bytes()
	}
	return b

}

func setIcon(iconName string, configFile string, settings *AppSettings, active ...bool) {
	systray.SetIcon(getIcon(iconName, len(active) != 0 && active[0]))
	if configFile != "" {

		*settings = AppSettings{iconName}
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

		about := systray.AddMenuItem("About AMM", "Information about the app")
		systray.AddSeparator()
		ammStart := systray.AddMenuItem("Start", "start the app")
		ammStop := systray.AddMenuItem("Stop", "stop the app")

		icons := systray.AddMenuItem("Icons", "icon of the app")
		mouse := icons.AddSubMenuItem("Mouse", "Mouse icon")

		mouse.SetIcon(getIcon("mouse", false))
		cloud := icons.AddSubMenuItem("Cloud", "Cloud icon")
		cloud.SetIcon(getIcon("cloud", false))
		man := icons.AddSubMenuItem("Man", "Man icon")
		man.SetIcon(getIcon("man", false))
		geometric := icons.AddSubMenuItem("Geometric", "Geometric")
		geometric.SetIcon(getIcon("geometric", false))

		ammStop.Disable()
		setIcon(settings.Icon, "", &settings, true)
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
				setIcon(settings.Icon, configFile, &settings, true)

			case <-ammStop.ClickedCh:
				log.Infof("stopping the app")
				ammStart.Enable()
				ammStop.Disable()
				mouseMover.Quit()
				setIcon(settings.Icon, configFile, &settings, false)

			case <-mQuit.ClickedCh:
				log.Infof("Requesting quit")
				mouseMover.Quit()
				systray.Quit()
				return
			case <-mouse.ClickedCh:
				setIcon("mouse", configFile, &settings, ammStart.Disabled())
			case <-cloud.ClickedCh:
				setIcon("cloud", configFile, &settings, ammStart.Disabled())
			case <-man.ClickedCh:
				setIcon("man", configFile, &settings, ammStart.Disabled())
			case <-geometric.ClickedCh:
				setIcon("geometric", configFile, &settings, ammStart.Disabled())
			case <-about.ClickedCh:
				log.Infof("Requesting about")
				robotgo.Alert("Automatic-mouse-mover app v1.3.1", "Developed by Prashant Gupta. \n\nMore info at: https://github.com/prashantgupta24/automatic-mouse-mover", "OK", "")
			}
		}

	}()
}

func onExit() {
	// clean up here
	log.Infof("Finished quitting")
}
