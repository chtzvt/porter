package door

import (
	"encoding/json"
	"porter/config"
	"porter/hw"
	"sync"
	"time"
)

// Type Door integrates the configuration parameters of a config.Door alongside a number of other
// properties that keep track of a garage door's running state
type Door struct {
	config.Door
	liftCtl        *hw.LiftController `json:"-"`
	sensor         hw.Pin             `json:"-"`
	initialized    bool               `json:"-"`
	mutex          sync.Mutex         `json:"-"`
	monitorStarted bool               `json:"-"`
	monitorCtl     chan bool          `json:"-"`
}

// MonitorSampleTime defines the interval at which a Door's sensor is polled for state changes
const MonitorSampleTime = 500 * time.Millisecond

const (
	TargetStateOpen int8 = iota
	TargetStateClosed
)

// New returns a new instance of a Door using the provided door configuration
func New(c config.Door) *Door {
	door := Door{
		Door:    c,
		liftCtl: hw.NewLiftController(c.LiftCtlPin, c.LiftCtlInactiveState, c.LiftCtlTripMs),
		sensor:  hw.Pin(c.SensorPin),
	}

	door.sensor.Input()

	door.initialized = true
	door.startMonitor()

	return &door
}

// GetSensors returns a Door's DoorSensor and LiftController, which are useful for debugging.
func (d *Door) GetSensors() (*hw.Pin, *hw.LiftController) {
	return &d.sensor, d.liftCtl
}

// String serializes a door's configuration as a JSON string
func (d *Door) String() string {
	bytes, err := json.Marshal(d)
	if err != nil {
		return ""
	}

	return string(bytes)
}
