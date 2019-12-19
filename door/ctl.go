package door

import (
	"fmt"
	"porter/hw"
	"time"
)

func (d *Door) GetSensors() (*hw.DoorSensor, *hw.LiftController) {
	return d.sensor, d.liftCtl
}

func (d *Door) IsOpen() bool {
	return d.State == Open
}

func (d *Door) IsClosed() bool {
	return d.State == Closed
}

func (d *Door) IsLocked() bool {
	return d.Locked
}

func (d *Door) Open() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.IsLocked() {
		return fmt.Errorf("%s is locked", d.Name)
	}

	return d.sendLiftCmd(Closed, false)
}

func (d *Door) Close() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.IsLocked() {
		return fmt.Errorf("%s is locked", d.Name)
	}

	return d.sendLiftCmd(Open, false)
}

func (d *Door) Lock() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.Locked = true
}

func (d *Door) Unlock() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.Locked = false
}

func (d *Door) Trip() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.IsLocked() {
		return fmt.Errorf("trip %s failed due to lockout", d.Name)
	}

	_ = d.sendLiftCmd(Closed, true)
	return nil
}

func (d *Door) readState() (State, error) {
	if !d.initialized {
		return Open, fmt.Errorf("%s has not been initialized", d.Name)
	}

	if d.sensor.Closed() {
		return Closed, nil
	}

	return Open, nil
}

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

	go d.liftCtl.Call()

	return nil
}

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

func (d *Door) stopMonitor() {
	if !d.monitorStarted {
		return
	}

	d.monitorCtl <- false

	d.monitorStarted = false
}
