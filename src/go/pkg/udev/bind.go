package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
	"time"
)

// Define necessary constants and structures
// These values might need to be adjusted based on your system

const absSize = 64

const (
	uinputMaxNameSize = 80
	uiDevCreate       = 0x5501
	uiDevDestroy      = 0x5502
	uiDevSetup        = 0x405c5503
	// this is for 64 length buffer to store name
	// for another length generate using : (len << 16) | 0x8000552C
	uiGetSysname = 0x8041552c
	uiSetEvBit   = 0x40045564
	uiSetKeyBit  = 0x40045565

	uiSetRelBit  = 0x40045566
	uiSetAbsBit  = 0x40045567
	uiSetAbsInfo = 0x40045568
	busUsb       = 0x03
)

const (
	UI_DEV_CREATE = 21761
	UI_SET_EVBIT  = 1074025828
	EV_ABS        = 3
	ABS_Z         = 2 // Typically represents the Z-axis (triggers)
)

type inputID struct {
	Bustype uint16
	Vendor  uint16
	Product uint16
	Version uint16
}

type uinputUserDev struct {
	Name       [uinputMaxNameSize]byte
	ID         inputID
	EffectsMax uint32
	Absmax     [absSize]int32
	Absmin     [absSize]int32
	Absfuzz    [absSize]int32
	Absflat    [absSize]int32
}

func main() {
	// Open the uinput device
	uinputFile, err := os.OpenFile("/dev/uinput", syscall.O_WRONLY|syscall.O_NONBLOCK, 0660)
	if err != nil {
		panic(err)
	}
	defer uinputFile.Close()

	//CHAT GPT, IS IT CORRECT, WHEN I RUN IT ABS Z IS NOT SIMULATED AS I EXPECTED
	//MY SYSTEM DO NOT RECOGNIZE IT
	if err := ioctl(uinputFile, UI_SET_EVBIT, EV_ABS); err != nil {
		panic(err)
	}
	if err := ioctl(uinputFile, UI_SET_EVBIT, ABS_Z); err != nil {
		panic(err)
	}

	// Setup key event for the 'A' button
	if err := ioctl(uinputFile, UI_SET_EVBIT, EV_KEY); err != nil {
		panic(err)
	}
	if err := ioctl(uinputFile, uiSetKeyBit, BTN_A); err != nil {
		panic(err)
	}

	// Create the device
	var uidev uinputUserDev
	copy(uidev.Name[:], []byte("Virtual Xbox Controller"))
	uidev.ID = inputID{
		Bustype: 0x03,
		Vendor:  0x045e,
		Product: 0x02d1,
		Version: 1,
	}
	uidev.Absmax[ABS_Z] = 255
	uidev.Absmin[ABS_Z] = -255
	uidev.Absfuzz[ABS_Z] = 0
	uidev.Absflat[ABS_Z] = 0

	uinputFile, err = createUsbDevice(uinputFile, uidev)
	if err != nil {
		panic(err)
	}

	// Emit a key press event for the 'A' button
	emit(uinputFile, EV_KEY, BTN_A, 1)
	err = syncEvents(uinputFile)
	if err != nil {
		panic(err)
	}

	// Emit a trigger press event
	if err := emit(uinputFile, EV_ABS, ABS_Z, 255); err != nil {
		panic(err)
	}

	err = syncEvents(uinputFile)
	if err != nil {
		panic(err)
	}
	time.Sleep(100 * time.Second)

	// Emit a trigger release event (if necessary)
	emit(uinputFile, EV_ABS, ABS_Z, 0)
}

func ioctl(deviceFile *os.File, cmd, ptr uintptr) error {
	_, _, errorCode := syscall.Syscall(syscall.SYS_IOCTL, deviceFile.Fd(), cmd, ptr)
	if errorCode != 0 {
		return errorCode
	}
	return nil
}

func emit(file *os.File, etype, code, value int) error {
	ev := InputEvent{
		Type:  uint16(etype),
		Code:  uint16(code),
		Value: int32(value),
	}

	buffer, err := inputEventToBuffer(ev)
	if err != nil {
		return err
	}

	if _, err := file.Write(buffer); err != nil {
		return err
	}
	return nil
}

type InputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

func createUsbDevice(deviceFile *os.File, dev uinputUserDev) (fd *os.File, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, dev)
	if err != nil {
		_ = deviceFile.Close()
		return nil, fmt.Errorf("failed to write user device buffer: %v", err)
	}
	_, err = deviceFile.Write(buf.Bytes())
	if err != nil {
		_ = deviceFile.Close()
		return nil, fmt.Errorf("failed to write uidev struct to device file: %v", err)
	}

	err = ioctl(deviceFile, uiDevCreate, uintptr(0))
	if err != nil {
		_ = deviceFile.Close()
		return nil, fmt.Errorf("failed to create device: %v", err)
	}

	time.Sleep(time.Millisecond * 200)

	return deviceFile, err
}

func syncEvents(deviceFile *os.File) (err error) {
	buf, err := inputEventToBuffer(InputEvent{
		Time:  syscall.Timeval{Sec: 0, Usec: 0},
		Type:  evSyn,
		Code:  uint16(synReport),
		Value: 0})
	if err != nil {
		return fmt.Errorf("writing sync event failed: %v", err)
	}
	_, err = deviceFile.Write(buf)
	return err
}

func inputEventToBuffer(iev InputEvent) (buffer []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, 24))
	err = binary.Write(buf, binary.LittleEndian, iev)
	if err != nil {
		return nil, fmt.Errorf("failed to write input event to buffer: %v", err)
	}
	return buf.Bytes(), nil
}

const (
	evSyn     = 0x00
	evKey     = 0x01
	evRel     = 0x02
	evAbs     = 0x03
	relX      = 0x0
	relY      = 0x1
	relHWheel = 0x6
	relWheel  = 0x8
	relDial   = 0x7

	absX     = 0x00
	absY     = 0x01
	absZ     = 0x02
	absRX    = 0x03
	absRY    = 0x04
	absRZ    = 0x05
	absHat0X = 0x10
	absHat0Y = 0x11

	absMtSlot       = 0x2f
	absMtTouchMajor = 0x30
	absMtPositionX  = 0x35
	absMtPositionY  = 0x36
	absMtTrackingId = 0x39

	synReport        = 0
	evMouseBtnLeft   = 0x110
	evMouseBtnRight  = 0x111
	evMouseBtnMiddle = 0x112
	evBtnTouch       = 0x14a
)

// Additional Constants for Key Events
const (
	EV_KEY = 0x01  // Event type for keys/buttons
	BTN_A  = 0x130 // Button code for 'A' button (example)
)
