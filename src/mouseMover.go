package main

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

const (
	timeToCheck = 3
)

type mousePos struct {
	mouseX int
	mouseY int
}

func moveMouse(comm chan struct{}, quit chan struct{}) {
	ticker := time.NewTicker(time.Second * timeToCheck)
	isIdle := true
	movePixel := 10
	lastMousePos := getMousePos()
	for {
		select {
		case <-ticker.C:
			fmt.Println("ticked : ", isIdle)
			currentMousePos := getMousePos()
			if isIdle && isPointerIdle(currentMousePos, lastMousePos) {
				fmt.Println("moving mouse")
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
			fmt.Println("val received: ", isIdle)
		case <-quit:
			return
		}
	}
}

func isPointerIdle(currentMousePos, lastMousePos *mousePos) bool {
	fmt.Println("current : ", currentMousePos)
	fmt.Println("last : ", lastMousePos)

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

func isMouseClicked(comm chan struct{}) {
	go func() {
		for {
			mleft := robotgo.AddEvent("mleft")
			if mleft == 0 {
				//fmt.Println("mleft clicked")
				comm <- struct{}{}
			}
		}
	}()
}

func main() {
	comm := make(chan struct{})
	quit := make(chan struct{})
	go isMouseClicked(comm)
	moveMouse(comm, quit)
}
