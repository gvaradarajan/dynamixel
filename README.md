# Dynamixel

[![Build Status](https://travis-ci.org/echolabsinc/dynamixel.svg?branch=master)](https://travis-ci.org/echolabsinc/dynamixel)
[![Go Report Card](https://goreportcard.com/badge/github.com/echolabsinc/dynamixel)](https://goreportcard.com/report/github.com/echolabsinc/dynamixel)
[![GoDoc](https://godoc.org/github.com/echolabsinc/dynamixel?status.svg)](https://godoc.org/github.com/echolabsinc/dynamixel)

This packages provides a Go interface to Dynamixel servos.


## Example

```go
package main

import (
  "log"
  "fmt"
  "github.com/jacobsa/go-serial/serial"
  "dynamixel/network"
  "dynamixel/servo"
  "dynamixel/servo/xl430"
)

func GoalAndTrack(s *servo.Servo, pos int) error{
  curPos, err := s.PresentPosition()
  if err != nil {
    return err
  }
  if curPos == pos{
    fmt.Println("Already at ", pos, " nothing to do")
  }else{
    err = s.SetGoalPosition(pos)
    if err != nil {
      return err
    }
    for(curPos != pos){
      curPos, err = s.PresentPosition()
      if err != nil {
        return err
      }
      fmt.Println("Goal: ", pos, " currently at: ", curPos)
    }
  }
  return nil
}

func main() {
  options := serial.OpenOptions{
    PortName: "/dev/ttyUSB1",
    BaudRate: 1000000,
    DataBits: 8,
    StopBits: 1,
    MinimumReadSize: 0,
    InterCharacterTimeout: 100,
  }

  serial, err := serial.Open(options)
  if err != nil {
    log.Fatalf("error opening serial port: %v\n", err)
  }

  network := network.New(serial)
  servo, err := xl430.New(network, 2)
  if err != nil {
    log.Fatalf("error initializing servo: %v\n", err)
  }

  err = servo.Ping()
  if err != nil {
    log.Fatalf("error pinging servo: %v\n", err)
  }
  var ver int
  ver, err = servo.ModelNumber()
  if err != nil {
    log.Fatalf("error getting model num: %v\n", err)
  }
  fmt.Println(ver)

  err = servo.SetTorqueEnable(true)
  if err != nil {
    log.Fatalf("error setting Torque on\n", err)
  }
  err = GoalAndTrack(servo, 950)
  if err != nil {
    log.Fatalf("error setting goal position: %v\n", err)
  }

  err = GoalAndTrack(servo, 850)
  if err != nil {
    log.Fatalf("error setting goal position: %v\n", err)
  }
  err = servo.SetTorqueEnable(false)
}

```

More examples can be found in the [examples] [examples] directory of this repo.


## Documentation

The docs can be found at [godoc.org] [docs], as usual.
The API is based on the Dynamixel [v1 protocol] [proto] docs.


## License

[MIT] [license], obv.


## Author

[Adam Mckaig] [adammck] made this just for you.
[Peter LoVerso] [biotinker] updated it so that it works on things made this decade.




[ax]:       http://support.robotis.com/en/product/dynamixel/ax_series/dxl_ax_actuator.htm
[xl]:       http://support.robotis.com/en/product/dynamixel/xl-series/xl-320.htm
[docs]:     https://godoc.org/github.com/echolabsinc/dynamixel
[examples]: https://github.com/echolabsinc/dynamixel/tree/master/examples
[proto]:    http://support.robotis.com/en/product/dynamixel/ax_series/dxl_ax_actuator.htm#Control_Table
[license]:  https://github.com/echolabsinc/dynamixel/blob/master/LICENSE
