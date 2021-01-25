package wx250s

import (
	"fmt"
	"time"
	"dynamixel/servo"
	"dynamixel/arm"
)

type Wx250s struct{
	Joints map[string][]*servo.Servo
}


// TODO: No need to duplicate the servos here, find a better way to apply things to all of them
func New(servos []*servo.Servo) *Wx250s{
	return &Wx250s{
		Joints: map[string][]*servo.Servo{
			"Waist":       []*servo.Servo{servos[0]},
			"Shoulder":    []*servo.Servo{servos[1],servos[2]},
			"Elbow":       []*servo.Servo{servos[3],servos[4]},
			"Forearm_rot": []*servo.Servo{servos[5]},
			"Wrist":       []*servo.Servo{servos[6]},
			"Wrist_rot":   []*servo.Servo{servos[7]},
			"Gripper":     []*servo.Servo{servos[8]},
		},
	}
}

// Note: there are a million and a half different ways to move servos
// GoalPosition, GoalCurrent, GoalPWM, etc
// To start with I'll just use GoalPosition
// TODO: Support other movement types
// TODO: Configurable waiting for movement to complete or not
// TODO: write more TODOS for nice-to-have features


// Grippers are special because they use PWM by default rather than position control
// Note that goal PWM values not in [-350:350] may cause the servo to overload, necessitating an arm reboot
// TODO: Track position or something rather than just have a timer
func (a *Wx250s) Close(block bool) error{
	err := a.Joints["Gripper"][0].SetGoalPWM(-350)
	if block{
		arm.WaitForMovement(a.Joints["Gripper"][0])
	}
	return err
}

// See Close()
func (a *Wx250s) Open() error{
	err := a.Joints["Gripper"][0].SetGoalPWM(250)
	if err != nil {
		return err
	}

	// We don't want to over-open
	atPos := false
	for !atPos{
		var pos int
		pos, err = a.Joints["Gripper"][0].PresentPosition()
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
	for i, s := range(a.GetAllServos()){
		pos, err := s.PresentPosition()
		if err != nil {
			return err
		}
		posString = fmt.Sprintf("%s || %d : %d", posString, i, pos)
	}
	fmt.Println(posString)
	return nil
}

// Return a slice containing all servos in the arm
func (a *Wx250s) GetAllServos() []*servo.Servo{
	var servos []*servo.Servo
	for _, v := range(a.Joints){
		for _, s := range v{
			servos = append(servos, s)
		}
	}
	return servos
}

// Return a slice containing all servos in the named joint
func (a *Wx250s) GetServos(joint string) []*servo.Servo{
	var servos []*servo.Servo
	for _, s := range(a.Joints[joint]){
		servos = append(servos, s)
	}
	return servos
}

// Set Acceleration for servos
func (a *Wx250s) SetAcceleration(accel int) error{
	for _, s := range(a.GetAllServos()){
		err := s.SetProfileAcceleration(accel)
		if err != nil {
			return err
		}
	}
	return nil
}

// Set Velocity for servos in travel time
// Recommended value 1000
func (a *Wx250s) SetVelocity(veloc int) error{
	for _, s := range(a.GetAllServos()){
		err := s.SetProfileVelocity(veloc)
		if err != nil {
			return err
		}
	}
	return nil
}

//Turn on torque for all servos
func (a *Wx250s) TorqueOn() error{
	for _, s := range(a.GetAllServos()){
		err := s.SetTorqueEnable(true)
		if err != nil {
			return err
		}
	}
	return nil
}

//Turn off torque for all servos
func (a *Wx250s) TorqueOff() error{
	for _, s := range(a.GetAllServos()){
		err := s.SetTorqueEnable(false)
		if err != nil {
			return err
		}
	}
	return nil
}

// Set a joint to a position
func (a *Wx250s) JointTo(joint string, pos int, block bool) error{
	fmt.Println("setting ", joint, " to ", pos)
	return arm.GoalAndTrack(pos, block, a.GetServos(joint)...)
}

//Go back to the sleep position, ready to turn off torque
func (a *Wx250s) SleepPosition() error{
	a.JointTo("Waist", 2034, true)
	a.JointTo("Shoulder", 840, true)
	a.JointTo("Wrist_rot", 2038, true)
	a.JointTo("Wrist", 2509, true)
	a.JointTo("Forearm_rot", 2064, true)
	a.JointTo("Elbow", 3090, true)
	return nil
}
