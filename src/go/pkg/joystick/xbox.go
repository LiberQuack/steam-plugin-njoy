package main

import (
	"github.com/bendahl/uinput"
	"math"
	"time"
)

type xboxController struct {
	gamepad uinput.Gamepad

	LeftStickMove  func(x float32, y float32) error
	RightStickMove func(x float32, y float32) error

	MoveLT func() error
	MoveRT func() error

	DpadPressUp      func() error
	DpadReleaseUp    func() error
	DpadPressDown    func() error
	DpadReleaseDown  func() error
	DpadPressLeft    func() error
	DpadReleaseLeft  func() error
	DpadPressRight   func() error
	DpadReleaseRight func() error

	ButtonPressA       func() error
	ButtonReleaseA     func() error
	ButtonPressB       func() error
	ButtonReleaseB     func() error
	ButtonPressX       func() error
	ButtonReleaseX     func() error
	ButtonPressY       func() error
	ButtonReleaseY     func() error
	ButtonPressLB      func() error
	ButtonReleaseLB    func() error
	ButtonPressRB      func() error
	ButtonReleaseRB    func() error
	ButtonPressBack    func() error
	ButtonReleaseBack  func() error
	ButtonPressStart   func() error
	ButtonReleaseStart func() error
	ButtonPressXbox    func() error
	ButtonReleaseXbox  func() error
	ButtonPressL3      func() error
	ButtonReleaseL3    func() error
	ButtonPressR3      func() error
	ButtonReleaseR3    func() error
}

type newXboxController func() (xboxController, error)

func makeNewUdevXboxController() newXboxController {
	return func() (xboxController, error) {
		gamepad, err := uinput.CreateGamepad("/dev/uinput", []byte("Virtual Xbox Controller"), 0x045e, 0x02d1)
		if err != nil {
			return xboxController{}, nil
		}

		return xboxController{gamepad: gamepad}, nil
	}
}

type bindXboxController func(controller xboxController) xboxController

func makeBindXboxController() bindXboxController {
	return func(xboxJoystick xboxController) xboxController {

		xboxJoystick.LeftStickMove = func(x float32, y float32) error {
			x = float32(math.Min(float64(x), uinput.MaximumAxisValue))
			x = float32(math.Max(float64(x), -uinput.MaximumAxisValue))
			y = float32(math.Min(float64(y), uinput.MaximumAxisValue))
			y = float32(math.Max(float64(y), -uinput.MaximumAxisValue))
			return xboxJoystick.gamepad.LeftStickMove(x, y)
		}

		xboxJoystick.RightStickMove = func(x float32, y float32) error {
			x = float32(math.Min(float64(x), uinput.MaximumAxisValue))
			x = float32(math.Max(float64(x), -uinput.MaximumAxisValue))
			y = float32(math.Min(float64(y), uinput.MaximumAxisValue))
			y = float32(math.Max(float64(y), -uinput.MaximumAxisValue))
			return xboxJoystick.gamepad.RightStickMove(x, y)
		}

		xboxJoystick.DpadPressUp = func() error {
			return xboxJoystick.gamepad.HatPress(uinput.HatUp, uinput.MaximumAxisValue)
		}

		xboxJoystick.DpadReleaseUp = func() error {
			return xboxJoystick.gamepad.HatRelease(uinput.HatUp)
		}

		xboxJoystick.DpadPressDown = func() error {
			return xboxJoystick.gamepad.HatPress(uinput.HatDown, uinput.MaximumAxisValue)
		}

		xboxJoystick.DpadReleaseDown = func() error {
			return xboxJoystick.gamepad.HatRelease(uinput.HatDown)
		}

		xboxJoystick.DpadPressLeft = func() error {
			return xboxJoystick.gamepad.HatPress(uinput.HatLeft, uinput.MaximumAxisValue)
		}

		xboxJoystick.DpadReleaseLeft = func() error {
			return xboxJoystick.gamepad.HatRelease(uinput.HatLeft)
		}

		xboxJoystick.DpadPressRight = func() error {
			return xboxJoystick.gamepad.HatPress(uinput.HatRight, uinput.MaximumAxisValue)
		}

		xboxJoystick.DpadReleaseRight = func() error {
			return xboxJoystick.gamepad.HatRelease(uinput.HatRight)
		}

		xboxJoystick.ButtonPressA = func() error {
			return xboxJoystick.gamepad.ButtonDown(304, 1)
		}

		xboxJoystick.ButtonReleaseA = func() error {
			return xboxJoystick.gamepad.ButtonUp(304)
		}

		xboxJoystick.ButtonPressB = func() error {
			return xboxJoystick.gamepad.ButtonDown(305, 1)
		}

		xboxJoystick.ButtonReleaseB = func() error {
			return xboxJoystick.gamepad.ButtonUp(305)
		}

		xboxJoystick.ButtonPressX = func() error {
			return xboxJoystick.gamepad.ButtonDown(307, 1)
		}

		xboxJoystick.ButtonReleaseX = func() error {
			return xboxJoystick.gamepad.ButtonUp(307)
		}

		xboxJoystick.ButtonPressY = func() error {
			return xboxJoystick.gamepad.ButtonDown(308, 1)
		}

		xboxJoystick.ButtonReleaseY = func() error {
			return xboxJoystick.gamepad.ButtonUp(308)
		}

		xboxJoystick.ButtonPressLB = func() error {
			return xboxJoystick.gamepad.ButtonDown(310, 1)
		}

		xboxJoystick.ButtonReleaseLB = func() error {
			return xboxJoystick.gamepad.ButtonUp(310)
		}

		xboxJoystick.ButtonPressRB = func() error {
			return xboxJoystick.gamepad.ButtonDown(311, 1)
		}

		xboxJoystick.ButtonReleaseRB = func() error {
			return xboxJoystick.gamepad.ButtonUp(311)
		}

		xboxJoystick.MoveLT = func() error {
			xboxJoystick.gamepad.TriggerMove(0x00, 255)
			time.Sleep(50 * time.Millisecond)
			xboxJoystick.gamepad.TriggerMove(0x01, 255)
			time.Sleep(50 * time.Millisecond)
			xboxJoystick.gamepad.TriggerMove(0x02, 255)
			time.Sleep(50 * time.Millisecond)
			xboxJoystick.gamepad.TriggerMove(0x03, 255)
			time.Sleep(50 * time.Millisecond)
			xboxJoystick.gamepad.TriggerMove(0x04, 255)
			time.Sleep(50 * time.Millisecond)
			xboxJoystick.gamepad.TriggerMove(0x05, 255)
			time.Sleep(50 * time.Millisecond)
			xboxJoystick.gamepad.TriggerMove(0x10, 255)
			time.Sleep(50 * time.Millisecond)
			xboxJoystick.gamepad.TriggerMove(0x11, 255)
			time.Sleep(50 * time.Millisecond)
			return nil
		}

		xboxJoystick.MoveRT = func() error {
			return xboxJoystick.gamepad.ButtonDown(0x05, 255)
		}

		xboxJoystick.ButtonPressBack = func() error {
			return xboxJoystick.gamepad.ButtonDown(314, 1)
		}

		xboxJoystick.ButtonReleaseBack = func() error {
			return xboxJoystick.gamepad.ButtonUp(314)
		}

		xboxJoystick.ButtonPressStart = func() error {
			return xboxJoystick.gamepad.ButtonDown(315, 1)
		}

		xboxJoystick.ButtonReleaseStart = func() error {
			return xboxJoystick.gamepad.ButtonUp(315)
		}

		xboxJoystick.ButtonPressXbox = func() error {
			return xboxJoystick.gamepad.ButtonDown(316, 1)
		}

		xboxJoystick.ButtonReleaseXbox = func() error {
			return xboxJoystick.gamepad.ButtonUp(316)
		}

		xboxJoystick.ButtonPressL3 = func() error {
			return xboxJoystick.gamepad.ButtonDown(317, 1)
		}

		xboxJoystick.ButtonReleaseL3 = func() error {
			return xboxJoystick.gamepad.ButtonUp(317)
		}

		xboxJoystick.ButtonPressR3 = func() error {
			return xboxJoystick.gamepad.ButtonDown(318, 1)
		}

		xboxJoystick.ButtonReleaseR3 = func() error {
			return xboxJoystick.gamepad.ButtonUp(318)
		}

		return xboxJoystick
	}
}

type NewXboxController func() (controller xboxController, err error)

func makeNewXboxController(newXboxController newXboxController, bindXboxController bindXboxController) NewXboxController {
	return func() (controller xboxController, err error) {
		controller, err = newXboxController()
		if err != nil {
			return
		}

		return bindXboxController(controller), nil
	}

}
