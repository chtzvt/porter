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

func (p *PorterClient) List() (map[string]config.Door, error) {
	list := new(map[string]config.Door)

	if err := p.Call("GET", V1GetDoorList, DoorIDEmpty, list); err != nil {
		return nil, err
	}

	return *list, nil
}

func (p *PorterClient) GetState(id string) (config.Door, error) {
	state := new(config.Door)

	if err := p.Call("GET", V1GetDoorState, id, state); err != nil {
		return config.Door{}, err
	}

	return *state, nil
}

func (p *PorterClient) SendCmd(cmd, id string) (api.StatusMsg, error) {
	status := new(api.StatusMsg)

	if err := p.Call("PUT", cmd, id, status); err != nil {
		return api.StatusMsg{}, err
	}

	return *status, nil
}

func (p *PorterClient) Lock(id string) (api.StatusMsg, error) {
	return p.SendCmd(V1LockDoor, id)
}

func (p *PorterClient) Unlock(id string) (api.StatusMsg, error) {
	return p.SendCmd(V1UnlockDoor, id)
}

func (p *PorterClient) Open(id string) (api.StatusMsg, error) {
	return p.SendCmd(V1OpenDoor, id)
}

func (p *PorterClient) Close(id string) (api.StatusMsg, error) {
	return p.SendCmd(V1CloseDoor, id)
}

func (p *PorterClient) Trip(id string) (api.StatusMsg, error) {
	return p.SendCmd(V1TripDoor, id)
}
