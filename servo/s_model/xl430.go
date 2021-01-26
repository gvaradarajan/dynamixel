package s_model

import (
	reg "dynamixel/registers"
)
// New returns a new XL-430 registry
// See: https://emanual.robotis.com/docs/en/dxl/x/xl430-w250/
func init_xl430() reg.Map{
	x := 0

	Registers := reg.Map{
		// Note that we use -1 for nonexistent default values
		// Currently no servos have any defaults below 0
		// EEPROM: Persisted
		reg.ModelNumber:             {0x00, 2, reg.RO, x, x},
		reg.ModelInfo:               {0x02, 4, reg.RO, x, x},
		reg.FirmwareVersion:         {0x06, 1, reg.RO, x, x},
		reg.ServoID:                 {0x07, 1, reg.RW, 0, 252}, // renamed from ID for clarity
		reg.BaudRate:                {0x08, 1, reg.RW, 0, 7},   // 0=9600, 1=57600, 2=115200, 3=1Mbps
		reg.ReturnDelayTime:         {0x09, 1, reg.RW, 0, 254}, // usec = value*2
		reg.DriveMode:               {0x0a, 1, reg.RW, 0, 5},    //
		reg.OperatingMode:           {0x0b, 1, reg.RW, 0, 16},    //
		reg.SecondaryID:             {0x0c, 1, reg.RW, 0, 252},    //
		reg.ProtocolType:            {0x0d, 1, reg.RW, 1, 2},    //
		reg.HomingOffset:            {0x14, 4, reg.RW, -1044479, 1044479},    //
		reg.MovingThreshold:         {0x18, 4, reg.RW, 0, 1023},    // 0.229 [rev/min]
		reg.HighestLimitTemperature: {0x1f, 1, reg.RW, 0, 100},  // docs says not to set higher than default of 72C
		reg.MaxLimitVoltage:         {0x20, 2, reg.RW, 60, 140}, // volt = value*0.1
		reg.MinLimitVoltage:         {0x22, 2, reg.RW, 60, 140}, // volt = value*0.1
		reg.PWMLimit:                {0x24, 2, reg.RW, 0, 885}, // 0.113 [%]
		reg.VelocityLimit:           {0x2c, 4, reg.RW, 0, 1023}, //  0.229 [rev/min]
		reg.MaxPositionLimit:        {0x30, 4, reg.RW, 0, 4095},
		reg.MinPositionLimit:        {0x34, 4, reg.RW, 0, 4095},
		reg.AlarmShutdown:           {0x3f, 1, reg.RW, 0, 255},  // enum; see docs

		// RAM: Reset to default when power-cycled
		reg.TorqueEnable:          {0x40, 1, reg.RW, 0, 1},
		reg.Led:                   {0x41, 1, reg.RW, 0, 1},
		reg.StatusReturnLevel:     {0x44, 1, reg.RW, 0, 2},
		reg.RegisteredInstruction: {0x45, 1, reg.RO, 0, 1},
		reg.HardwareErrorStatus:   {0x46, 1, reg.RO, x, x},
		reg.VIGain:                {0x4c, 2, reg.RW, 0, 16383},
		reg.VDGain:                {0x4e, 2, reg.RW, 0, 16383},
		reg.DGain:                 {0x50, 2, reg.RW, 0, 16383},
		reg.IGain:                 {0x52, 2, reg.RW, 0, 16383},
		reg.PGain:                 {0x54, 2, reg.RW, 0, 16383},
		reg.Feedforward2ndGain:    {0x56, 2, reg.RW, 0, 16383},
		reg.Feedforward1stGain:    {0x58, 2, reg.RW, 0, 16383},
		reg.BusWatchdog:           {0x62, 1, reg.RW, 1, 127}, // 20 [msec]
		reg.GoalPWM:               {0x64, 2, reg.RW, -885, 885}, // 0.113 [%]
		reg.GoalVelocity:          {0x68, 4, reg.RW, -1023, 1023}, // 0.229 [rev/min]
		reg.ProfileAcceleration:   {0x6c, 4, reg.RW, 0, 32737}, // 214.577 [rev/min2]
		reg.ProfileVelocity:       {0x70, 4, reg.RW, 0, 32767}, // 0.229 [rev/min]
		reg.GoalPosition:          {0x74, 4, reg.RW, 0, 4095},
		reg.RealtimeTick:          {0x78, 2, reg.RO, 0, 32767},
		reg.Moving:                {0x7a, 1, reg.RO, 0, 1},
		reg.MovingStatus:          {0x7b, 1, reg.RO, x, x},
		reg.PresentPWM:            {0x7c, 2, reg.RO, x, x},
		reg.PresentLoad:           {0x7e, 2, reg.RO, -1000, 1000}, // 0.1%
		reg.PresentVelocity:       {0x80, 4, reg.RO, x, x},
		reg.PresentPosition:       {0x84, 4, reg.RO, x, x},
		reg.VelocityTrajectory:    {0x88, 4, reg.RO, x, x},
		reg.PositionTrajectory:    {0x8c, 4, reg.RO, x, x},
		reg.PresentVoltage:        {0x90, 2, reg.RO, x, x},
		reg.PresentTemperature:    {0x92, 1, reg.RO, x, x},
		// TODO: add Indirect Address and Indirect Data fields https://emanual.robotis.com/docs/en/dxl/x/xl430-w250/#indirect-address
	}
	return Registers
}

