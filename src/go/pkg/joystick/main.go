package main

import (
	"fmt"
	"time"
)

func main() {
	newController := makeNewUdevXboxController()
	controller, err := newController()
	if err != nil {
		fmt.Println("Could not create udev controller", err)
	}

	bindController := makeBindXboxController()
	controller = bindController(controller)

	// Simulate button presses and releases
	buttonActions := []func() error{
		//controller.ButtonPressA, controller.ButtonReleaseA,
		//controller.ButtonPressB, controller.ButtonReleaseB,
		//controller.ButtonPressX, controller.ButtonReleaseX,
		//controller.ButtonPressY, controller.ButtonReleaseY,
		//controller.ButtonPressLB, controller.ButtonReleaseLB,
		//controller.ButtonPressRB, controller.ButtonReleaseRB,
		controller.MoveLT,
		controller.MoveRT,
		//controller.ButtonPressBack, controller.ButtonReleaseBack,
		//controller.ButtonPressStart, controller.ButtonReleaseStart,
		//controller.ButtonPressL3, controller.ButtonReleaseL3,
		//controller.ButtonPressR3, controller.ButtonReleaseR3,
	}

	// Simulate D-pad presses and releases
	dpadActions := []func() error{
		//controller.DpadPressUp, controller.DpadReleaseUp,
		//controller.DpadPressDown, controller.DpadReleaseDown,
		//controller.DpadPressLeft, controller.DpadReleaseLeft,
		//controller.DpadPressRight, controller.DpadReleaseRight,
	}

	// Simulate stick movements
	stickMovements := []func(float32, float32) error{
		//controller.LeftStickMove, controller.RightStickMove,
	}

	stickPositions := []struct{ x, y float32 }{
		//{32767, 0},
		//{-32767, 0},
		//{0, 32767},
		//{0, -32767},
		//{0, 0},
	}

	for {
		for _, action := range buttonActions {
			if err := action(); err != nil {
				panic(err)
			}
			time.Sleep(500 * time.Millisecond)
		}

		for _, action := range dpadActions {
			if err := action(); err != nil {
				panic(err)
			}
			time.Sleep(500 * time.Millisecond)
		}

		for _, move := range stickMovements {
			for _, pos := range stickPositions {
				if err := move(pos.x, pos.y); err != nil {
					panic(err)
				}
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}
