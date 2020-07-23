package commands

import (
	"stahl4/service/internal/events"
	"stahl4/service/internal/stahl4"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

// Globals
var client = &stahl4.RemoteClient{
	Mux:        sync.Mutex{},
	Id:         23,
	Ip:         nil,
	Connection: nil,
}

var server = &stahl4.RemoteServer{
	Mux:        sync.Mutex{},
	Ip:         nil,
	Connection: nil,
	Timestamp:  time.Time{},
	Typ:        7,
}

func TestCommands(t *testing.T) {
	// setup
	t.Run("CommandTestGroup", func(t *testing.T) {
		t.Run("TestGetCommand", testGetCommand)
		t.Run("TestGetMaterialsCommand", testGetMaterialsCommand)
		t.Run("testGetProductsCommand", testGetProductsCommand)
		t.Run("testGetMaterialsBroadcastCommand", testGetMaterialsBroadcastCommand)
		t.Run("testGetProductsBroadcastCommand", testGetProductsBroadcastCommand)
	})
	t.Run("CommandParserTests", func(t *testing.T) {
		t.Run("TestParseCommand", testParseCommand)
		t.Run("TestParseInvalidCommand", testParseInvalidCommand)
	})
	// teardown
	os.RemoveAll("./data/")
}

// Command tests
func testGetCommand(t *testing.T) {
	t.Parallel()
	var productionState = &stahl4.State{
		Production: &stahl4.Production{
			ProductCapacity:    300,
			MaterialCapacity:   300,
			ProductAmount:      100,
			MaterialAmount:     0,
			AskedForMaterials:  false,
			AskedToGetProducts: false,
		},
		Logistic: nil,
		Mux:      sync.Mutex{},
	}

	var logisticsState = &stahl4.State{
		Production: nil,
		Logistic: &stahl4.Logistic{
			Busy: false,
			TotalMaterialCapacity: 100,
			MaterialsStored:       100,
			TotalProductCapacity:  100,
			ProductsStored:        0,
		},
		Mux: sync.Mutex{},
	}
	var command = &GetDataCommand{
		AskedId: 23,
		Client:  nil,
		Server:  nil,
	}
	marshal := "{\"asked_id\":\"23\",\"type\":\"GetDataCommand\"}"
	t.Run("SetGetClient", func(t *testing.T) {
		helpSetGetClient(t, command)
	})
	t.Run("SetGetServer", func(t *testing.T) {
		helpSetGetServer(t, command)
	})

	t.Run("Marshal", func(t *testing.T) {
		//helpMarshal(t, command, marshal)
		m, err := json.Marshal(*command)
		if err != nil {
			t.Errorf("Error on marshalling")
		}
		if string(marshal) != string(m) {
			t.Errorf("Marshal error: expected: %s - got %s", string(marshal), string(m))
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		var object GetDataCommand
		err := json.Unmarshal([]byte(marshal), &object)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if object.AskedId != command.AskedId {
			t.Errorf("expected id: %d - got: %d", command.AskedId, object.AskedId)
		}
		if object.AskedId != command.AskedId {
			t.Errorf("expected id: %d - got: %d", command.AskedId, object.AskedId)
		}
	})
	t.Run("ToEvent", func(t *testing.T) {
		var message = "23,test123"
		err := ioutil.WriteFile("./data/messages.data", []byte(message), 0644)
		if err != nil {
			t.Skipf("Skipping ToEvent tests because there was a problem creating the file")
		}
		event := command.ToEvent(productionState)
		getEvent, ok := event.(*events.GetEvent)
		if !ok {
			t.Errorf("event was not of type GetEvent")
		} else {
			if getEvent.Message != message {
				t.Errorf("expected message: %s - got: %s", message, getEvent.Message)
			}
		}
		event = command.ToEvent(logisticsState)
		getEvent, ok = event.(*events.GetEvent)
		if !ok {
			t.Errorf("event was not of type GetEvent")
		} else {
			if getEvent.Message != message {
				t.Errorf("expected message: %s - got: %s", message, getEvent.Message)
			}
		}
		os.Remove("./data/messages.data")
		os.Remove("./data")
		// expect error message
		message = "File does not exist on the server!\n"
		event = command.ToEvent(productionState)
		getEvent, ok = event.(*events.GetEvent)
		if !ok {
			t.Errorf("event was not of type GetEvent")
		} else {
			if getEvent.Message != message {
				t.Errorf("expected message: %s - got: %s", message, getEvent.Message)
			}
		}
		event = command.ToEvent(logisticsState)
		getEvent, ok = event.(*events.GetEvent)
		if !ok {
			t.Errorf("event was not of type GetEvent")
		} else {
			if getEvent.Message != message {
				t.Errorf("expected message: %s - got: %s", message, getEvent.Message)
			}
		}
		// check for correct check of id
		command.AskedId = 0
		message = "Your are not authorized to view that file!\n"
		event = command.ToEvent(productionState)
		getEvent, ok = event.(*events.GetEvent)
		if !ok {
			t.Errorf("event was not of type GetEvent")
		} else {
			if getEvent.Message != message {
				t.Errorf("expected message: %s - got: %s", message, getEvent.Message)
			}
		}
		event = command.ToEvent(logisticsState)
		getEvent, ok = event.(*events.GetEvent)
		if !ok {
			t.Errorf("event was not of type GetEvent")
		} else {
			if getEvent.Message != message {
				t.Errorf("expected message: %s - got: %s", message, getEvent.Message)
			}
		}
	})
}

func testGetMaterialsCommand(t *testing.T) {
	t.Parallel()
	var productionState = &stahl4.State{
		Production: &stahl4.Production{
			ProductCapacity:    300,
			MaterialCapacity:   300,
			ProductAmount:      100,
			MaterialAmount:     0,
			AskedForMaterials:  false,
			AskedToGetProducts: false,
		},
		Logistic: nil,
		Mux:      sync.Mutex{},
	}

	var logisticsState = &stahl4.State{
		Production: nil,
		Logistic: &stahl4.Logistic{
			Busy: false,
			TotalMaterialCapacity: 100,
			MaterialsStored:       100,
			TotalProductCapacity:  100,
			ProductsStored:        0,
		},
		Mux: sync.Mutex{},
	}
	amount := uint8(100)
	message := "kekse"
	var command = &GetMaterialsCommand{
		Amount:  amount,
		Message: message,
		Client:  nil,
		Server:  nil,
	}
	marshal := "{\"amount\":\"" + strconv.Itoa(int(amount)) + "\",\"message\":\"" + message + "\",\"type\":\"GetMaterialsCommand\"}"
	t.Run("SetGetClient", func(t *testing.T) {
		helpSetGetClient(t, command)
	})
	t.Run("SetGetServer", func(t *testing.T) {
		helpSetGetServer(t, command)
	})
	t.Run("Marshal", func(t *testing.T) {
		m, err := json.Marshal(*command)
		if err != nil {
			t.Errorf("Error on marshalling")
		}
		if string(marshal) != string(m) {
			t.Errorf("Marshal error: expected: %s - got %s", string(marshal), string(m))
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		var object GetMaterialsCommand
		err := json.Unmarshal([]byte(marshal), &object)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if object.Message != command.Message {
			t.Errorf("expected message: %s - got: %s", command.Message, object.Message)
		}
		if object.Amount != command.Amount {
			t.Errorf("expected amount: %d - got: %d", command.Amount, object.Amount)
		}
	})
	t.Run("ToEvent", func(t *testing.T) {
		event := command.ToEvent(logisticsState)
		deliveredEvent, ok := event.(*events.DeliveredMaterialsEvent)
		if !ok {
			t.Errorf("event was not of type DeliveredMaterialsEvent but %s", reflect.TypeOf(event))
		} else {
			if deliveredEvent.Amount != amount {
				t.Errorf("expected amount: %d - got: %d", amount, deliveredEvent.Amount)
			}
		}
		logisticsState.Logistic.MaterialsStored = 5
		event = command.ToEvent(logisticsState)
		_, ok = event.(*events.NoMaterialsEvent)
		if !ok {
			t.Errorf("event was not of type NoMaterialsEvent but %s", reflect.TypeOf(event))
		}

		// expect ignore
		event = command.ToEvent(productionState)
		if event != nil {
			t.Errorf("expected no event when passing a production state.")
		}
		logisticsState.Logistic.Busy = true
		event = command.ToEvent(logisticsState)
		_, ok = event.(*events.BusyEvent)
		if !ok {
			t.Errorf("event was not of type BusyEvent but %s", reflect.TypeOf(event))
		}
		logisticsState.Logistic.Busy = false
	})
}

func testGetProductsCommand(t *testing.T) {
	t.Parallel()
	var productionState = &stahl4.State{
		Production: &stahl4.Production{
			ProductCapacity:    300,
			MaterialCapacity:   300,
			ProductAmount:      100,
			MaterialAmount:     0,
			AskedForMaterials:  false,
			AskedToGetProducts: false,
		},
		Logistic: nil,
		Mux:      sync.Mutex{},
	}

	var logisticsState = &stahl4.State{
		Production: nil,
		Logistic: &stahl4.Logistic{
			Busy: false,
			TotalMaterialCapacity: 100,
			MaterialsStored:       100,
			TotalProductCapacity:  100,
			ProductsStored:        0,
		},
		Mux: sync.Mutex{},
	}
	amount := uint8(100)
	message := "kekse"
	var command = &GetProductsCommand{
		Amount:  amount,
		Message: message,
		Client:  nil,
		Server:  nil,
	}
	marshal := "{\"amount\":\"" + strconv.Itoa(int(amount)) + "\",\"message\":\"" + message + "\",\"type\":\"GetProductsCommand\"}"
	t.Run("SetGetClient", func(t *testing.T) {
		helpSetGetClient(t, command)
	})
	t.Run("SetGetServer", func(t *testing.T) {
		helpSetGetServer(t, command)
	})
	t.Run("Marshal", func(t *testing.T) {
		m, err := json.Marshal(*command)
		if err != nil {
			t.Errorf("Error on marshalling")
		}
		if string(marshal) != string(m) {
			t.Errorf("Marshal error: expected: %s - got %s", string(marshal), string(m))
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		var object GetProductsCommand
		err := json.Unmarshal([]byte(marshal), &object)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if object.Message != command.Message {
			t.Errorf("expected message: %s - got: %s", command.Message, object.Message)
		}
		if object.Amount != command.Amount {
			t.Errorf("expected amount: %d - got: %d", command.Amount, object.Amount)
		}
	})
	t.Run("ToEvent", func(t *testing.T) {
		event := command.ToEvent(logisticsState)
		deliveredEvent, ok := event.(*events.FetchedProductsEvent)
		if !ok {
			t.Errorf("event was not of type FetchedProductsEvent but %s", reflect.TypeOf(event))
		} else {
			if deliveredEvent.Amount != amount {
				t.Errorf("expected amount: %d - got: %d", amount, deliveredEvent.Amount)
			}
		}
		logisticsState.Logistic.ProductsStored = 95
		event = command.ToEvent(logisticsState)
		_, ok = event.(*events.NoProductSpaceEvent)
		if !ok {
			t.Errorf("event was not of type NoProductSpaceEvent but %s", reflect.TypeOf(event))
		}

		// expect ignore
		event = command.ToEvent(productionState)
		if event != nil {
			t.Errorf("expected no event when passing a production state.")
		}
		logisticsState.Logistic.Busy = true
		event = command.ToEvent(logisticsState)
		_, ok = event.(*events.BusyEvent)
		if !ok {
			t.Errorf("event was not of type BusyEvent but %s", reflect.TypeOf(event))
		}
		logisticsState.Logistic.Busy = false
	})
}

func testGetMaterialsBroadcastCommand(t *testing.T) {
	t.Parallel()
	var productionState = &stahl4.State{
		Production: &stahl4.Production{
			ProductCapacity:    300,
			MaterialCapacity:   300,
			ProductAmount:      100,
			MaterialAmount:     0,
			AskedForMaterials:  false,
			AskedToGetProducts: false,
		},
		Logistic: nil,
		Mux:      sync.Mutex{},
	}

	var logisticsState = &stahl4.State{
		Production: nil,
		Logistic: &stahl4.Logistic{
			Busy: false,
			TotalMaterialCapacity: 100,
			MaterialsStored:       100,
			TotalProductCapacity:  100,
			ProductsStored:        0,
		},
		Mux: sync.Mutex{},
	}
	id := "1"
	ip := "123.456.789.012"
	message := "kekse"
	var command = &GetMaterialsBroadcastCommand{
		Message: message,
		ID:      id,
		IP:      ip,
		Client:  nil,
		Server:  nil,
	}
	marshal := "{\"id\":\"" + id + "\",\"ip\":\"" + ip + "\",\"message\":\"" + message + "\",\"type\":\"GetMaterialsBroadcastCommand\"}"
	t.Run("SetGetClient", func(t *testing.T) {
		helpSetGetClient(t, command)
	})
	t.Run("SetGetServer", func(t *testing.T) {
		helpSetGetServer(t, command)
	})
	t.Run("Marshal", func(t *testing.T) {
		m, err := json.Marshal(*command)
		if err != nil {
			t.Errorf("Error on marshalling")
		}
		if string(marshal) != string(m) {
			t.Errorf("Marshal error: expected: %s - got %s", string(marshal), string(m))
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		var object GetMaterialsBroadcastCommand
		err := json.Unmarshal([]byte(marshal), &object)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if object.Message != command.Message {
			t.Errorf("expected message: %s - got: %s", command.Message, object.Message)
		}
		if object.ID != command.ID {
			t.Errorf("expected id: %s - got: %s", command.ID, object.ID)
		}
		if object.IP != command.IP {
			t.Errorf("expected ip: %s - got: %s", command.IP, object.IP)
		}
	})
	t.Run("ToEvent", func(t *testing.T) {
		logisticsState.Logistic.MaterialsStored = 51
		event := command.ToEvent(logisticsState)
		acceptEvent, ok := event.(*events.GetMaterialsBroadcastAcceptEvent)
		if !ok {
			t.Errorf("event was not of type GetMaterialsBroadcastAcceptEvent")
		} else {
			if acceptEvent.ID != id {
				t.Errorf("expected id: %s - got: %s", id, acceptEvent.ID)
			}
		}

		// expect ignore
		logisticsState.Logistic.MaterialsStored = 49
		event = command.ToEvent(logisticsState)
		if event != nil {
			t.Errorf("expected no event when passing a production state.")
		}

		event = command.ToEvent(productionState)
		if event != nil {
			t.Errorf("expected no event when passing a production state.")
		}
	})
}

func testGetProductsBroadcastCommand(t *testing.T) {
	t.Parallel()
	var productionState = &stahl4.State{
		Production: &stahl4.Production{
			ProductCapacity:    300,
			MaterialCapacity:   300,
			ProductAmount:      100,
			MaterialAmount:     0,
			AskedForMaterials:  false,
			AskedToGetProducts: false,
		},
		Logistic: nil,
		Mux:      sync.Mutex{},
	}

	var logisticsState = &stahl4.State{
		Production: nil,
		Logistic: &stahl4.Logistic{
			Busy: false,
			TotalMaterialCapacity: 100,
			MaterialsStored:       100,
			TotalProductCapacity:  100,
			ProductsStored:        0,
		},
		Mux: sync.Mutex{},
	}
	id := "1"
	ip := "123.456.789.012"
	message := "kekse"
	var command = &GetProductsBroadcastCommand{
		Message: message,
		ID:      id,
		IP:      ip,
		Client:  nil,
		Server:  nil,
	}
	marshal := "{\"id\":\"" + id + "\",\"ip\":\"" + ip + "\",\"message\":\"" + message + "\",\"type\":\"GetProductsBroadcastCommand\"}"
	t.Run("SetGetClient", func(t *testing.T) {
		helpSetGetClient(t, command)
	})
	t.Run("SetGetServer", func(t *testing.T) {
		helpSetGetServer(t, command)
	})
	t.Run("Marshal", func(t *testing.T) {
		m, err := json.Marshal(*command)
		if err != nil {
			t.Errorf("Error on marshalling")
		}
		if string(marshal) != string(m) {
			t.Errorf("Marshal error: expected: %s - got %s", string(marshal), string(m))
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		var object GetMaterialsBroadcastCommand
		err := json.Unmarshal([]byte(marshal), &object)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if object.Message != command.Message {
			t.Errorf("expected message: %s - got: %s", command.Message, object.Message)
		}
		if object.ID != command.ID {
			t.Errorf("expected id: %s - got: %s", command.ID, object.ID)
		}
		if object.IP != command.IP {
			t.Errorf("expected ip: %s - got: %s", command.IP, object.IP)
		}
	})
	t.Run("ToEvent", func(t *testing.T) {
		logisticsState.Logistic.ProductsStored = 49
		event := command.ToEvent(logisticsState)
		acceptEvent, ok := event.(*events.GetProductsBroadcastAcceptEvent)
		if !ok {
			t.Errorf("event was not of type GetProductsBroadcastAcceptEvent")
		} else {
			if acceptEvent.ID != id {
				t.Errorf("expected id: %s - got: %s", id, acceptEvent.ID)
			}
		}

		// expect ignore
		logisticsState.Logistic.ProductsStored = 51
		event = command.ToEvent(logisticsState)
		if event != nil {
			t.Errorf("expected no event when passing a production state.")
		}

		event = command.ToEvent(productionState)
		if event != nil {
			t.Errorf("expected no event when passing a production state.")
		}
	})
}

// Parser tests
func testParseCommand(t *testing.T) {
	t.Parallel()
	t.Run("ParseGetCommand", func(t *testing.T) {
		//parser will reverse the order of the json contents
		bufferStr := "{\"asked_id\":\"23\",\"type\":\"GetDataCommand\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		command := ParseCommand(client, buffer, n)
		if command == nil {
			t.Errorf("parser error: no parser output")
		}
		commandBytes, err := json.Marshal(command)
		if err != nil {
			t.Errorf("parser error: can not unmarshal the parser output (maybe empty?)")
		}
		if (bytes.Compare(commandBytes, buffer)) != 0 {
			t.Errorf("parser - expected command: %s - got: %s", string(buffer), string(commandBytes))
		}
	})

	t.Run("ParseGetMaterialsCommand", func(t *testing.T) {
		//parser will reverse the order of the json contents
		bufferStr := "{\"amount\":\"50\",\"message\":\"Getting new materials.\",\"type\":\"GetMaterialsCommand\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		// read from server queue
		command := ParseCommand(client, buffer, n)
		if command == nil {
			t.Errorf("parser error: no parser output")
		}
		commandBytes, err := command.MarshalJSON()
		if err != nil {
			t.Errorf("parser error: can not unmarshal the queue output (maybe empty?)")
		}
		if (bytes.Compare(commandBytes, buffer)) != 0 {
			t.Errorf("parser - expected command: %s - got: %s", string(buffer), string(commandBytes))
		}
	})

	t.Run("ParseGetProductsCommand", func(t *testing.T) {
		//parser will reverse the order of the json contents
		bufferStr := "{\"amount\":\"50\",\"message\":\"Fetching new products.\",\"type\":\"GetProductsCommand\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		command := ParseCommand(client, buffer, n)
		if command == nil {
			t.Errorf("parser error: no parser output")
		}
		commandBytes, err := json.Marshal(command)
		if err != nil {
			t.Errorf("parser error: can not unmarshal the queue output (maybe empty?)")
		}
		if (bytes.Compare(commandBytes, buffer)) != 0 {
			t.Errorf("parser - expected command: %s - got: %s", string(buffer), string(commandBytes))
		}
	})
	t.Run("ParseGetMaterialsBroadcastCommand", func(t *testing.T) {
		//parser will reverse the order of the json contents
		bufferStr := "{\"id\":\"123\",\"ip\":\"10.20.30.40\",\"message\":\"Asking for new materials.\",\"type\":\"GetMaterialsBroadcastCommand\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		command := ParseCommand(client, buffer, n)
		if command == nil {
			t.Errorf("parser error: no parser output")
		}
		commandBytes, err := json.Marshal(command)
		if err != nil {
			t.Errorf("parser error: can not unmarshal the queue output (maybe empty?)")
		}
		if (bytes.Compare(commandBytes, buffer)) != 0 {
			t.Errorf("parser - expected command: %s - got: %s", string(buffer), string(commandBytes))
		}
	})
	t.Run("ParseGetProductsBroadcastCommand", func(t *testing.T) {
		//parser will reverse the order of the json contents
		bufferStr := "{\"id\":\"123\",\"ip\":\"10.20.30.40\",\"message\":\"Asking to fetch new products.\",\"type\":\"GetProductsBroadcastCommand\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		command := ParseCommand(client, buffer, n)
		if command == nil {
			t.Errorf("parser error: no parser output")
		}
		commandBytes, err := json.Marshal(command)
		if err != nil {
			t.Errorf("parser error: can not unmarshal the queue output (maybe empty?)")
		}
		if (bytes.Compare(commandBytes, buffer)) != 0 {
			t.Errorf("parser - expected command: %s - got: %s", string(buffer), string(commandBytes))
		}
	})
}

func testParseInvalidCommand(t *testing.T) {
	t.Parallel()
	/*
		t.Run("ParseGetCommandNoId", func(t *testing.T) {
			//parser will reverse the order of the json contents
			bufferStr := "{\"type\":\"GetDataCommand\"}"
			buffer := []byte(bufferStr)
			n := len(bufferStr)
			command := ParseCommand(client, buffer, n)
			if command != nil {
				t.Errorf("parser error: expected nil for invalid input")
			}
		})
	*/
	t.Run("ParseGetMaterialsCommandNoAmount", func(t *testing.T) {
		//parser will reverse the order of the json contents
		bufferStr := "{\"message\":\"Getting new materials.\",\"type\":\"GetMaterialsCommand\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		// read from server queue
		command := ParseCommand(client, buffer, n)
		if command != nil {
			t.Errorf("parser error: expected nil for invalid input")
		}
	})

	t.Run("ParseGetProductsCommandNoAmountNoMsg", func(t *testing.T) {
		//parser will reverse the order of the json contents
		bufferStr := "{\"50\"\"type\":\"GetProductsCommand\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		command := ParseCommand(client, buffer, n)
		if command != nil {
			t.Errorf("parser error: expected nil for invalid input")
		}
	})
	t.Run("ParseGetMaterialsBroadcastCommandNoIp", func(t *testing.T) {
		//parser will reverse the order of the json contents
		bufferStr := "{\"id\":\"1337\",\"message\":\"Asking for new materials.\",\"type\":\"GetMaterialsBroadcastCommand\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		command := ParseCommand(client, buffer, n)
		if command != nil {
			t.Errorf("parser error: expected nil for invalid input")
		}
	})

	t.Run("ParseGetMaterialsBroadcastCommandNoId", func(t *testing.T) {
		//parser will reverse the order of the json contents
		bufferStr := "{\"ip\":\"10.20.30.40\",\"message\":\"Asking for new materials.\",\"type\":\"GetMaterialsBroadcastCommand\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		command := ParseCommand(client, buffer, n)
		if command != nil {
			t.Errorf("parser error: expected nil for invalid input")
		}
	})

	t.Run("ParseGetProductsBroadcastCommandNoIp", func(t *testing.T) {
		//parser will reverse the order of the json contents
		bufferStr := "{\"id\":\"123\",\"message\":\"Asking to fetch new products.\",\"type\":\"GetProductsBroadcastCommand\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		command := ParseCommand(client, buffer, n)
		if command != nil {
			t.Errorf("parser error: no parser output")
		}
	})

	t.Run("ParseGetProductsBroadcastCommandNoId", func(t *testing.T) {
		//parser will reverse the order of the json contents
		bufferStr := "{\"ip\":\"10.20.30.40\",\"message\":\"Asking to fetch new products.\",\"type\":\"GetProductsBroadcastCommand\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		command := ParseCommand(client, buffer, n)
		if command != nil {
			t.Errorf("parser error: no parser output")
		}
	})
}

// Helper Functions
func helpSetGetClient(t *testing.T, command stahl4.Command) {
	command.SetClient(client)
	getClient := command.GetClient()
	if client.Id != getClient.Id {
		t.Errorf("expected id: %d - got: %d", client.Id, getClient.Id)
	}
	if client != getClient {
		t.Errorf("expected client: %p - got: %p", client, getClient)
	}
}

func helpSetGetServer(t *testing.T, command stahl4.Command) {
	command.SetServer(server)
	getServer := command.GetServer()
	if server != getServer {
		t.Errorf("expected server: %p - got: %p", server, getServer)
	}
}
