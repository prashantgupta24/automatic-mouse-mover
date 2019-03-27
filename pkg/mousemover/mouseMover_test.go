package mousemover

import (
	"fmt"
	"math"
	"testing"

	"github.com/go-vgo/robotgo"
	"github.com/stretchr/testify/assert"
)

func TestMouseMove(t *testing.T) {
	fmt.Println("starting test")
	movePixel := 10
	currentX, _ := robotgo.GetMousePos()
	commCh := make(chan bool, 1)
	moveMouse(movePixel, commCh)
	movedX, _ := robotgo.GetMousePos()
	assert.Equal(t, float64(movePixel), math.Abs(float64(movedX-currentX)))
}
