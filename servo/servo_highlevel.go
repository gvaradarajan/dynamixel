package servo

import (
	"time"
)

// High-level interface. (Most of this should be removed, or moved to a separate
// type which embeds or interacts with the servo type.)

func (s *Servo) posToAngle(p int) float64 {
	return (positionToAngle * float64(p)) - s.zeroAngle
}

func (s *Servo) angleToPos(angle float64) int {
	return int((s.zeroAngle + angle) * angleToPosition)
}

// Sets the origin angle (in degrees).
func (s *Servo) SetZero(offset float64) {
	s.zeroAngle = offset
}

// Returns the current position of the servo, relative to the zero angle.
func (s *Servo) Angle() (float64, error) {
	p, err := s.Position()

	if err != nil {
		return 0, err
	}

	return s.posToAngle(p), nil
}

// MoveTo sets the goal position of the servo by angle (in degrees), where zero
// is the midpoint, 150 deg is max left (clockwise), and -150 deg is max right
// (counter-clockwise). This is generally preferable to calling SetGoalPosition,
// which uses the internal uint16 representation.
func (s *Servo) MoveTo(angle float64) error {
	p := s.angleToPos(normalizeAngle(angle))
	return s.SetGoalPosition(p)
}

// Voltage returns the current voltage supplied. Unlike the underlying register,
// this is the actual voltage, not multiplied by ten.
func (s *Servo) Voltage() (float64, error) {
	val, err := s.PresentVoltage()
	if err != nil {
		return 0.0, err
	}

	// Convert the return value into actual volts.
	return (float64(val) / 10), nil
}

//
func normalizeAngle(d float64) float64 {
	if d > 180 {
		return normalizeAngle(d - 360)

	} else if d < -180 {
		return normalizeAngle(d + 360)

	} else {
		return d
	}
}

// GoalAndTrack takes a goal position and some servos, and moves all servos to that position
// It will optionally track the joint positions and block until the servos are done moving
func GoalAndTrack(goal int, block bool, servos ...*Servo) error{
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
func WaitForMovementVar(servos ...*Servo) error{
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
func TorqueOnVar(servos ...*Servo) error{
	for _, s := range(servos){
		err := s.SetTorqueEnable(true)
		if err != nil {
			return err
		}
	}
	return nil
}

//Turn on torque for a select few servos
func TorqueOffVar(servos ...*Servo) error{
	for _, s := range(servos){
		err := s.SetTorqueEnable(false)
		if err != nil {
			return err
		}
	}
	return nil
}
