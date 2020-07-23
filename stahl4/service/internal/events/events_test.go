package events

import (
	"stahl4/service/internal/stahl4"
	"bytes"
	"encoding/json"
	"os"
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

func TestEvents(t *testing.T) {
	// setup
	t.Run("EventTestGroup", func(t *testing.T) {
		t.Run("TestGetEvent", testGetEvent)
		t.Run("BusyEvent", testBusyEvent)
		t.Run("TestNoProductSpaceEvent", testNoProductSpaceEvent)
		t.Run("TestNoMaterialsEvent", testNoMaterialsEvent)
		t.Run("TestDeliveredMaterialsEvent", testDeliveredMaterialsEvent)
		t.Run("TestFetchedProductsEvent", testFetchedProductsEvent)
		t.Run("TestGetMaterialsBroadcastAcceptEvent", testGetMaterialsBroadcastAcceptEvent)
		t.Run("TestGetProductsBroadcastAcceptEvent", testGetProductsBroadcastAcceptEvent)
	})
	t.Run("EventParserTests", func(t *testing.T) {
		t.Run("TestParseEvent", testParseEvent)
		t.Run("TestParseInvalidEvent", testParseInvalidEvent)
	})
	// teardown
	os.RemoveAll("./data")
}

// Event Tests
func testGetProductsBroadcastAcceptEvent(t *testing.T) {
	t.Parallel()
	var event = &GetProductsBroadcastAcceptEvent{
		Message: "bla",
		ID:      "1234",
		Client:  nil,
		Server:  nil,
	}
	marshal := "{\"id\":\"1234\",\"message\":\"bla\",\"type\":\"GetProductsBroadcastAcceptEvent\"}"
	t.Run("SetGetClient", func(t *testing.T) {
		helpSetGetClient(t, event)
	})
	t.Run("SetGetServer", func(t *testing.T) {
		helpSetGetServer(t, event)
	})
	t.Run("Marshal", func(t *testing.T) {
		helpMarshal(t, event, marshal)
	})
	t.Run("Unmarshal", func(t *testing.T) {
		var object GetProductsBroadcastAcceptEvent
		err := json.Unmarshal([]byte(marshal), &object)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if object.ID != event.ID {
			t.Errorf("expected id: %s - got: %s", event.ID, object.ID)
		}
		if object.Message != event.Message {
			t.Errorf("expected message: %s - got: %s", event.Message, object.Message)
		}
	})
	t.Run("UpdateState", func(t *testing.T) {
		// nothing happens here
	})
}

func testGetMaterialsBroadcastAcceptEvent(t *testing.T) {
	t.Parallel()
	var event = &GetMaterialsBroadcastAcceptEvent{
		Message: "bla",
		ID:      "1234",
		Client:  nil,
		Server:  nil,
	}
	marshal := "{\"id\":\"1234\",\"message\":\"bla\",\"type\":\"GetMaterialsBroadcastAcceptEvent\"}"
	t.Run("SetGetClient", func(t *testing.T) {
		helpSetGetClient(t, event)
	})
	t.Run("SetGetServer", func(t *testing.T) {
		helpSetGetServer(t, event)
	})
	t.Run("Marshal", func(t *testing.T) {
		helpMarshal(t, event, marshal)
	})
	t.Run("Unmarshal", func(t *testing.T) {
		var object GetProductsBroadcastAcceptEvent
		err := json.Unmarshal([]byte(marshal), &object)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if object.ID != event.ID {
			t.Errorf("expected id: %s - got: %s", event.ID, object.ID)
		}
		if object.Message != event.Message {
			t.Errorf("expected message: %s - got: %s", event.Message, object.Message)
		}
	})
	t.Run("UpdateState", func(t *testing.T) {
		// nothing happens here
	})
}

func testFetchedProductsEvent(t *testing.T) {
	t.Parallel()
	var productionState = &stahl4.State{
		Production: &stahl4.Production{
			ProductCapacity:    300,
			MaterialCapacity:   300,
			ProductAmount:      100,
			MaterialAmount:     0,
			AskedForMaterials:  true,
			AskedToGetProducts: true,
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
	var event = &FetchedProductsEvent{
		Message: "keks",
		Amount:  10,
		Client:  nil,
		Server:  nil,
	}
	marshal := "{\"amount\":\"10\",\"message\":\"keks\",\"type\":\"FetchedProductsEvent\"}"
	t.Run("SetGetClient", func(t *testing.T) {
		helpSetGetClient(t, event)
	})
	t.Run("SetGetServer", func(t *testing.T) {
		helpSetGetServer(t, event)
	})
	t.Run("Marshal", func(t *testing.T) {
		helpMarshal(t, event, marshal)
	})
	t.Run("Unmarshal", func(t *testing.T) {
		var object FetchedProductsEvent
		err := json.Unmarshal([]byte(marshal), &object)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if object.Message != event.Message {
			t.Errorf("expected message: %s - got: %s", event.Message, object.Message)
		}
		if object.Amount != event.Amount {
			t.Errorf("expected amount: %d - got: %d", event.Amount, object.Amount)
		}
	})
	t.Run("UpdateState", func(t *testing.T) {
		event.UpdateState(productionState)
		if productionState.Production.ProductAmount != 90 {
			t.Errorf("expected Products: %d - got: %d", 90, event.Amount)
		}
		if productionState.Production.AskedToGetProducts {
			t.Errorf("expected AskedToGetProducts: false - got: true")
		}
		event.Amount = 100
		event.UpdateState(productionState)
		if productionState.Production.ProductAmount != 90 {
			t.Errorf("expected Products: %d - got: %d", 90, event.Amount)
		}
		if productionState.Production.AskedToGetProducts {
			t.Errorf("expected AskedToGetProducts: false - got: true")
		}
		// this should not do anything
		event.UpdateState(logisticsState)
		if productionState.Production.ProductAmount != 90 {
			t.Errorf("expected Products: %d - got: %d", 90, event.Amount)
		}
		if productionState.Production.AskedToGetProducts {
			t.Errorf("expected AskedToGetProducts: false - got: true")
		}
	})
}

func testDeliveredMaterialsEvent(t *testing.T) {
	t.Parallel()
	var productionState = &stahl4.State{
		Production: &stahl4.Production{
			ProductCapacity:    300,
			MaterialCapacity:   300,
			ProductAmount:      100,
			MaterialAmount:     0,
			AskedForMaterials:  true,
			AskedToGetProducts: true,
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
	var event = &DeliveredMaterialsEvent{
		Message: "keks",
		Amount:  100,
		Client:  nil,
		Server:  nil,
	}
	marshal := "{\"amount\":\"100\",\"message\":\"keks\",\"type\":\"DeliveredMaterialsEvent\"}"
	t.Run("SetGetClient", func(t *testing.T) {
		helpSetGetClient(t, event)
	})
	t.Run("SetGetServer", func(t *testing.T) {
		helpSetGetServer(t, event)
	})
	t.Run("Marshal", func(t *testing.T) {
		helpMarshal(t, event, marshal)
	})
	t.Run("Unmarshal", func(t *testing.T) {
		var object DeliveredMaterialsEvent
		err := json.Unmarshal([]byte(marshal), &object)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if object.Message != event.Message {
			t.Errorf("expected message: %s - got: %s", event.Message, object.Message)
		}
		if object.Amount != event.Amount {
			t.Errorf("expected amount: %d - got: %d", event.Amount, object.Amount)
		}
	})
	t.Run("UpdateState", func(t *testing.T) {
		event.UpdateState(productionState)
		if productionState.Production.MaterialAmount != 100 {
			t.Errorf("expected Materials: %d - got: %d", 100, productionState.Production.MaterialAmount)
		}
		if productionState.Production.AskedForMaterials {
			t.Errorf("expected AskedForMaterials: false - got: true")
		}
		event.Amount = 255
		event.UpdateState(productionState)
		if productionState.Production.MaterialAmount != 100 {
			t.Errorf("expected Materials: %d - got: %d", 100, productionState.Production.MaterialAmount)
		}
		if productionState.Production.AskedForMaterials {
			t.Errorf("expected AskedForMaterials: false - got: true")
		}
		// this should not do anything
		event.UpdateState(logisticsState)
		if productionState.Production.MaterialAmount != 100 {
			t.Errorf("expected Materials: %d - got: %d", 100, productionState.Production.MaterialAmount)
		}
		if productionState.Production.AskedForMaterials {
			t.Errorf("expected AskedForMaterials: false - got: true")
		}
	})
}

func testNoMaterialsEvent(t *testing.T) {
	t.Parallel()
	var productionState = &stahl4.State{
		Production: &stahl4.Production{
			ProductCapacity:    300,
			MaterialCapacity:   300,
			ProductAmount:      100,
			MaterialAmount:     0,
			AskedForMaterials:  true,
			AskedToGetProducts: true,
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
	var event = &NoMaterialsEvent{
		Message: "keks",
		Client:  nil,
		Server:  nil,
	}
	marshal := "{\"message\":\"keks\",\"type\":\"NoMaterialsEvent\"}"
	t.Run("SetGetClient", func(t *testing.T) {
		helpSetGetClient(t, event)
	})
	t.Run("SetGetServer", func(t *testing.T) {
		helpSetGetServer(t, event)
	})
	t.Run("Marshal", func(t *testing.T) {
		helpMarshal(t, event, marshal)
	})
	t.Run("Unmarshal", func(t *testing.T) {
		var object NoMaterialsEvent
		err := json.Unmarshal([]byte(marshal), &object)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if object.Message != event.Message {
			t.Errorf("expected message: %s - got: %s", event.Message, object.Message)
		}
	})
	t.Run("UpdateState", func(t *testing.T) {
		// this should set AskedForMaterials to false
		event.UpdateState(productionState)
		if productionState.Production.AskedForMaterials {
			t.Errorf("expected AskedForMaterials: false - got: true")
		}
		// this should not change the state
		event.UpdateState(productionState)
		if productionState.Production.AskedForMaterials {
			t.Errorf("expected AskedForMaterials: false - got: true")
		}
		// expect nothing to happen because we are logistic
		event.UpdateState(logisticsState)
	})
}

func testNoProductSpaceEvent(t *testing.T) {
	t.Parallel()
	var productionState = &stahl4.State{
		Production: &stahl4.Production{
			ProductCapacity:    300,
			MaterialCapacity:   300,
			ProductAmount:      100,
			MaterialAmount:     0,
			AskedForMaterials:  true,
			AskedToGetProducts: true,
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
	var event = &NoProductSpaceEvent{
		Message: "keks",
		Client:  nil,
		Server:  nil,
	}
	marshal := "{\"message\":\"keks\",\"type\":\"NoProductSpaceEvent\"}"
	t.Run("SetGetClient", func(t *testing.T) {
		helpSetGetClient(t, event)
	})
	t.Run("SetGetServer", func(t *testing.T) {
		helpSetGetServer(t, event)
	})
	t.Run("Marshal", func(t *testing.T) {
		helpMarshal(t, event, marshal)
	})
	t.Run("Unmarshal", func(t *testing.T) {
		var object NoProductSpaceEvent
		err := json.Unmarshal([]byte(marshal), &object)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if object.Message != event.Message {
			t.Errorf("expected message: %s - got: %s", event.Message, object.Message)
		}
	})
	t.Run("UpdateState", func(t *testing.T) {
		// this should set AskedForProducts to false
		event.UpdateState(productionState)
		if productionState.Production.AskedToGetProducts {
			t.Errorf("expected AskedForProducts: false - got: true")
		}
		// this should not change the state
		event.UpdateState(productionState)
		if productionState.Production.AskedToGetProducts {
			t.Errorf("expected AskedForProducts: false - got: true")
		}
		// expect nothing to happen because we are logistic
		event.UpdateState(logisticsState)
	})
}

func testBusyEvent(t *testing.T) {
	t.Parallel()
	var productionState = &stahl4.State{
		Production: &stahl4.Production{
			ProductCapacity:    300,
			MaterialCapacity:   300,
			ProductAmount:      100,
			MaterialAmount:     0,
			AskedForMaterials:  true,
			AskedToGetProducts: true,
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
	var event = &BusyEvent{
		Message: "keks",
		Option:  1,
		Client:  nil,
		Server:  nil,
	}
	marshal := "{\"message\":\"keks\",\"option\":\"1\",\"type\":\"BusyEvent\"}"
	t.Run("SetGetClient", func(t *testing.T) {
		helpSetGetClient(t, event)
	})
	t.Run("SetGetServer", func(t *testing.T) {
		helpSetGetServer(t, event)
	})
	t.Run("Marshal", func(t *testing.T) {
		helpMarshal(t, event, marshal)
	})
	t.Run("Unmarshal", func(t *testing.T) {
		var object BusyEvent
		err := json.Unmarshal([]byte(marshal), &object)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if object.Message != event.Message {
			t.Errorf("expected message: %s - got: %s", event.Message, object.Message)
		}
	})
	t.Run("UpdateState", func(t *testing.T) {
		// this should set AskedForProducts to false
		event.UpdateState(productionState)
		if productionState.Production.AskedToGetProducts {
			t.Errorf("expected AskedForProducts: false - got: true")
		}
		// this should not change the state
		event.UpdateState(productionState)
		if productionState.Production.AskedToGetProducts {
			t.Errorf("expected AskedForProducts: false - got: true")
		}
		// expect nothing to happen because we are logistic
		event.UpdateState(logisticsState)
		// set option to 2 such that AskedForMaterials is indicated
		event.Option = 2
		// this should set AskedForMaterials to false
		event.UpdateState(productionState)
		if productionState.Production.AskedForMaterials {
			t.Errorf("expected AskedForMaterials: false - got: true")
		}
		// this should not change the state
		event.UpdateState(productionState)
		if productionState.Production.AskedForMaterials {
			t.Errorf("expected AskedForMaterials: false - got: true")
		}
		// expect nothing to happen because we are logistic
		event.UpdateState(logisticsState)
	})
}

func testGetEvent(t *testing.T) {
	t.Parallel()
	var event = &GetEvent{
		Message: "keks",
		Client:  nil,
		Server:  nil,
	}
	marshal := "{\"message\":\"keks\",\"type\":\"GetEvent\"}"
	t.Run("SetGetClient", func(t *testing.T) {
		helpSetGetClient(t, event)
	})
	t.Run("SetGetServer", func(t *testing.T) {
		helpSetGetServer(t, event)
	})
	t.Run("Marshal", func(t *testing.T) {
		helpMarshal(t, event, marshal)
	})
	t.Run("Unmarshal", func(t *testing.T) {
		var object GetProductsBroadcastAcceptEvent
		err := json.Unmarshal([]byte(marshal), &object)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if object.Message != event.Message {
			t.Errorf("expected id: %s - got: %s", event.Message, object.Message)
		}
	})
	t.Run("UpdateState", func(t *testing.T) {
		// nothing happens here
	})
}

// Helper Functions
func helpSetGetClient(t *testing.T, command stahl4.Event) {
	command.SetClient(client)
	getClient := command.GetClient()
	if client.Id != getClient.Id {
		t.Errorf("expected id: %d - got: %d", client.Id, getClient.Id)
	}
	if client != getClient {
		t.Errorf("expected client: %p - got: %p", client, getClient)
	}
}

func helpSetGetServer(t *testing.T, command stahl4.Event) {
	command.SetServer(server)
	getServer := command.GetServer()
	if server != getServer {
		t.Errorf("expected server: %p - got: %p", server, getServer)
	}
}

func helpMarshal(t *testing.T, command stahl4.Event, input string) {
	jsonText, err := json.Marshal(command)
	if err != nil {
		t.Errorf("there was an error when calling marshal: %s", err.Error())
	}
	if string(jsonText) != input {
		t.Errorf("expected json: %s - got: %s", input, string(jsonText))
	}
}

func testParseEvent(t *testing.T) {
	t.Parallel()
	t.Run("ParseGetEvent", func(t *testing.T) {
		bufferStr := "{\"message\":\"23\",\"type\":\"GetEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		if server == nil {
			print("server is nil")
		}
		event := ParseEvent(server, buffer, n)
		if event == nil {
			t.Errorf("parser error: no parser output")
		}
		eventBytes, err := json.Marshal(event)
		if err != nil {
			t.Errorf("parser error: can not unmarshal the parser output (maybe empty?)")
		}
		if (bytes.Compare(eventBytes, buffer)) != 0 {
			t.Errorf("parser - expected event: %s - got: %s", string(buffer), string(eventBytes))
		}
	})

	t.Run("ParseBusyEvent", func(t *testing.T) {
		bufferStr := "{\"message\":\"test\",\"option\":\"1\",\"type\":\"BusyEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event == nil {
			t.Errorf("parser error: no parser output")
		}
		eventBytes, err := json.Marshal(event)
		if err != nil {
			t.Errorf("parser error: can not unmarshal the parser output (maybe empty?)")
		}
		if (bytes.Compare(eventBytes, buffer)) != 0 {
			t.Errorf("parser - expected event: %s - got: %s", string(buffer), string(eventBytes))
		}
	})

	t.Run("ParseNoProductSpaceEvent", func(t *testing.T) {
		bufferStr := "{\"message\":\"test\",\"type\":\"NoProductSpaceEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event == nil {
			t.Errorf("parser error: no parser output")
		}
		eventBytes, err := json.Marshal(event)
		if err != nil {
			t.Errorf("parser error: can not unmarshal the parser output (maybe empty?)")
		}
		if (bytes.Compare(eventBytes, buffer)) != 0 {
			t.Errorf("parser - expected event: %s - got: %s", string(buffer), string(eventBytes))
		}
	})

	t.Run("ParseNoMaterialsEvent", func(t *testing.T) {
		bufferStr := "{\"message\":\"test\",\"type\":\"NoMaterialsEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event == nil {
			t.Errorf("parser error: no parser output")
		}
		eventBytes, err := json.Marshal(event)
		if err != nil {
			t.Errorf("parser error: can not unmarshal the parser output (maybe empty?)")
		}
		if (bytes.Compare(eventBytes, buffer)) != 0 {
			t.Errorf("parser - expected event: %s - got: %s", string(buffer), string(eventBytes))
		}
	})

	t.Run("ParseGetMaterialsBroadcastAcceptEvent", func(t *testing.T) {
		bufferStr := "{\"id\":\"1337\",\"message\":\"test\",\"type\":\"GetMaterialsBroadcastAcceptEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event == nil {
			t.Errorf("parser error: no parser output")
		}
		eventBytes, err := json.Marshal(event)
		if err != nil {
			t.Errorf("parser error: can not unmarshal the parser output (maybe empty?)")
		}
		if (bytes.Compare(eventBytes, buffer)) != 0 {
			t.Errorf("parser - expected event: %s - got: %s", string(buffer), string(eventBytes))
		}
	})

	t.Run("ParseGetProductsBroadcastAcceptEvent", func(t *testing.T) {
		bufferStr := "{\"id\":\"1337\",\"message\":\"test\",\"type\":\"GetProductsBroadcastAcceptEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event == nil {
			t.Errorf("parser error: no parser output")
		}
		eventBytes, err := json.Marshal(event)
		if err != nil {
			t.Errorf("parser error: can not unmarshal the parser output (maybe empty?)")
		}
		if (bytes.Compare(eventBytes, buffer)) != 0 {
			t.Errorf("parser - expected event: %s - got: %s", string(buffer), string(eventBytes))
		}
	})
	t.Run("ParseDeliveredMaterialsEvent", func(t *testing.T) {
		bufferStr := "{\"amount\":\"50\",\"message\":\"test\",\"type\":\"DeliveredMaterialsEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event == nil {
			t.Errorf("parser error: no parser output")
		}
		eventBytes, err := json.Marshal(event)
		if err != nil {
			t.Errorf("parser error: can not unmarshal the parser output (maybe empty?)")
		}
		if (bytes.Compare(eventBytes, buffer)) != 0 {
			t.Errorf("parser - expected event: %s - got: %s", string(buffer), string(eventBytes))
		}
	})

	t.Run("ParseFetchedProductsEvent", func(t *testing.T) {
		bufferStr := "{\"amount\":\"50\",\"message\":\"test\",\"type\":\"FetchedProductsEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event == nil {
			t.Errorf("parser error: no parser output")
		}
		eventBytes, err := json.Marshal(event)
		if err != nil {
			t.Errorf("parser error: can not unmarshal the parser output (maybe empty?)")
		}
		if (bytes.Compare(eventBytes, buffer)) != 0 {
			t.Errorf("parser - expected event: %s - got: %s", string(buffer), string(eventBytes))
		}
	})

}

func testParseInvalidEvent(t *testing.T) {
	t.Parallel()
	t.Run("ParseNoTypeSet", func(t *testing.T) {
		bufferStr := "{\"message\":\"23\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		if server == nil {
			print("server is nil")
		}
		event := ParseEvent(server, buffer, n)
		if event != nil {
			t.Errorf("parser error: expected nil for invalid input")
		}
	})

	t.Run("ParseGetEventNoMsgSet", func(t *testing.T) {
		bufferStr := "{\"type\":\"GetEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		if server == nil {
			print("server is nil")
		}
		event := ParseEvent(server, buffer, n)
		if event != nil {
			t.Errorf("parser error: expected nil for invalid input")
		}
	})

	t.Run("ParseBusyEventNoOptionSet", func(t *testing.T) {
		bufferStr := "{\"message\":\"test\",\"type\":\"BusyEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event != nil {
			t.Errorf("parser error: expected nil for invalid input")
		}
	})

	t.Run("ParseNoProductSpaceEventNoMsgSet", func(t *testing.T) {
		bufferStr := "{\"type\":\"NoProductSpaceEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event != nil {
			t.Errorf("parser error: no parser output")
		}
	})

	t.Run("ParseNoMaterialsEventNoMsgSet", func(t *testing.T) {
		bufferStr := "{\"type\":\"NoMaterialsEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event != nil {
			t.Errorf("parser error: expected nil for invalid input")
		}
	})

	t.Run("ParseGetMaterialsBroadcastAcceptEventNoIDNoMsgGiven", func(t *testing.T) {
		bufferStr := "{\"type\":\"GetMaterialsBroadcastAcceptEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event != nil {
			t.Errorf("parser error: expected nil for invalid input")
		}
	})

	t.Run("ParseGetProductsBroadcastAcceptEventNoIdNoMsgGiven", func(t *testing.T) {
		bufferStr := "{\"type\":\"GetProductsBroadcastAcceptEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event != nil {
			t.Errorf("parser error: expected nil for invalid input")
		}
	})

	t.Run("ParseDeliveredMaterialsEventNoMsgGiven", func(t *testing.T) {
		bufferStr := "{\"amount\":\"50\",\"type\":\"DeliveredMaterialsEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event != nil {
			t.Errorf("parser error: expected nil for invalid input")
		}
	})

	t.Run("ParseFetchedProductsEventNoAmountGiven", func(t *testing.T) {
		bufferStr := "{\"message\":\"test\",\"type\":\"FetchedProductsEvent\"}"
		buffer := []byte(bufferStr)
		n := len(bufferStr)
		event := ParseEvent(server, buffer, n)
		if event != nil {
			t.Errorf("parser error: expected nil for invalid input")
		}
	})

}
