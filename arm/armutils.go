package arm

import (
	//~ "fmt"
	"time"
	"dynamixel/servo"
)

// GoalAndTrack takes a goal position and some servos, and moves all servos to that position
// It will optionally track the joint positions and block until the servos are done moving
func GoalAndTrack(goal int, block bool, servos ...*servo.Servo) error{
	for _, s := range(servos){
		err := s.SetGoalPosition(goal)
		if err != nil {
			return err
		}
	}
	if block{
		allAtPos := false
		for !allAtPos{
			time.Sleep(200 * time.Millisecond)
			allAtPos = true
			for _, s := range(servos){
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
	}
	return nil
}

// WaitForMovement takes some servos, and will block until the servos are done moving
func WaitForMovement(servos ...*servo.Servo) error{
	allAtPos := false

	for !allAtPos{
		time.Sleep(200 * time.Millisecond)
		allAtPos = true
		for _, s := range(servos){
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

//Turn on torque for a select few servos
func TorqueOnVar(servos ...*servo.Servo) error{
	for _, s := range(servos){
		err := s.SetTorqueEnable(true)
		if err != nil {
			return err
		}
	}
	return nil
}

//Turn on torque for a select few servos
func TorqueOffVar(servos ...*servo.Servo) error{
	for _, s := range(servos){
		err := s.SetTorqueEnable(false)
		if err != nil {
			return err
		}
	}
	return nil
}
