**Note: Known issues with MacOS Sonoma, will work on a fix asap**

# Presenting the minimalistic Automatic-Mouse-Mover(AMM) app!

[![version][version-badge]][releases] [![Go Report Card](https://goreportcard.com/badge/github.com/prashantgupta24/automatic-mouse-mover)](https://goreportcard.com/report/github.com/prashantgupta24/automatic-mouse-mover) [![godoc-badge][godoc-badge]][godoc-link] [![codecov](https://codecov.io/gh/prashantgupta24/automatic-mouse-mover/branch/master/graph/badge.svg)](https://codecov.io/gh/prashantgupta24/automatic-mouse-mover)

Ever felt the need to keep your machine awake without actually having to move the mouse pointer manually at regular intervals? **Well, not anymore!**

Introducing the simplest app that has the sole purpose of **moving your mouse pointer at regular intervals so that your machine is kept awake!** And best of all, it works **ONLY** when you are not working, so be rest assured that the mouse won't start moving on its own without the machine actually being idle.

**Table of contents**

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [How it's different from other apps](#how-its-different-from-other-apps)
- [How I use it](#how-i-use-it)
  - [Work from home benefit](#work-from-home-benefit)
- [Demo](#demo)
- [How to install](#how-to-install)
  - [Install from binary](#install-from-binary)
  - [Install from source](#install-from-source)
- [Granting access for moving the mouse cursor](#granting-access-for-moving-the-mouse-cursor)
- [How it works](#how-it-works)

<!-- /code_chunk_output -->

## How it's different from other apps

The main difference between this app and other apps is that this app **keeps your machine awake**, whereas other apps keep your machine from going to sleep. Let me explain.

If you use a messaging app like Slack, they are programmed to automatically change your status to `Away` when you don't interact with your machine for a certain amount of time. Apps which keep your machine from going to sleep will not be able to stop that. This app will prevent that, since **it is actively keeping your machine awake by moving the mouse pointer**.

**This comes very close to simulating an actual user using the machine, which the other apps cannot do.**

So if you want something that will keep your mac awake as long as you don't manually put it to sleep, then this is for you!

## How I use it

I always have this app working in the background for me whenever I work from home, so that I can take a break from work, strech my legs, go for a short walk, come back and still have my slack open and set to active! (not having to type in my password every time is awesome, for machines without fingerprint sensors anyway).

Also if I need to go out for longer, I just close the lid, and off goes my mac to sleep!

### Work from home benefit

Since this app actively keeps your machine awake, your messaging app will never go on idle and put your status as `Away`.

## Demo

You just click on `Start`, and AMM will take care of moving your mouse whenever it feels that the system has been left idle for a minute. It's as simple as this.

![](https://github.com/prashantgupta24/automatic-mouse-mover/blob/master/resources/amm-demo.gif)

## How to install

### Install from binary

1. Download the latest `amm.app.zip` from the [releases](https://github.com/prashantgupta24/automatic-mouse-mover/releases) page, unzip it, and copy the .app to your `Applications` folder like any other application.

1. Since the application is not notarized, you will need to right click on the .app and choose Open.

1. You will see a scary message that warns you about all the bad things that the app can do to your computer. If you are paranoid (fair enough, you don't really know me that well) then you can skip to the section which builds the app from source. That way you can see what exactly the app does (You can check that the application makes no connections to the internet whatsoever).

1. In case you do trust me, once you click on `Open`, you might encounter an initial `Access request` which I've discussed in the next section.

### Install from source

Make sure you have `go` installed. Once that is done, clone this repo and run `Make`, it should create the `amm.app` and open the folder where it was built for you. Copy the .app to your `Applications` folder like any other application.

Double click on the app, and the cute `mouse` should appear on your taskbar on top of your screen. Once you click on `Start`, you might encounter an initial `Access request` which I've discussed in the next section. If not, then you are all set!

## Granting access for moving the mouse cursor

While starting the app, you might see a message like the one below or an error stating `Mouse pointer cannot be moved`.

![](https://github.com/prashantgupta24/automatic-mouse-mover/blob/master/resources/request.jpg)

Don't worry, it's nothing sinister, it's just that Mac doesn't allow apps to gain accessibility to the computer by default (even standard apps like Automator, Firefox etc. who might want to access some features need to go through the same process).

In order to resolve this error you need to:

> Go to System Preferences -> Security & Privacy -> Privacy -> Accessibility and allow the `amm` app to gain access.

If you still see the error, try to quit and start the app again (the age-old way of fixing everything).

## How it works

Every 60 seconds, AMM uses [Activity tracker](https://github.com/prashantgupta24/activity-tracker) to track the various changes that happened in your system during that time, like cursor movement, mouse clicks, screen changes etc. Whenever `AMM` detects a change in the system, it knows that the system is busy and will not do anything. If not, it moves the mouse cursor ever so slightly, enough to keep your Mac awake for eternity.

> All code is public and open-sourced so no worrying if there's nefarious intention involved in recording your activity or not.

[version-badge]: https://img.shields.io/github/release/prashantgupta24/automatic-mouse-mover.svg
[releases]: https://github.com/prashantgupta24/automatic-mouse-mover/releases
[godoc-badge]: https://img.shields.io/badge/godoc-reference-blue.svg
[godoc-link]: https://godoc.org/github.com/prashantgupta24/automatic-mouse-mover/pkg/mousemover

## Contributions welcome!

Please feel free to contribute and make this open source app even better! Adding more pluggable activities to [Activity tracker](https://github.com/prashantgupta24/activity-tracker) will make sure that AMM works even better!
