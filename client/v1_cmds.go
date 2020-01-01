package client

import (
	"porter/api"
)

func (p *Client) List() (map[string]V1DoorState, error) {
	list := new(map[string]V1DoorState)

	if err := p.Call("GET", V1GetDoorList, DoorIDEmpty, list); err != nil {
		return nil, err
	}

	return *list, nil
}

func (p *Client) GetState(id string) (V1DoorState, error) {
	state := new(V1DoorState)

	if err := p.Call("GET", V1GetDoorState, id, state); err != nil {
		return V1DoorState{}, err
	}

	return *state, nil
}

func (p *Client) SendCmd(cmd, id string) (api.StatusMsg, error) {
	status := new(api.StatusMsg)

	if err := p.Call("PUT", cmd, id, status); err != nil {
		return api.StatusMsg{}, err
	}

	return *status, nil
}

func (p *Client) Lock(id string) (api.StatusMsg, error) {
	return p.SendCmd(V1LockDoor, id)
}

func (p *Client) Unlock(id string) (api.StatusMsg, error) {
	return p.SendCmd(V1UnlockDoor, id)
}

func (p *Client) Open(id string) (api.StatusMsg, error) {
	return p.SendCmd(V1OpenDoor, id)
}

func (p *Client) Close(id string) (api.StatusMsg, error) {
	return p.SendCmd(V1CloseDoor, id)
}

func (p *Client) Trip(id string) (api.StatusMsg, error) {
	return p.SendCmd(V1TripDoor, id)
}
