package hw

import (
	"sync"
	"time"
)

// LiftController provides a simple mechanism to control and manage
// a garage door lift with configurable parameters
type LiftController struct {
	HwPin         Pin
	InactiveState State
	TripTime      time.Duration
	mutex         sync.Mutex
}

// NewLiftController returns a new LiftController instance using the provided configuration
func NewLiftController(pin, inactiveState, tripTime int) *LiftController {
	controller := new(LiftController)
	controller.InactiveState = State(inactiveState)
	controller.HwPin = Pin(pin)
	controller.HwPin.Output()
	controller.HwPin.Write(controller.InactiveState)
	controller.TripTime = (time.Duration)(tripTime) * time.Millisecond
	return controller
}

// Call signals a LiftController to operate the physical garage door lift
func (l *LiftController) Call() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.HwPin.Toggle() // Transition to active state
	time.Sleep(l.TripTime)
	l.HwPin.Toggle() // Return to inactive state

	time.Sleep(100 * time.Millisecond) // slight debounce before releasing mutex
}
