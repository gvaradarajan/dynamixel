package s_model


import (
	"io"
	"fmt"
	"github.com/viamrobotics/dynamixel/protocol/v2"
	reg "github.com/viamrobotics/dynamixel/registers"
	"github.com/viamrobotics/dynamixel/servo"
	"github.com/viamrobotics/dynamixel/utils"
)

var Registers reg.Map

// Determines the model of a servo, and configures it with the proper registry
func New(network io.ReadWriter, ID int) (*servo.Servo, error) {
	//So far, all servos I know of have their version number in the two bytes at 0x00
	b, err := v2.New(network).ReadData(ID, 0, 2)
	if err != nil {
		return nil, fmt.Errorf("error getting version for servo %d: %v\n", ID, err)
	}
	v, err2 := utils.BytesToInt(b)
	if err2 != nil {
		return nil, fmt.Errorf("error converting version bytes for servo %d: %v\n", ID, err2)
	}
	// Set Registry based on model
	//Note that the AX model uses protocol V1, currently not supported
	switch v {
	case 350:
		Registers = init_xl320()
	case 1020:
		Registers = init_xm430()
	case 1060:
		Registers = init_xl430()

	default:
		return nil, fmt.Errorf("Servo id %d version not supported: %d", ID, v)
	}
	s := servo.New(v2.New(network), Registers, ID)
	//~ s.ID = s.ServoID
	return s, nil
}

//~ func NewWithReturnLevel(network io.ReadWriter, ID int, returnLevel int) (*servo.Servo, error) {
	//~ s :=  servo.New(v2.New(network), Registers, ID), nil
	//~ return s, nil
//~ }
