package wx250s

import (
	"fmt"
	"time"
	"math"
	"dynamixel/servo"
)

type Wx250s struct{
	Waist *servo.Servo
	Shoulder [2]*servo.Servo
	Elbow [2]*servo.Servo
	Forearm_rot *servo.Servo
	Wrist *servo.Servo
	Wrist_rot *servo.Servo
	Gripper *servo.Servo
	servos []*servo.Servo
}


// TODO: No need to duplicate the servos here, find a better way to apply things to all of them
func New(servos []*servo.Servo) *Wx250s{
	return &Wx250s{
		Waist:       servos[0],
		Shoulder:    [2]*servo.Servo{servos[1],servos[2]},
		Elbow:       [2]*servo.Servo{servos[3],servos[4]},
		Forearm_rot: servos[5],
		Wrist:       servos[6],
		Wrist_rot:   servos[7],
		Gripper:     servos[8],
		servos:      servos,
	}
}

// Note: there are a million and a half different ways to move servos
// GoalPosition, GoalCurrent, GoalPWM, etc
// To start with I'll just use GoalPosition
// TODO: Support other movement types
// TODO: Configurable waiting for movement to complete or not
// TODO: write more TODOS for nice-to-have features


// GoalAndTrack takes a goal position and some servos, and moves all servos to that position
// It will track the joint positions and block until the servos are within a margin of the goal
func GoalAndTrack(goal int, servos ...*servo.Servo) error{
	for _, s := range(servos){
		err := s.SetGoalPosition(goal)
		if err != nil {
			return err
		}
	}

	allAtPos := false

	for !allAtPos{
		time.Sleep(10 * time.Millisecond)
		allAtPos = true
		for _, s := range(servos){
			curPos, err := s.PresentPosition()
			if err != nil {
				return err
			}
			// TODO: Make this configurable
			if math.Abs(float64(goal - curPos)) > 40{
				allAtPos = false
			}
		}
	}
	return nil
}

// Grippers are special because they use PWM by default rather than position control
// GoalAndTrack calls are commented out because I didn't realize that at first
// TODO: Track position or something rather than just have a timer
func (a *Wx250s) Close() error{
	//~ err := GoalAndTrack(1490, a.Gripper)
	err := a.Gripper.SetGoalPWM(-350)
	time.Sleep(1500 * time.Millisecond)
	return err
}

// See Close()
func (a *Wx250s) Open() error{
	//~ err := GoalAndTrack(3000, a.Gripper)
	err := a.Gripper.SetGoalPWM(250)
	if err != nil {
		return err
	}

	atPos := false
	for !atPos{
		var pos int
		pos, err = a.Gripper.PresentPosition()
		if err != nil {
			return err
		}
		// TODO: Don't harcode
		if pos < 3000{
			time.Sleep(50 * time.Millisecond)
		}else{
			atPos = true
		}
	}
	return err
}

// Print positions of all servos
// TODO: Print joint names, not just servo numbers
func (a *Wx250s) PrintPositions() error{
	posString := ""
	for i, s := range(a.servos){
		pos, err := s.PresentPosition()
		if err != nil {
			return err
		}
		posString = fmt.Sprintf("%s || %d : %d", posString, i, pos)
	}
	fmt.Println(posString)
	return nil
}

//Turn on torque for all servos
func (a *Wx250s) TorqueOn() error{
	for _, s := range(a.servos){
		err := s.SetTorqueEnable(true)
		if err != nil {
			return err
		}
	}
	return nil
}

//Turn off torque for all servos
func (a *Wx250s) TorqueOff() error{
	for _, s := range(a.servos){
		err := s.SetTorqueEnable(false)
		if err != nil {
			return err
		}
	}
	return nil
}

//Turn on torque for a select few servos
func (a *Wx250s) TorqueOnVar(servos ...*servo.Servo) error{
	for _, s := range(servos){
		err := s.SetTorqueEnable(true)
		if err != nil {
			return err
		}
	}
	return nil
}

//Turn on torque for a select few servos
func (a *Wx250s) TorqueOffVar(servos ...*servo.Servo) error{
	for _, s := range(servos){
		err := s.SetTorqueEnable(false)
		if err != nil {
			return err
		}
	}
	return nil
}

