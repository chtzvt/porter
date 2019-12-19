package hw

import (
	"sync"
	"time"
)

type LiftController struct {
	HwPin         Pin
	InactiveState State
	TripTime      time.Duration
	mutex         sync.Mutex
}

func NewLiftController(pin, inactiveState, tripTime int) *LiftController {
	controller := new(LiftController)
	controller.InactiveState = State(inactiveState)
	controller.HwPin = Pin(pin)
	controller.HwPin.Output()
	controller.HwPin.Write(controller.InactiveState)
	controller.TripTime = (time.Duration)(tripTime) * time.Millisecond
	return controller
}

func (l *LiftController) Call() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.HwPin.Toggle() // Transition to active state
	time.Sleep(l.TripTime)
	l.HwPin.Toggle() // Return to inactive state
	l.HwPin.Write(0)
}
