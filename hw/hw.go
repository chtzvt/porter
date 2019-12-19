// Package hw provides an abstraction layer for hardware-specific libraries, along with
// several higher-level control and sensor devices built from those abstracted primitives.
package hw

import (
	"github.com/stianeikeland/go-rpio"
)

// Init performs any required operations to initialize the hardware's GPIO interface
func Init() error {
	return rpio.Open()
}

// Close performs any required operations to un-initialize the hardware's GPIO interface
func Close() error {
	return rpio.Close()
}
