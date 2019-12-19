package door

import (
	"encoding/json"
	"porter/config"
	"porter/hw"
	"sync"
	"time"
)

type Door struct {
	config.Door
	liftCtl *hw.LiftController `json:"-"`
	sensor  *hw.DoorSensor     `json:"-"`

	State  State `json:"door_sensor_current_state"`
	Locked bool  `json:"locked"`

	initialized bool `json:"-"`

	mutex sync.Mutex `json:"-"`

	monitorStarted bool      `json:"-"`
	monitorCtl     chan bool `json:"-"`

	LastCmdTimestamp         time.Time `json:"last_cmd_ts"`
	LastStateChangeTimestamp time.Time `json:"last_state_change_ts"`
}

type State int8

const (
	Open State = iota
	Closed
)

const MonitorSampleTime = 1 * time.Second

func New(c config.Door) *Door {
	door := Door{
		Door:    c,
		liftCtl: hw.NewLiftController(c.LiftCtlPin, c.LiftCtlInactiveState, c.LiftCtlTripMs),
		sensor:  hw.NewDoorSensor(c.SensorPin, c.SensorClosedState),
	}

	door.initialized = true
	door.startMonitor()

	return &door
}

func (d *Door) String() string {
	bytes, err := json.Marshal(d)
	if err != nil {
		return ""
	}

	return string(bytes)
}
