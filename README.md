## Presenting the minimalistic Automatic-Mouse-Mover(AMM) app!

![GitHub release](https://img.shields.io/github/release/prashantgupta24/automatic-mouse-mover.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/prashantgupta24/automatic-mouse-mover)](https://goreportcard.com/report/github.com/prashantgupta24/automatic-mouse-mover)

Ever felt the need to keep your machine awake without having to resort to the age-old methods of installing an app that you don't trust or playing a video? **Well, not anymore!**

Introducing the simplest app on the market that has the sole purpose of moving your mouse pointer at regular intervals so that your machine never sleeps! And best of all, it works **ONLY** when you are not working, so be rest assured that the mouse won't start moving on its own without the machine actually being idle.

## Demo

You just click on `Start`, and AMM will take care of moving your mouse whenever it feels that the system has been left idle for a long time. It's as simple as this. 

![](https://github.com/prashantgupta24/automatic-mouse-mover/tree/master/resources/amm-demo.gif)

## How to install

Make sure you have `go` installed. Once that is done, clone this repo and run `Make`, it should create the `amm.app` and open the folder where it was built for you. You just have to drag and drop it to the `Applications` folder on your mac. 

Double click on the app, and the cute `mouse` should appear on your taskbar on top of your screen. Once you click on `Start`, you might encounter an initial `request` which I've discussed in the next section. If not, then you are all set! 

## Granting access for moving the mouse cursor

While starting the app, you might see a message like the one below or an error stating `Mouse pointer cannot be moved.`.

![](https://github.com/prashantgupta24/automatic-mouse-mover/tree/master/resources/request.jpg)

Don't worry, it's nothing sinister, it's just that Mac doesn't allow apps to gain accessibility to the computer by default (even standard apps like Automator, Firefox etc. who might want to access some features need to go through the same process).

In order to resolve this error you need to:

> Go to System Preferences -> Security & Privacy -> Privacy -> Accessibility and allow the `amm` app to gain access.

## How it works

AMM uses [Activity tracker](https://github.com/prashantgupta24/activity-tracker) to track various changes to your system. All code is public and open-sourced so no worrying if there's nefarious intention involved or not.
