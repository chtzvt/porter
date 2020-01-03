package client

import "time"

const (
	V1GetDoorList string = "/api/v1/list"
	V1LockDoor    string = "/api/v1/lock/"
	V1UnlockDoor  string = "/api/v1/unlock/"
	V1OpenDoor    string = "/api/v1/open/"
	V1CloseDoor   string = "/api/v1/close/"
	V1TripDoor    string = "/api/v1/trip/"
)

type V1DoorState struct {
	Name                     string    `json:"name"`
	LiftCtlPin               int       `json:"lift_ctl_pin"`
	LiftCtlInactiveState     int       `json:"lift_ctl_inactive_state"`
	LiftCtlTripMs            int       `json:"lift_ctl_trip_time_ms"`
	SensorPin                int       `json:"door_sensor_pin"`
	SensorClosedState        int       `json:"door_sensor_closed_state"`
	State                    int       `json:"door_sensor_current_state"`
	Locked                   bool      `json:"locked"`
	LastCmdTimestamp         time.Time `json:"last_cmd_ts"`
	LastStateChangeTimestamp time.Time `json:"last_state_change_ts"`
}
