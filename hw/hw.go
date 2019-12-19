// Package hw provides a rudimentary shim driver for hardware-specific libraries, enabling easier portability
package hw

import (
	"github.com/stianeikeland/go-rpio"
)

func Init() error {
	return rpio.Open()
}

func Close() error {
	return rpio.Close()
}
