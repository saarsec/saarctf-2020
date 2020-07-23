package commands

import (
	"stahl4/service/internal/events"
	"stahl4/service/internal/stahl4"
	"encoding/json"
	"strconv"
)

type GetMaterialsCommand struct {
	Amount  uint8                `json:"amount,string"`
	Message string               `json:"message"`
	Client  *stahl4.RemoteClient `json:"-"` //Ignore in Marshal
	Server  *stahl4.RemoteServer `json:"-"` //Ignore in Marshal
}

func (getMaterialsCommand *GetMaterialsCommand) ToEvent(state *stahl4.State) stahl4.Event {
	var answer stahl4.Event
	if state.Logistic == nil {
		// we are not a logistic server, hence we ignore the request
		answer = nil
	} else {
		if state.Logistic.Busy {
			msg := "Already busy, ignoring your message."
			// option 2 because we were asked for materials instead of products
			answer = events.NewBusyEvent(msg, 2, getMaterialsCommand.Client, getMaterialsCommand.Server)
		} else if state.Logistic.MaterialsStored < uint16(getMaterialsCommand.Amount) {
			//not enough materials available
			msg := "Not enough materials available"
			state.Mux.Lock()
			state.Logistic.NumberOfNoMaterialAnswers += 1
			state.Mux.Unlock()
			answer = events.NewNoMaterialsEvent(msg, getMaterialsCommand.Client, getMaterialsCommand.Server)
		} else {
			//everything is fine and we can fetch materials

			//set state to busy
			state.Mux.Lock()
			state.Logistic.Busy = true
			state.Mux.Unlock()

			//sleep to simulate movement
			stahl4.MovementSleep()
			msg := "Delivered new materials."
			//update state
			state.Mux.Lock()
			state.Logistic.MaterialsStored = state.Logistic.MaterialsStored - uint16(getMaterialsCommand.Amount)
			answer = events.NewDeliveredMaterialsEvent(
				msg,
				getMaterialsCommand.Client,
				getMaterialsCommand.Server,
				getMaterialsCommand.Amount)
			//not busy anymore
			state.Logistic.Busy = false
			state.Mux.Unlock()
		}
	}
	return answer
}

/*
func (getMaterialsCommand *GetMaterialsCommand) UpdateState(state *stahl4.State, message string, option int) {
	// NOT TESTED
	if option == 0 {
		//fetching materials
		state.Logistic.MaterialsStored = state.Logistic.MaterialsStored + getMaterialsCommand.Amount

		log.Print("fetched " + strconv.Itoa(int(getMaterialsCommand.Amount)) + "materials. for " +
			getMaterialsCommand.Client.Ip.String())
	} else if option == 2 {
		//not enough materials
		panic("not implemented")
	} else {
		//busy
		log.Println("Server is busy, ignoring request of " + getMaterialsCommand.Client.Ip.String())
	}
	state.Logistic.LastMessage = message
	panic("not used anymore")
}
*/
func (getMaterialsCommand *GetMaterialsCommand) ToByteArray() []byte {
	return []byte(getMaterialsCommand.Message)
}

// method is used on non-pointer objects -> no pointer-receiver
func (getMaterialsCommand GetMaterialsCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":    "GetMaterialsCommand",
		"message": getMaterialsCommand.Message,
		"amount":  strconv.Itoa(int(getMaterialsCommand.Amount)),
	})
}

func (getMaterialsCommand *GetMaterialsCommand) GetClient() *stahl4.RemoteClient {
	return getMaterialsCommand.Client
}

func (getMaterialsCommand *GetMaterialsCommand) SetClient(client *stahl4.RemoteClient) {
	getMaterialsCommand.Client = client
}

func (getMaterialsCommand *GetMaterialsCommand) GetServer() *stahl4.RemoteServer {
	return getMaterialsCommand.Server
}

func (getMaterialsCommand *GetMaterialsCommand) SetServer(server *stahl4.RemoteServer) {
	getMaterialsCommand.Server = server
}

//
//
//

type GetProductsCommand struct {
	Amount  uint8                `json:"amount,string"`
	Message string               `json:"message"`
	Client  *stahl4.RemoteClient `json:"-"` //Ignore in Marshal
	Server  *stahl4.RemoteServer `json:"-"` //Ignore in Marshal
}

func (getProductsCommand *GetProductsCommand) ToEvent(state *stahl4.State) stahl4.Event {
	// brings products from production to logistic and checks if logistic got enough storage for them
	var answer stahl4.Event
	if state.Logistic == nil {
		// we are not a logistic server, hence we ignore the request
		answer = nil
	} else {
		freeSpace := state.Logistic.TotalProductCapacity - state.Logistic.ProductsStored
		if state.Logistic.Busy {
			msg := "Already busy, ignoring your message."
			// option 1 because we were asked for materials instead of products
			answer = events.NewBusyEvent(msg, 1, getProductsCommand.Client, getProductsCommand.Server)
		} else if freeSpace < uint16(getProductsCommand.Amount) {
			msg := "Not enough space for products."
			answer = events.NewNoProductSpaceEvent(msg, getProductsCommand.Client, getProductsCommand.Server)
		} else {
			//everything is correct

			//set state to busy
			state.Mux.Lock()
			state.Logistic.Busy = true
			state.Mux.Unlock()
			//sleep to simulate movement
			stahl4.MovementSleep()
			//update state
			state.Mux.Lock()
			state.Logistic.ProductsStored = state.Logistic.ProductsStored + uint16(getProductsCommand.Amount)
			msg := "fetched products."
			answer = events.NewFetchedProductsEvent(
				msg,
				getProductsCommand.Client,
				getProductsCommand.Server,
				getProductsCommand.Amount)
			//not busy anymore
			state.Logistic.Busy = false
			state.Mux.Unlock()
		}
	}
	return answer
}

/*
func (getProductsCommand *GetProductsCommand) UpdateState(state *stahl4.State, message string, option int) {
	panic("not used anymore")
}
*/
func (getProductsCommand *GetProductsCommand) ToByteArray() []byte {
	return []byte(getProductsCommand.Message)
}

// method is used on non-pointer objects -> no pointer-receiver
func (getProductsCommand GetProductsCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":    "GetProductsCommand",
		"message": getProductsCommand.Message,
		"amount":  strconv.Itoa(int(getProductsCommand.Amount)),
	})
}

func (getProductsCommand *GetProductsCommand) GetClient() *stahl4.RemoteClient {
	return getProductsCommand.Client
}

func (getProductsCommand *GetProductsCommand) SetClient(client *stahl4.RemoteClient) {
	getProductsCommand.Client = client
}

func (getProductsCommand *GetProductsCommand) GetServer() *stahl4.RemoteServer {
	return getProductsCommand.Server
}

func (getProductsCommand *GetProductsCommand) SetServer(server *stahl4.RemoteServer) {
	getProductsCommand.Server = server
}
