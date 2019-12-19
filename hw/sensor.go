package hw

type DoorSensor struct {
	HwPin       Pin
	ClosedState State
}

func NewDoorSensor(pin, closedState int) *DoorSensor {
	sensor := new(DoorSensor)
	sensor.HwPin = Pin(pin)
	sensor.HwPin.Input()
	sensor.ClosedState = State(closedState)
	return sensor
}

func (s *DoorSensor) Open() bool {
	return s.HwPin.Read() == s.ClosedState
}

func (s *DoorSensor) Closed() bool {
	return s.HwPin.Read() != s.ClosedState
}
