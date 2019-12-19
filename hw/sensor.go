package hw

// DoorSensor provides a simple mechanism to access and manage
// a garage door sensor with configurable parameters
type DoorSensor struct {
	HwPin       Pin
	ClosedState State
}

// NewDoorSensor returns a new DoorSensor instance using the provided configuration
func NewDoorSensor(pin, closedState int) *DoorSensor {
	sensor := new(DoorSensor)
	sensor.HwPin = Pin(pin)
	sensor.HwPin.Input()
	sensor.ClosedState = State(closedState)
	return sensor
}

// Open returns true if a DoorSensor detects that the garage door is open
func (s *DoorSensor) Open() bool {
	return s.HwPin.Read() == s.ClosedState
}

// Closed returns true if a DoorSensor detects that the garage door is closed
func (s *DoorSensor) Closed() bool {
	return s.HwPin.Read() != s.ClosedState
}
