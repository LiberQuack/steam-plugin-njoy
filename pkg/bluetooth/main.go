package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
)

func main() {
	go serveLocal()

	d, err := dev.NewDevice("default")
	if err != nil {
		fmt.Printf("can't create new device : %v", err)
	}

	ble.SetDefaultDevice(d)

	joystickService1 := newJoystickService("12345678-1234-1234-1234-123456789ab1")
	joystickService2 := newJoystickService("12345678-1234-1234-1234-123456789ab2")
	ble.AddService(joystickService1)
	ble.AddService(joystickService2)

	deviceInfoService := newDeviceInfoService()

	ble.AddService(deviceInfoService)

	ctx := context.Background()

	go func() {
		ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		defer cancel()

		ctx = ble.WithSigHandler(ctx, cancel)
		fmt.Println("Advertising...")
		ble.AdvertiseNameAndServices(ctx, "XboxLikeJoystick", joystickService1.UUID, joystickService2.UUID)

	}()

	select {
	case <-ctx.Done():
		fmt.Println("stopped")
	}
}

func newDeviceInfoService() *ble.Service {
	deviceInfoService := ble.NewService(ble.MustParse("180A")) // Standard UUID for Device Information Service

	// Manufacturer Name String characteristic
	manufacturerCharacteristic := ble.NewCharacteristic(ble.MustParse("2A29"))
	manufacturerCharacteristic.HandleRead(ble.ReadHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		_ = binary.Write(rsp, binary.LittleEndian, []byte("Your Manufacturer Name\x00")) // Add your manufacturer name here
	}))
	deviceInfoService.AddCharacteristic(manufacturerCharacteristic)

	// Model Number String characteristic
	modelNumberCharacteristic := ble.NewCharacteristic(ble.MustParse("2A24"))
	modelNumberCharacteristic.HandleRead(ble.ReadHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		_ = binary.Write(rsp, binary.LittleEndian, []byte("Model 123\x00")) // Add your model number here
	}))
	deviceInfoService.AddCharacteristic(modelNumberCharacteristic)
	return deviceInfoService
}

func newJoystickService(uuid string) *ble.Service {
	joystickService := ble.NewService(ble.MustParse(uuid))

	// Buttons characteristic (simulated as a single byte for simplicity)
	buttonsCharacteristic := ble.NewCharacteristic(ble.MustParse("12345678-1234-1234-1234-123456789abd"))
	buttonsCharacteristic.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		fmt.Println("Button clicked")
		// Simulate button state (e.g., 0x01 for button A pressed)
		_ = binary.Write(rsp, binary.LittleEndian, byte(0x01))
	}))
	joystickService.AddCharacteristic(buttonsCharacteristic)

	// Left joystick characteristic (simulated as two bytes: X and Y axes)
	leftJoystickCharacteristic := ble.NewCharacteristic(ble.MustParse("12345678-1234-1234-1234-123456789abe"))
	leftJoystickCharacteristic.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		fmt.Println("Left analog activated")
		// Simulate joystick position (e.g., X=50, Y=100)
		_ = binary.Write(rsp, binary.LittleEndian, []byte{50, 100})
	}))
	joystickService.AddCharacteristic(leftJoystickCharacteristic)

	// Right joystick characteristic (similar to the left joystick)
	rightJoystickCharacteristic := ble.NewCharacteristic(ble.MustParse("12345678-1234-1234-1234-123456789abf"))
	rightJoystickCharacteristic.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		fmt.Println("Right analog activated")
		// Simulate joystick position (e.g., X=70, Y=80)
		_ = binary.Write(rsp, binary.LittleEndian, []byte{70, 80})
	}))
	joystickService.AddCharacteristic(rightJoystickCharacteristic)
	return joystickService
}

func serveLocal() error {
	fs := http.FileServer(http.Dir("."))

	// Handle all requests by serving a file of the same name.
	http.Handle("/", fs)

	// Specify the address and port to listen on.
	log.Println("Listening on http://localhost:8080/...")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
