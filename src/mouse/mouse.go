package mouse

import "github.com/go-vgo/robotgo"

type Position struct {
	MouseX int
	MouseY int
}

func GetPosition() *Position {
	x, y := robotgo.GetMousePos()
	return &Position{
		MouseX: x,
		MouseY: y,
	}
}
