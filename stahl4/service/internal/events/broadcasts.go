package events

import (
	"stahl4/service/internal/stahl4"
	"encoding/json"
)

type GetMaterialsBroadcastAcceptEvent struct {
	Message string               `json:"message"`
	ID      string               `json:"id"`
	Client  *stahl4.RemoteClient `json:"-"` //Ignore in Marshal
	Server  *stahl4.RemoteServer `json:"-"` //Ignore in Marshal
}

func (getMaterialsBroadcastAcceptEvent *GetMaterialsBroadcastAcceptEvent) UpdateState(state *stahl4.State) {
	//no state changes are needed but function is needed for interface
}

func (getMaterialsBroadcastAcceptEvent *GetMaterialsBroadcastAcceptEvent) ToByteArray() []byte {
	return []byte(getMaterialsBroadcastAcceptEvent.Message)
}

func (getMaterialsBroadcastAcceptEvent *GetMaterialsBroadcastAcceptEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":    "GetMaterialsBroadcastAcceptEvent",
		"message": getMaterialsBroadcastAcceptEvent.Message,
		"id":      getMaterialsBroadcastAcceptEvent.ID,
	})
}

func (getMaterialsBroadcastAcceptEvent *GetMaterialsBroadcastAcceptEvent) GetClient() *stahl4.RemoteClient {
	return getMaterialsBroadcastAcceptEvent.Client
}

func (getMaterialsBroadcastAcceptEvent *GetMaterialsBroadcastAcceptEvent) SetClient(client *stahl4.RemoteClient) {
	getMaterialsBroadcastAcceptEvent.Client = client
}

func (getMaterialsBroadcastAcceptEvent *GetMaterialsBroadcastAcceptEvent) GetServer() *stahl4.RemoteServer {
	return getMaterialsBroadcastAcceptEvent.Server
}

func (getMaterialsBroadcastAcceptEvent *GetMaterialsBroadcastAcceptEvent) SetServer(server *stahl4.RemoteServer) {
	getMaterialsBroadcastAcceptEvent.Server = server
}

func (getMaterialsBroadcastAcceptEvent *GetMaterialsBroadcastAcceptEvent) GetID() string {
	return getMaterialsBroadcastAcceptEvent.ID
}

func (getMaterialsBroadcastAcceptEvent *GetMaterialsBroadcastAcceptEvent) SetID(id string) {
	getMaterialsBroadcastAcceptEvent.ID = id
}

//
//
//

type GetProductsBroadcastAcceptEvent struct {
	Message string               `json:"message"`
	ID      string               `json:"id"`
	Client  *stahl4.RemoteClient `json:"-"` //Ignore in Marshal
	Server  *stahl4.RemoteServer `json:"-"` //Ignore in Marshal
}

func (getProductsBroadcastAcceptEvent *GetProductsBroadcastAcceptEvent) UpdateState(state *stahl4.State) {
	//no state changes are needed but function is needed for interface
}

func (getProductsBroadcastAcceptEvent *GetProductsBroadcastAcceptEvent) ToByteArray() []byte {
	return []byte(getProductsBroadcastAcceptEvent.Message)
}

func (getProductsBroadcastAcceptEvent *GetProductsBroadcastAcceptEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":    "GetProductsBroadcastAcceptEvent",
		"message": getProductsBroadcastAcceptEvent.Message,
		"id":      getProductsBroadcastAcceptEvent.ID,
	})
}

func (getProductsBroadcastAcceptEvent *GetProductsBroadcastAcceptEvent) GetClient() *stahl4.RemoteClient {
	return getProductsBroadcastAcceptEvent.Client
}

func (getProductsBroadcastAcceptEvent *GetProductsBroadcastAcceptEvent) SetClient(client *stahl4.RemoteClient) {
	getProductsBroadcastAcceptEvent.Client = client
}

func (getProductsBroadcastAcceptEvent *GetProductsBroadcastAcceptEvent) GetServer() *stahl4.RemoteServer {
	return getProductsBroadcastAcceptEvent.Server
}

func (getProductsBroadcastAcceptEvent *GetProductsBroadcastAcceptEvent) SetServer(server *stahl4.RemoteServer) {
	getProductsBroadcastAcceptEvent.Server = server
}

func (getProductsBroadcastAcceptEvent *GetProductsBroadcastAcceptEvent) GetID() string {
	return getProductsBroadcastAcceptEvent.ID
}

func (getProductsBroadcastAcceptEvent *GetProductsBroadcastAcceptEvent) SetID(id string) {
	getProductsBroadcastAcceptEvent.ID = id
}
