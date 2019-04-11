## Presenting the minimalistic Automatic-Mouse-Mover(AMM) app!

![GitHub release](https://img.shields.io/github/release/prashantgupta24/automatic-mouse-mover.svg)

Ever felt the need to keep your machine awake without having to resort to the age-old methods of playing a video or installing an app that you don't trust? **Well, not anymore!**

> Introducing the simplest app on the market that has the sole purpose of moving your mouse pointer at regular intervals so that your machine never sleeps! And best of all, it works **ONLY** when you are not working, so be rest assured that the mouse won't start moving on its own without the machine actually being idle.

## Demo

## How it works

See [Activity tracker](https://github.com/prashantgupta24/activity-tracker) for a full description of how the idleness is tracked.

## Error while running the app

In case you get an error from the app saying `mouse pointer cannot be moved.`, you need to give the app permission to control your mouse. Don't worry, it's nothing sinister, but Mac doesn't allow apps to gain accessibility to the computer by default (even standard apps like Automator, Firefox etc. who might want to access some features need to go through the same process)

In order to resolve this error you need to:

> Go to Security & Privacy -> Privacy -> Accessibility and allow the `amm` app to gain access.
