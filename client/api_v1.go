package client

import (
	"porter/api"
	"porter/config"
)

const (
	V1GetDoorList  string = "/api/v1/list"
	V1GetDoorState string = "/api/v1/state/"
	V1LockDoor     string = "/api/v1/lock/"
	V1UnlockDoor   string = "/api/v1/unlock/"
	V1OpenDoor     string = "/api/v1/open/"
	V1CloseDoor    string = "/api/v1/close/"
	V1TripDoor     string = "/api/v1/trip/"
)

func (p *PorterClient) ListDoors() (map[string]config.Door, error) {
	list := new(map[string]config.Door)

	if err := p.Call("GET", V1GetDoorList, DoorIDEmpty, list); err != nil {
		return nil, err
	}

	return *list, nil
}

func (p *PorterClient) GetDoorState(id string) (config.Door, error) {
	state := new(config.Door)

	if err := p.Call("GET", V1GetDoorState, id, state); err != nil {
		return config.Door{}, err
	}

	return *state, nil
}

func (p *PorterClient) SendDoorCmd(cmd, id string) (api.StatusMsg, error) {
	status := new(api.StatusMsg)

	if err := p.Call("PUT", cmd, id, status); err != nil {
		return api.StatusMsg{}, err
	}

	return *status, nil
}

func (p *PorterClient) LockDoor(id string) (api.StatusMsg, error) {
	return p.SendDoorCmd(V1LockDoor, id)
}

func (p *PorterClient) UnlockDoor(id string) (api.StatusMsg, error) {
	return p.SendDoorCmd(V1UnlockDoor, id)
}

func (p *PorterClient) OpenDoor(id string) (api.StatusMsg, error) {
	return p.SendDoorCmd(V1OpenDoor, id)
}

func (p *PorterClient) CloseDoor(id string) (api.StatusMsg, error) {
	return p.SendDoorCmd(V1CloseDoor, id)
}

func (p *PorterClient) TripDoor(id string) (api.StatusMsg, error) {
	return p.SendDoorCmd(V1TripDoor, id)
}
