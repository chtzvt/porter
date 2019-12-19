package hw

import (
	"github.com/stianeikeland/go-rpio"
)

// Pin represents an abstraction of hardware-specific libraries, and this is used by all
// hardware-level components in Porter. This allows easy porting of Porter to run on other
// platforms or use GPIO libraries other than go-rpio, without having to modify any other code.
// As support for other platforms is developed, this process could be automatically integrated
// through the use of build tags which select the appropriate Pin drivers for each board
type Pin uint8
type State uint8

// Read a Pin State
func (p Pin) Read() State {
	return State(rpio.ReadPin(rpio.Pin(p)))
}

// Write a State to a Pin
func (p Pin) Write(state State) {
	rpio.WritePin(rpio.Pin(p), rpio.State(state))
}

// Input configures a GPIO Pin to act as an input, allowing it to be read
func (p Pin) Input() {
	rpio.PinMode(rpio.Pin(p), rpio.Input)
}

// Output configures a GPIO Pin to act as an output, allowing it to be written
func (p Pin) Output() {
	rpio.PinMode(rpio.Pin(p), rpio.Output)
}

// Toggle inverts the state of an output Pin
func (p Pin) Toggle() {
	rpio.TogglePin(rpio.Pin(p))
}
