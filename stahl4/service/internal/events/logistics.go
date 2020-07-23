package events

import (
	"stahl4/service/internal/stahl4"
	"encoding/json"
	"strconv"
)

type BusyEvent struct {
	Message string               `json:"message"`
	Option  uint8                `json:"option,string"`
	Client  *stahl4.RemoteClient `json:"-"`
	Server  *stahl4.RemoteServer `json:"-"`
}

func (busyEvent *BusyEvent) UpdateState(state *stahl4.State) {
	if state.Production != nil {
		//only do something if we are production server
		if busyEvent.Option == 1 {
			// Option 1 is set when we were asked for Products
			state.Mux.Lock()
			state.Production.AskedToGetProducts = false
			state.Mux.Unlock()
		} else {
			// Option 2 is set when we were asked for Materials
			state.Mux.Lock()
			state.Production.AskedForMaterials = false
			state.Mux.Unlock()
		}
	}
}

func (busyEvent *BusyEvent) ToByteArray() []byte {
	return []byte(busyEvent.Message)
}

func (busyEvent *BusyEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":    "BusyEvent",
		"message": busyEvent.Message,
		"option":  strconv.Itoa(int(busyEvent.Option)),
	})
}

func (busyEvent *BusyEvent) GetClient() *stahl4.RemoteClient {
	return busyEvent.Client
}

func (busyEvent *BusyEvent) SetClient(client *stahl4.RemoteClient) {
	busyEvent.Client = client
}

func (busyEvent *BusyEvent) GetServer() *stahl4.RemoteServer {
	return busyEvent.Server
}

func (busyEvent *BusyEvent) SetServer(server *stahl4.RemoteServer) {
	busyEvent.Server = server
}

//
//
//

type NoProductSpaceEvent struct {
	Message string               `json:"message"`
	Client  *stahl4.RemoteClient `json:"-"`
	Server  *stahl4.RemoteServer `json:"-"`
}

func (noProductSpaceEvent *NoProductSpaceEvent) UpdateState(state *stahl4.State) {
	if state.Production == nil {
		//we are not a production server hence we ignore the change because something went wrong
	} else {
		//no significant state changes are needed, hence only set request back to false
		state.Mux.Lock()
		state.Production.AskedToGetProducts = false
		state.Mux.Unlock()
	}
}

func (noProductSpaceEvent *NoProductSpaceEvent) ToByteArray() []byte {
	return []byte(noProductSpaceEvent.Message)
}

func (noProductSpaceEvent *NoProductSpaceEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":    "NoProductSpaceEvent",
		"message": noProductSpaceEvent.Message,
	})
}

func (noProductSpaceEvent *NoProductSpaceEvent) GetClient() *stahl4.RemoteClient {
	return noProductSpaceEvent.Client
}

func (noProductSpaceEvent *NoProductSpaceEvent) SetClient(client *stahl4.RemoteClient) {
	noProductSpaceEvent.Client = client
}

func (noProductSpaceEvent *NoProductSpaceEvent) GetServer() *stahl4.RemoteServer {
	return noProductSpaceEvent.Server
}

func (noProductSpaceEvent *NoProductSpaceEvent) SetServer(server *stahl4.RemoteServer) {
	noProductSpaceEvent.Server = server
}

//
//
//

type NoMaterialsEvent struct {
	Message string               `json:"message"`
	Client  *stahl4.RemoteClient `json:"-"`
	Server  *stahl4.RemoteServer `json:"-"`
}

func (noProductSpaceEvent *NoMaterialsEvent) UpdateState(state *stahl4.State) {
	if state.Production == nil {
		//we are not a production server hence we ignore the change because something went wrong
	} else {
		//no significant state changes are needed, hence only set request back to false
		state.Mux.Lock()
		state.Production.AskedForMaterials = false
		state.Mux.Unlock()
	}
}

func (noProductSpaceEvent *NoMaterialsEvent) ToByteArray() []byte {
	return []byte(noProductSpaceEvent.Message)
}

func (noProductSpaceEvent *NoMaterialsEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":    "NoMaterialsEvent",
		"message": noProductSpaceEvent.Message,
	})
}

func (noProductSpaceEvent *NoMaterialsEvent) GetClient() *stahl4.RemoteClient {
	return noProductSpaceEvent.Client
}

func (noProductSpaceEvent *NoMaterialsEvent) SetClient(client *stahl4.RemoteClient) {
	noProductSpaceEvent.Client = client
}

func (noProductSpaceEvent *NoMaterialsEvent) GetServer() *stahl4.RemoteServer {
	return noProductSpaceEvent.Server
}

func (noProductSpaceEvent *NoMaterialsEvent) SetServer(server *stahl4.RemoteServer) {
	noProductSpaceEvent.Server = server
}

//
//
//

type DeliveredMaterialsEvent struct {
	Message string               `json:"message"`
	Amount  uint8                `json:"amount,string"`
	Client  *stahl4.RemoteClient `json:"-"`
	Server  *stahl4.RemoteServer `json:"-"`
}

func (deliveredMaterialsEvent *DeliveredMaterialsEvent) UpdateState(state *stahl4.State) {
	//we got new materials from logistic
	if state.Production == nil {
		//we are not a production server hence we ignore the change because something went wrong
	} else {
		// need to cast uint8 to uint16
		newAmount := uint16(deliveredMaterialsEvent.Amount) + state.Production.MaterialAmount
		if newAmount > state.Production.MaterialCapacity {
			//no state changes, because this should never happen
		} else {
			state.Mux.Lock()
			state.Production.MaterialAmount = newAmount
			//set AskedForMaterials back to false (server sets it to true)
			state.Production.AskedForMaterials = false
			state.Mux.Unlock()
		}
		//log.Print("fetched " + strconv.Itoa(int(deliveredMaterialsEvent.Amount)) + "materials. for " +
		//	deliveredMaterialsEvent.Client.Ip.String())
	}
}

func (deliveredMaterialsEvent *DeliveredMaterialsEvent) ToByteArray() []byte {
	return []byte(deliveredMaterialsEvent.Message)
}

func (deliveredMaterialsEvent *DeliveredMaterialsEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":    "DeliveredMaterialsEvent",
		"message": deliveredMaterialsEvent.Message,
		"amount":  strconv.Itoa(int(deliveredMaterialsEvent.Amount)),
	})
}

func (deliveredMaterialsEvent *DeliveredMaterialsEvent) GetClient() *stahl4.RemoteClient {
	return deliveredMaterialsEvent.Client
}

func (deliveredMaterialsEvent *DeliveredMaterialsEvent) SetClient(client *stahl4.RemoteClient) {
	deliveredMaterialsEvent.Client = client
}

func (deliveredMaterialsEvent *DeliveredMaterialsEvent) GetServer() *stahl4.RemoteServer {
	return deliveredMaterialsEvent.Server
}

func (deliveredMaterialsEvent *DeliveredMaterialsEvent) SetServer(server *stahl4.RemoteServer) {
	deliveredMaterialsEvent.Server = server
}

//
//
//

type FetchedProductsEvent struct {
	Message string               `json:"message"`
	Amount  uint8                `json:"amount,string"`
	Client  *stahl4.RemoteClient `json:"-"` //ignore in Marshal
	Server  *stahl4.RemoteServer `json:"-"` //ignore in Marshal
}

func (fetchedProductsEvent *FetchedProductsEvent) UpdateState(state *stahl4.State) {
	//we sent products to logistic
	if state.Production == nil {
		//we are not a production server hence we ignore the change because something went wrong
	} else {
		// need to cast uint8 to uint16
		newAmount := state.Production.ProductAmount - uint16(fetchedProductsEvent.Amount)
		if state.Production.ProductAmount < uint16(fetchedProductsEvent.Amount) {
			//this should never happen, so we ignore the request
		} else {
			state.Mux.Lock()
			state.Production.ProductAmount = newAmount
			//set AskedToGetProducts back to false (server sets it to true)
			state.Production.AskedToGetProducts = false
			state.Mux.Unlock()
		}
		//log.Print("fetched " + strconv.Itoa(int(fetchedProductsEvent.Amount)) + "products. for " +
		//	fetchedProductsEvent.Client.Ip.String())
	}
}

func (fetchedProductsEvent *FetchedProductsEvent) ToByteArray() []byte {
	return []byte(fetchedProductsEvent.Message)
}

func (fetchedProductsEvent *FetchedProductsEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":    "FetchedProductsEvent",
		"message": fetchedProductsEvent.Message,
		"amount":  strconv.Itoa(int(fetchedProductsEvent.Amount)),
	})
}

func (fetchedProductsEvent *FetchedProductsEvent) GetClient() *stahl4.RemoteClient {
	return fetchedProductsEvent.Client
}

func (fetchedProductsEvent *FetchedProductsEvent) SetClient(client *stahl4.RemoteClient) {
	fetchedProductsEvent.Client = client
}

func (fetchedProductsEvent *FetchedProductsEvent) GetServer() *stahl4.RemoteServer {
	return fetchedProductsEvent.Server
}

func (fetchedProductsEvent *FetchedProductsEvent) SetServer(server *stahl4.RemoteServer) {
	fetchedProductsEvent.Server = server
}
