package mousemover

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

// func isPointerIdle(comm chan bool) {
// 	for {
// 		x1, y1 := robotgo.GetMousePos()
// 		time.Sleep(time.Second * 3)
// 		x2, y2 := robotgo.GetMousePos()

// 		if x1 == x2 && y1 == y2 {
// 			fmt.Println("idle")
// 			//comm <- true
// 		} else {
// 			fmt.Println("moving")
// 			comm <- false
// 		}
// 	}
// }

func checkIfMouseMoved(x1, y1, x2, y2 int, comm chan bool) {
	if x1 == x2 && y1 == y2 {
		fmt.Println("idle")
		//return false
		//comm <- true
	} else {
		fmt.Println("moving")
		comm <- false
		//return true
	}
}

func isMouseClick(comm chan bool) {

}

func moveMouse(comm chan bool) {
	ticker := time.NewTicker(time.Second * 3)
	val := true
	movePixel := 10
	x1, y1 := robotgo.GetMousePos()
	for {
		select {
		case <-ticker.C:
			fmt.Println("ticked")
			x2, y2 := robotgo.GetMousePos()
			checkIfMouseMoved(x1, y1, x2, y2, comm)
			if val {
				fmt.Println("moving mouse because idle")
				//x1, y1 := robotgo.GetMousePos()
				robotgo.Move(x2+movePixel, y2+movePixel)
				movePixel *= -1
			} else {
				val = true
			}
			x1 = x2
			y1 = y2
		case val = <-comm:
			fmt.Println("val received: ", val)
		}
	}
}
func main() {
	// log.SetOutput(os.Stdout)
	// log.Println("starting")
	// log.Println("logging")
	// //robotgo.ScrollMouse(100, "up")
	// robotgo.Move(100, 100)
	// x, y := robotgo.GetMousePos()
	// fmt.Println("pos: ", x, y)

	comm := make(chan bool)
	moveMouse(comm)
	//isPointerIdle(comm)

	// for {
	// 	// wheelDown := robotgo.AddEvent("wheelDown")
	// 	// wheelRight := robotgo.AddEvent("wheelRight")

	// 	// fmt.Println("wheelDown : ", wheelDown)
	// 	// fmt.Println("wheelRight : ", wheelRight)
	// 	count := 0
	// 	go func(count *int) {
	// 		for {
	// 			mleft := robotgo.AddEvent("mleft")
	// 			if mleft == 0 {
	// 				*count++
	// 				fmt.Println("mleft : ", *count)
	// 				time.Sleep(time.Millisecond * 500)
	// 			}
	// 		}

	// 	}(&count)

	// 	mright := robotgo.AddEvent("mright")

	// 	fmt.Println("mright : ", mright)

	// 	// if mleft {
	// 	// 	fmt.Println("you press... ", "mouse left button")
	// 	// }
	// }

}
