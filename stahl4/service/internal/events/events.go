package events

import (
	"encoding/json"
	"stahl4/service/internal/stahl4"
)

type GetEvent struct {
	Message string               `json:"message"`
	Client  *stahl4.RemoteClient `json:"-"`
	Server  *stahl4.RemoteServer `json:"-"`
}

func (getEvent *GetEvent) UpdateState(state *stahl4.State) {
	//no state changes are needed but function is needed for interface
}

func (getEvent *GetEvent) ToByteArray() []byte {
	return []byte(getEvent.Message)
}

func (getEvent *GetEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":    "GetEvent",
		"message": getEvent.Message,
	})
}

func (getEvent *GetEvent) GetClient() *stahl4.RemoteClient {
	return getEvent.Client
}

func (getEvent *GetEvent) SetClient(client *stahl4.RemoteClient) {
	getEvent.Client = client
}

func (getEvent *GetEvent) GetServer() *stahl4.RemoteServer {
	return getEvent.Server
}

func (getEvent *GetEvent) SetServer(server *stahl4.RemoteServer) {
	getEvent.Server = server
}

func ParseEvent(server *stahl4.RemoteServer, buffer []byte, length int) stahl4.Event {
	// create a variable to hold the json encoded object
	var data map[string]string
	// decode the outer json layer and save the object into the variable
	err := json.Unmarshal(json.RawMessage(buffer[:length]), &data)
	// if there was an error unmarshalling the message
	if err != nil {
		// write the error to the log
		/*file.WriteContent(int(server.Id), []byte("A Server send an invalid message with the content "+
		string(buffer[:length])+" - "+err.Error()), false)*/
		// return
		return nil
	}
	// decide what to do based on the type
	switch data["type"] {
	case "GetEvent":
		// parse the fields
		message, err := stahl4.GetStringValueFromMap(data, "message")
		// if there was no error during parsing return the event
		if err == nil {
			return NewGetEvent(message, nil, server)
		} /*else { // if an error occurred create a log entry
			file.WriteContent(int(server.Id), []byte("A server send an invalid GetEvent with the content "+
				string(buffer[:length])+" - "+err.Error()), false)
		}
		*/
	case "BusyEvent":
		// parse the fields
		message, err := stahl4.GetStringValueFromMap(data, "message")
		option, oerr := stahl4.GetIntValueFromMap(data, "option")
		// if there was no error during parsing return the event
		if err == nil && oerr == nil {
			return NewBusyEvent(message, uint8(option), nil, server)
		} /*else { // if an error occurred create a log entry
			msg := "A server send an invalid BusyEvent with the content " + string(buffer[:length])
			if err != nil {
				msg = msg + " - " + err.Error()
			}
			if oerr != nil {
				msg = msg + " - " + oerr.Error()
			}
			file.WriteContent(int(server.Id), []byte(msg), false)
		}
		*/
	case "FetchedProductsEvent":
		// parse the fields
		message, err := stahl4.GetStringValueFromMap(data, "message")
		amount, oerr := stahl4.GetIntValueFromMap(data, "amount")
		// if there was no error during parsing return the event
		if err == nil && oerr == nil {
			return NewFetchedProductsEvent(message, nil, server, uint8(amount))
		} /*else { // if an error occurred create a log entry
			msg := "A server send an invalid FetchedProductsEvent with the content " + string(buffer[:length])
			if err != nil {
				msg = msg + " - " + err.Error()
			}
			if oerr != nil {
				msg = msg + " - " + oerr.Error()
			}
			file.WriteContent(int(server.Id), []byte(msg), false)
		}
		*/
	case "DeliveredMaterialsEvent":
		// parse the fields
		message, err := stahl4.GetStringValueFromMap(data, "message")
		amount, oerr := stahl4.GetIntValueFromMap(data, "amount")
		// if there was no error during parsing return the event
		if err == nil && oerr == nil {
			return NewDeliveredMaterialsEvent(message, nil, server, uint8(amount))
		} /*else { // if an error occurred create a log entry
			msg := "A server send an invalid DeliveredMaterialsEvent with the content " + string(buffer[:length])
			if err != nil {
				msg = msg + " - " + err.Error()
			}
			if oerr != nil {
				msg = msg + " - " + oerr.Error()
			}
			file.WriteContent(int(server.Id), []byte(msg), false)
		}*/
	case "NoProductSpaceEvent":
		// parse the fields
		message, err := stahl4.GetStringValueFromMap(data, "message")
		// if there was no error during parsing return the event
		if err == nil {
			return NewNoProductSpaceEvent(message, nil, server)
		} /*else { // if an error occurred create a log entry
			file.WriteContent(int(server.Id), []byte("A server send an invalid NoProductSpaceEvent with the content "+
				string(buffer[:length])+" - "+err.Error()), false)
		}*/
	case "NoMaterialsEvent":
		// parse the fields
		message, err := stahl4.GetStringValueFromMap(data, "message")
		// if there was no error during parsing return the event
		if err == nil {
			return NewNoMaterialsEvent(message, nil, server)
		} /*else { // if an error occurred create a log entry
			file.WriteContent(int(server.Id), []byte("A server send an invalid NoMaterialsEvent with the content "+
				string(buffer[:length])+" - "+err.Error()), false)
		}*/
	case "GetMaterialsBroadcastAcceptEvent":
		// parse the fields
		message, err := stahl4.GetStringValueFromMap(data, "message")
		id, iderr := stahl4.GetStringValueFromMap(data, "id")
		// if there was no error during parsing return the event
		if err == nil && iderr == nil {
			return NewGetMaterialsBroadcastAcceptEvent(message, id, nil, server)
		} /*else { // if an error occurred create a log entry
			msg := "A server send an invalid GetMaterialsBroadcastAcceptEvent with the content " + string(buffer[:length])
			if err != nil {
				msg = msg + " - " + err.Error()
			}
			if iderr != nil {
				msg = msg + " - " + iderr.Error()
			}
			file.WriteContent(int(server.Id), []byte(msg), false)
		}*/
	case "GetProductsBroadcastAcceptEvent":
		// parse the fields
		message, err := stahl4.GetStringValueFromMap(data, "message")
		id, iderr := stahl4.GetStringValueFromMap(data, "id")
		// if there was no error during parsing return the event
		if err == nil && iderr == nil {
			return NewGetProductsBroadcastAcceptEvent(message, id, nil, server)
		} /*else { // if an error occurred create a log entry
			msg := "A server send an invalid GetProductsBroadcastAcceptEvent with the content " + string(buffer[:length])
			if err != nil {
				msg = msg + " - " + err.Error()
			}
			if iderr != nil {
				msg = msg + " - " + iderr.Error()
			}
			file.WriteContent(int(server.Id), []byte(msg), false)
		}*/
		// if no event matches log it
	default:
		/*file.WriteContent(int(server.Id), []byte("A server send an invalid Event with the content "+
		string(buffer[:length])), false)*/
	}
	return nil
}
