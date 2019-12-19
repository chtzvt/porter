package hw

import (
	"github.com/stianeikeland/go-rpio"
)

type Pin uint8
type State uint8

func (p Pin) Read() State {
	return State(rpio.ReadPin(rpio.Pin(p)))
}

func (p Pin) Write(state State) {
	rpio.WritePin(rpio.Pin(p), rpio.State(state))
}

func (p Pin) Input() {
	rpio.PinMode(rpio.Pin(p), rpio.Input)
}

func (p Pin) Output() {
	rpio.PinMode(rpio.Pin(p), rpio.Output)
}

func (p Pin) Toggle() {
	rpio.TogglePin(rpio.Pin(p))
}
