package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jacobsa/go-serial/serial"
	"go.viam.com/dynamixel/network"
	"go.viam.com/dynamixel/servo"
	"go.viam.com/dynamixel/servo/s_model"
)

var (
	portName = flag.String("port", "/dev/serial/by-id/usb-FTDI_USB__-__Serial_Converter_FT4TFT52-if00-port0", "the serial port path")
	servoID  = flag.Int("id", 8, "the ID of the servo to move")
	debug    = flag.Bool("debug", false, "show serial traffic")
)

func main() {
	flag.Parse()

	options := serial.OpenOptions{
		PortName:              *portName,
		BaudRate:              1000000,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 100,
	}

	s, err := serial.Open(options)

	if err != nil {
		fmt.Printf("error opening serial port: %s\n", err.Error())
		return
	}

	n := network.New(s)
	if *debug {
		n.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}

	_servo, err := s_model.New(n, *servoID)

	if err != nil {
		fmt.Printf("error initializing servo %d: %s", *servoID, err.Error())
		os.Exit(1)
	}

	err = _servo.Ping()
	if err != nil {
		fmt.Printf("ping error: %s\n", err)
		os.Exit(1)
	}

	pos, err := _servo.PresentPosition()
	fmt.Println(pos)
	if err != nil {
		fmt.Printf("pos error: %s\n", err)
		os.Exit(1)
	}

	var newPos int
	if pos < 1000 {
		newPos = pos + 500
	} else {
		newPos = pos - 300
	}
	fmt.Println(newPos)

	// err = servo.GoalAndTrack(newPos, true, _servo)
	err = _servo.SetGoalPosition(newPos)
	if err != nil {
		fmt.Printf("move error: %s\n", err)
	}

	pos, err = _servo.PresentPosition()
	if err != nil {
		fmt.Printf("pos error 2: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(pos)
}
