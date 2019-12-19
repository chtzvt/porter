// Package door provides a management plane that integrates hardware components with control and monitoring logic.
package door

import (
	"fmt"
	"time"
)

// IsOpen returns true if the door is open
func (d *Door) IsOpen() bool {
	return d.State == Open
}

// IsClosed returns true if the door is closed
func (d *Door) IsClosed() bool {
	return d.State == Closed
}

// IsLocked returns true if the door is locked
func (d *Door) IsLocked() bool {
	return d.Locked
}

// Open signals the door to open if it is unlocked, closed, and no other operations are pending
func (d *Door) Open() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.IsLocked() {
		return fmt.Errorf("%s is locked", d.Name)
	}

	return d.sendLiftCmd(Closed, false)
}

// Open signals the door to close if it is unlocked, open, and no other operations are pending
func (d *Door) Close() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.IsLocked() {
		return fmt.Errorf("%s is locked", d.Name)
	}

	return d.sendLiftCmd(Open, false)
}

// Lock enables a soft lock on the door, blocking any Open, Close, or Trip operations
func (d *Door) Lock() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.Locked = true
}

// Unlock removes a door's locked state
func (d *Door) Unlock() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.Locked = false
}

// Trip sends a Call to the door's LiftController, overriding state checks
// Trip requires a Door to be unlocked and have no other operations are pending
func (d *Door) Trip() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.IsLocked() {
		return fmt.Errorf("trip %s failed due to lockout", d.Name)
	}

	_ = d.sendLiftCmd(Closed, true)
	return nil
}

// readState returns the DoorSensor's State
func (d *Door) readState() (State, error) {
	if !d.initialized {
		return Open, fmt.Errorf("%s has not been initialized", d.Name)
	}

	if d.sensor.Closed() {
		return Closed, nil
	}

	return Open, nil
}

// sendLiftCmd enforces state checks on calls to the LiftController
func (d *Door) sendLiftCmd(requiredInitialState State, bypassStateCheck bool) error {
	if !d.initialized {
		return fmt.Errorf("%s has not been initialized", d.Name)
	}

	doorState, err := d.readState()
	if err != nil {
		return err
	}

	if bypassStateCheck || doorState == requiredInitialState {
		return fmt.Errorf("%s is already in the requested state", d.Name)
	}

	d.LastCmdTimestamp = time.Now()

	d.liftCtl.Call()

	return nil
}

// startMonitor starts a goroutine that periodically samples a door's DoorSensor
// This slightly improves performance because individual API calls don't trigger
// relatively expensive GPIO read operations to fetch the state of the door.
// It also allows state changes to be centrally tracked, so the server can produce
// a timestamp of the last detected change in state.
func (d *Door) startMonitor() {
	if d.monitorStarted {
		return
	}

	go (func() {
		for {
			select {
			case <-d.monitorCtl:
				return
			default:
				time.Sleep(MonitorSampleTime)
				if sample, err := d.readState(); err == nil {
					if d.State != sample {
						d.LastStateChangeTimestamp = time.Now()
						d.State = sample
					}
				}
			}
		}
	})()

	d.monitorStarted = true
}

// stopMonitor kills the door state monitor
func (d *Door) stopMonitor() {
	if !d.monitorStarted {
		return
	}

	d.monitorCtl <- false

	d.monitorStarted = false
}
