package wx250s

import (
	"fmt"
	"time"
	"github.com/echolabsinc/dynamixel/servo"
	"github.com/echolabsinc/kinematics/joint"
)

type Wx250s struct{
	Joints map[string][]*servo.Servo
	LimbLengths map[string]float64
}


// TODO: Probably want to add/replace things with maybe Joint objects (doesn't yet exist)
// TODO: Find a better way of representing things like sleep angles, home angles, and limb lengths
// Currently the above is hacked together to make a proof of concept for fk/ik
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
		LimbLengths
	}
}

// Required methods for kinematics

// GetAllAngles will return a list of the angles of each servo
func (a *Wx250s) GetAllAngles() []int{
	var angles []int
	for i, s := range(a.GetAllServos()){
		pos, err := s.PresentPosition()
		if err != nil {
			fmt.Println(err)
		}
		angles = append(angles, pos)
	}
	return angles
}

func (a *Wx250s) LimbLengths() map[string]float64{
	return a.LimbLengths
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
		servo.WaitForMovementVar(a.Joints["Gripper"][0])
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
	return servo.GoalAndTrack(pos, block, a.GetServos(joint)...)
}

//Go back to the sleep position, ready to turn off torque
func (a *Wx250s) SleepPosition() error{
	sleepWait := false
	a.JointTo("Waist", 2048, sleepWait)
	a.JointTo("Shoulder", 840, sleepWait)
	a.JointTo("Wrist_rot", 2048, sleepWait)
	a.JointTo("Wrist", 2509, sleepWait)
	a.JointTo("Forearm_rot", 2048, sleepWait)
	a.JointTo("Elbow", 3090, sleepWait)
	a.WaitForMovement()
	return nil
}

//Go to the home position
func (a *Wx250s) HomePosition() error{
	wait := false
	for jointName, _ := range(a.Joints){
		a.JointTo(jointName, 2048, wait)
	}
	a.WaitForMovement()
	return nil
}

// WaitForMovement takes some servos, and will block until the servos are done moving
func (a *Wx250s) WaitForMovement() error{
	allAtPos := false

	for !allAtPos{
		time.Sleep(200 * time.Millisecond)
		allAtPos = true
		for _, s := range(a.GetAllServos()){
			isMoving, err := s.Moving()
			if err != nil {
				return err
			}
			// TODO: Make this configurable
			if isMoving != 0{
				allAtPos = false
			}
		}
	}
	return nil
}
