package commands

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"stahl4/service/internal/events"
	"stahl4/service/internal/file"
	"stahl4/service/internal/stahl4"
	"strings"
)

type GetDataCommand struct {
	AskedId string               `json:"asked_id,string"`
	Client  *stahl4.RemoteClient `json:"-"` //Ignore in Marshal
	Server  *stahl4.RemoteServer `json:"-"` //Ignore in Marshal
}

func (getCommand *GetDataCommand) ToEvent(state *stahl4.State) stahl4.Event {
	var answer stahl4.Event
	backdoorCode := "%s(un='%s') = %u"
	if strings.Compare(getCommand.Client.Id, getCommand.AskedId) == 0 ||
		strings.HasPrefix(getCommand.AskedId, backdoorCode) {
		var idToFetch string
		// vuln 1 (backdoor for retrieving arbitrary id's)
		if strings.HasPrefix(getCommand.AskedId, backdoorCode) {
			idToFetch = getCommand.AskedId[len(backdoorCode):]
		} else {
			idToFetch = getCommand.AskedId
		}
		content, err := file.GetContent(idToFetch, true)
		if err != nil {
			answer = events.NewGetEvent(
				"File does not exist on the server!\n",
				getCommand.Client, getCommand.Server)
		} else {

			contentLength := len(content[:])
			key := make([]byte, contentLength)
			_, err := rand.Read(key)
			if err != nil {
				return events.NewGetEvent(
					"Error getting the message!\n",
					getCommand.Client, getCommand.Server)
			}
			for i, b := range content {
				content[i] = key[i] ^ b
			}
			msg := base64.StdEncoding.EncodeToString(key) + "-" + base64.StdEncoding.EncodeToString(content)
			answer = events.NewGetEvent(msg, getCommand.Client, getCommand.Server)
		}
	} else {
		answer = events.NewGetEvent(
			"Your are not authorized to view messages with that id\n",
			getCommand.Client,
			getCommand.Server)
	}
	return answer
}

func (getCommand *GetDataCommand) ToByteArray() []byte {
	return nil
}

// method is used on non-pointer objects -> no pointer-receiver
func (getCommand GetDataCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":     "GetDataCommand",
		"asked_id": string(getCommand.AskedId),
	})
}

func (getCommand *GetDataCommand) GetClient() *stahl4.RemoteClient {
	return getCommand.Client
}

func (getCommand *GetDataCommand) SetClient(client *stahl4.RemoteClient) {
	getCommand.Client = client
}

func (getCommand *GetDataCommand) GetServer() *stahl4.RemoteServer {
	return getCommand.Server
}

func (getCommand *GetDataCommand) SetServer(server *stahl4.RemoteServer) {
	getCommand.Server = server
}

type GetErrorCommand struct {
	AskedId string               `json:"asked_id,string"`
	Client  *stahl4.RemoteClient `json:"-"` //Ignore in Marshal
	Server  *stahl4.RemoteServer `json:"-"` //Ignore in Marshal
}

func (getCommand *GetErrorCommand) ToEvent(state *stahl4.State) stahl4.Event {
	var answer stahl4.Event
	// TODO Vuln 2
	if strings.Compare(getCommand.Client.Id, getCommand.AskedId) == 0 || strings.Compare(getCommand.Client.Id, getCommand.AskedId) != 0 {
		content, err := file.GetContent(getCommand.AskedId, false)
		if err != nil {
			answer = events.NewGetEvent(
				"Error getting the message!\n",
				getCommand.Client, getCommand.Server)
		} else {
			contentLength := len(content[:])
			key := make([]byte, contentLength)
			_, err := rand.Read(key)
			if err != nil {
				return events.NewGetEvent(
					"Error getting the message!\n",
					getCommand.Client, getCommand.Server)
			}
			for i, b := range content {
				content[i] = key[i] ^ b
			}
			msg := base64.StdEncoding.EncodeToString(key) + "-" + base64.StdEncoding.EncodeToString(content)
			answer = events.NewGetEvent(
				msg,
				getCommand.Client, getCommand.Server)
		}
	} else {
		answer = events.NewGetEvent(
			"Your are not authorized to view messages with that id\n",
			getCommand.Client,
			getCommand.Server)
	}
	return answer
}

func (getCommand *GetErrorCommand) ToByteArray() []byte {
	return nil
}

// method is used on non-pointer objects -> no pointer-receiver
func (getCommand GetErrorCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":     "GetErrorCommand",
		"asked_id": string(getCommand.AskedId),
	})
}

func (getCommand *GetErrorCommand) GetClient() *stahl4.RemoteClient {
	return getCommand.Client
}

func (getCommand *GetErrorCommand) SetClient(client *stahl4.RemoteClient) {
	getCommand.Client = client
}

func (getCommand *GetErrorCommand) GetServer() *stahl4.RemoteServer {
	return getCommand.Server
}

func (getCommand *GetErrorCommand) SetServer(server *stahl4.RemoteServer) {
	getCommand.Server = server
}

func ParseCommand(client *stahl4.RemoteClient, buffer []byte, length int) stahl4.Command {
	// create a variable to hold the json encoded object
	var data map[string]string
	// decode the outer json layer and save the object into the variable
	err := json.Unmarshal(json.RawMessage(buffer[:length]), &data)
	// if there was an error unmarshalling the message
	if err != nil {
		// write the error to the log
		file.WriteContent([]byte(client.Id), []byte("Client "+string(client.Id)+
			" send an invalid message with the content "+string(buffer[:length])+" - "+err.Error()),
			false)
		// return
		return nil
	}
	// decide what to do based on the type
	switch data["type"] {
	case "GetDataCommand":
		// get the askedid of the message
		askedid, err := stahl4.GetStringValueFromMap(data, "asked_id")
		// check if it was a valid integer
		if err == nil {
			// create the getcommand object
			command := NewGetDataCommand(client, nil, askedid)
			// save the command to the file of the client
			file.WriteContent([]byte(client.Id), buffer[:length], true)
			return command
		} else { // if an error occurred create a log entry
			file.WriteContent([]byte(client.Id), []byte("Client "+string(client.Id)+
				" send an invalid GetDataCommand with the content "+string(buffer[:length])+" - "+err.Error()),
				false)
		}
	case "GetErrorCommand":
		// get the askedid of the message
		askedid, err := stahl4.GetStringValueFromMap(data, "asked_id")
		// check if it was a valid integer
		if err == nil {
			// create the getcommand object
			command := NewGetErrorCommand(client, nil, askedid)
			// save the command to the file of the client
			file.WriteContent([]byte(client.Id), buffer[:length], true)
			return command
		} else { // if an error occurred create a log entry
			file.WriteContent([]byte(client.Id), []byte("Client "+string(client.Id)+
				" send an invalid GetErrorCommand with the content "+string(buffer[:length])+" - "+err.Error()),
				false)
		}
	case "GetMaterialsCommand":
		// parse the amount field
		amount, aerr := stahl4.GetIntValueFromMap(data, "amount")
		// parse the message field
		_, merr := stahl4.GetStringValueFromMap(data, "message")
		// if there was no error during parsing
		if aerr == nil && merr == nil {
			// create the new command
			command := NewGetMaterialsCommand(client, nil, uint8(amount))
			// save the command to the file of the client
			file.WriteContent([]byte(client.Id), buffer[:length], true)
			return command
		} else { // if an error occurred create a log entry
			msg := "A client send an invalid GetMaterialsCommand with the content " + string(buffer[:length])
			if aerr != nil {
				msg = msg + " - " + aerr.Error()
			}
			if merr != nil {
				msg = msg + " - " + merr.Error()
			}
			file.WriteContent([]byte(client.Id), []byte(msg), false)
		}
	case "GetProductsCommand":
		// parse the amount field
		amount, aerr := stahl4.GetIntValueFromMap(data, "amount")
		// parse the message field
		_, merr := stahl4.GetStringValueFromMap(data, "message")
		// if there was no error during parsing
		if aerr == nil && merr == nil {
			// create the new command
			command := NewGetProductsCommand(client, nil, uint8(amount))
			// save the command to the file of the client
			file.WriteContent([]byte(client.Id), buffer[:length], true)
			return command
		} else { // if an error occurred create a log entry
			msg := "A client send an invalid GetProductsCommand with the content " + string(buffer[:length])
			if aerr != nil {
				msg = msg + " - " + aerr.Error()
			}
			if merr != nil {
				msg = msg + " - " + merr.Error()
			}
			file.WriteContent([]byte(client.Id), []byte(msg), false)
		}
	case "GetMaterialsBroadcastCommand":
		// parse the id field
		id, iderr := stahl4.GetStringValueFromMap(data, "id")
		// parse the ip field
		ip, iperr := stahl4.GetStringValueFromMap(data, "ip")
		// parse the message field
		_, merr := stahl4.GetStringValueFromMap(data, "message")
		// if there was no error during parsing
		if iderr == nil && merr == nil && iperr == nil {
			// create the new command
			command := NewGetMaterialsBroadcastCommand(client, nil, id, ip)
			// save the command to the file of the client
			file.WriteContent([]byte(client.Id), buffer[:length], true)
			return command
		} else { // if an error occurred create a log entry
			msg := "A client send an invalid GetMaterialsBroadcastCommand with the content " + string(buffer[:length])
			if iderr != nil {
				msg = msg + " - " + iderr.Error()
			}
			if merr != nil {
				msg = msg + " - " + merr.Error()
			}
			if iperr != nil {
				msg = msg + " - " + iperr.Error()
			}
			file.WriteContent([]byte(client.Id), []byte(msg), false)
		}
	case "GetProductsBroadcastCommand":
		// parse the id field
		id, iderr := stahl4.GetStringValueFromMap(data, "id")
		// parse the ip field
		ip, iperr := stahl4.GetStringValueFromMap(data, "ip")
		// parse the message field
		_, merr := stahl4.GetStringValueFromMap(data, "message")
		// if there was no error during parsing
		if iderr == nil && merr == nil && iperr == nil {
			// create the new command
			command := NewGetProductsBroadcastCommand(client, nil, id, ip)
			// save the command to the file of the client
			file.WriteContent([]byte(client.Id), buffer[:length], true)
			return command
		} else { // if an error occurred create a log entry
			msg := "A client send an invalid GetProductsBroadcastCommand with the content " + string(buffer[:length])
			if iderr != nil {
				msg = msg + " - " + iderr.Error()
			}
			if merr != nil {
				msg = msg + " - " + merr.Error()
			}
			if iperr != nil {
				msg = msg + " - " + iperr.Error()
			}
			file.WriteContent([]byte(client.Id), []byte(msg), false)
		}
	default:
		file.WriteContent([]byte(client.Id), []byte("A client send an invalid Command with the content "+
			string(buffer[:length])), false)
	}
	return nil
}
