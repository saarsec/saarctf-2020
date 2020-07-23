package commands

import (
	"stahl4/service/internal/events"
	"stahl4/service/internal/stahl4"
	"encoding/json"
)

type GetMaterialsBroadcastCommand struct {
	//Amount  uint8                `json:"amount,string"`
	Message string               `json:"message"`
	ID      string               `json:"id"`
	IP      string               `json:"ip"`
	Client  *stahl4.RemoteClient `json:"-"` //Ignore in Marshal
	Server  *stahl4.RemoteServer `json:"-"` //Ignore in Marshal
}

func (getMaterialsBroadcastCommand *GetMaterialsBroadcastCommand) ToEvent(state *stahl4.State) stahl4.Event {
	var answer stahl4.Event
	if state.Logistic == nil {
		// we are not a logistic server, hence we ignore the request
		answer = nil
	} else if state.Logistic.Busy {
		//we do not answer because we can not accept the request
		answer = nil
	} else if state.Logistic.MaterialsStored < state.Logistic.TotalMaterialCapacity/2 {
		//we do not answer because we can do not have enough materials
		answer = nil
	} else {
		//everything is fine and we can accept the broadcast
		msg := "Accepting your request."
		answer = events.NewGetMaterialsBroadcastAcceptEvent(msg, getMaterialsBroadcastCommand.ID, getMaterialsBroadcastCommand.Client,
			getMaterialsBroadcastCommand.Server)
	}
	return answer
}

/*
func (getMaterialsBroadcastCommand *GetMaterialsBroadcastCommand) UpdateState(state *stahl4.State, message string, option int) {
	panic("not used anymore")
}
*/

func (getMaterialsBroadcastCommand *GetMaterialsBroadcastCommand) ToByteArray() []byte {
	return []byte(getMaterialsBroadcastCommand.Message)
}

// method is used on non-pointer objects -> no pointer-receiver
func (getMaterialsBroadcastCommand GetMaterialsBroadcastCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":    "GetMaterialsBroadcastCommand",
		"message": getMaterialsBroadcastCommand.Message,
		"id":      getMaterialsBroadcastCommand.ID,
		"ip":      getMaterialsBroadcastCommand.IP,
		//"amount":  strconv.Itoa(int(getMaterialsBroadcastCommand.Amount)),
	})
}

func (getMaterialsBroadcastCommand *GetMaterialsBroadcastCommand) GetClient() *stahl4.RemoteClient {
	return getMaterialsBroadcastCommand.Client
}

func (getMaterialsBroadcastCommand *GetMaterialsBroadcastCommand) SetClient(client *stahl4.RemoteClient) {
	getMaterialsBroadcastCommand.Client = client
}

func (getMaterialsBroadcastCommand *GetMaterialsBroadcastCommand) GetServer() *stahl4.RemoteServer {
	return getMaterialsBroadcastCommand.Server
}

func (getMaterialsBroadcastCommand *GetMaterialsBroadcastCommand) SetServer(server *stahl4.RemoteServer) {
	getMaterialsBroadcastCommand.Server = server
}

func (getMaterialsBroadcastCommand *GetMaterialsBroadcastCommand) GetID() string {
	return getMaterialsBroadcastCommand.ID
}

func (getMaterialsBroadcastCommand *GetMaterialsBroadcastCommand) SetID(id string) {
	getMaterialsBroadcastCommand.ID = id
}

func (getMaterialsBroadcastCommand *GetMaterialsBroadcastCommand) GetIP() string {
	return getMaterialsBroadcastCommand.IP
}

func (getMaterialsBroadcastCommand *GetMaterialsBroadcastCommand) SetIP(ip string) {
	getMaterialsBroadcastCommand.IP = ip
}

//
//
//

type GetProductsBroadcastCommand struct {
	//Amount  uint8                `json:"amount,string"`
	Message string               `json:"message"`
	ID      string               `json:"id"`
	IP      string               `json:"ip"`
	Client  *stahl4.RemoteClient `json:"-"` //Ignore in Marshal
	Server  *stahl4.RemoteServer `json:"-"` //Ignore in Marshal
}

func (getProductsBroadcastCommand *GetProductsBroadcastCommand) ToEvent(state *stahl4.State) stahl4.Event {
	//
	var answer stahl4.Event
	if state.Logistic == nil {
		// we are not a logistic server, hence we ignore the request
		answer = nil
	} else if state.Logistic.Busy {
		//we do not answer because we can not accept the request
		answer = nil
	} else if (state.Logistic.TotalProductCapacity - state.Logistic.ProductsStored) < state.Logistic.TotalProductCapacity/2 {
		//we do not answer because we do not have enough free space
		answer = nil
	} else {
		//everything is correct and we can accept the broadcast
		msg := "Accepting your request."
		answer = events.NewGetProductsBroadcastAcceptEvent(
			msg,
			getProductsBroadcastCommand.ID,
			getProductsBroadcastCommand.Client,
			getProductsBroadcastCommand.Server)
	}
	return answer
}

/*
func (getProductsBroadcastCommand *GetProductsBroadcastCommand) UpdateState(state *stahl4.State, message string, option int) {
	panic("not used anymore")
}
*/
func (getProductsBroadcastCommand *GetProductsBroadcastCommand) ToByteArray() []byte {
	return []byte(getProductsBroadcastCommand.Message)
}

// method is used on non-pointer objects -> no pointer-receiver
func (getProductsBroadcastCommand GetProductsBroadcastCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":    "GetProductsBroadcastCommand",
		"message": getProductsBroadcastCommand.Message,
		"id":      getProductsBroadcastCommand.ID,
		"ip":      getProductsBroadcastCommand.IP,
		//"amount":  strconv.Itoa(int(getProductsBroadcastCommand.Amount)),
	})
}

func (getProductsBroadcastCommand *GetProductsBroadcastCommand) GetClient() *stahl4.RemoteClient {
	return getProductsBroadcastCommand.Client
}

func (getProductsBroadcastCommand *GetProductsBroadcastCommand) SetClient(client *stahl4.RemoteClient) {
	getProductsBroadcastCommand.Client = client
}

func (getProductsBroadcastCommand *GetProductsBroadcastCommand) GetServer() *stahl4.RemoteServer {
	return getProductsBroadcastCommand.Server
}

func (getProductsBroadcastCommand *GetProductsBroadcastCommand) SetServer(server *stahl4.RemoteServer) {
	getProductsBroadcastCommand.Server = server
}

func (getProductsBroadcastCommand *GetProductsBroadcastCommand) GetID() string {
	return getProductsBroadcastCommand.ID
}

func (getProductsBroadcastCommand *GetProductsBroadcastCommand) SetID(id string) {
	getProductsBroadcastCommand.ID = id
}
func (getProductsBroadcastCommand *GetProductsBroadcastCommand) GetIP() string {
	return getProductsBroadcastCommand.IP
}

func (getProductsBroadcastCommand *GetProductsBroadcastCommand) SetIP(ip string) {
	getProductsBroadcastCommand.IP = ip
}
