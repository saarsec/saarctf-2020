package network

import (
	"stahl4/service/internal/commands"
	"stahl4/service/internal/events"
	"stahl4/service/internal/stahl4"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net"
	"os"
	"testing"
	"time"
)

/*
var client = &stahl4.RemoteClient{
	Mux:        sync.Mutex{},
	Id:         23,
	Ip:         nil,
	Connection: nil,
}

var server = &stahl4.RemoteServer{
	Mux:        sync.Mutex{},
	Id:         42,
	Ip:         nil,
	Connection: nil,
	Timestamp:  time.Time{},
	Typ:        7,
}
*/
//port := 31337
//logistic := 1
//server := network.CreateServer(logistic, port)
// Globals

// (will be filled by testNetworkCreate and used by other tests)
var client *Client
var server *Server

/*
var rem_client = &stahl4.RemoteClient{
	Mux:        sync.Mutex{},
	Id:         23,
	Ip:         nil,
	Connection: nil,
}

var rem_server = &stahl4.RemoteServer{
	Mux:        sync.Mutex{},
	Id:         42,
	Ip:         nil,
	Connection: nil,
	Timestamp:  time.Time{},
	Typ:        7,
}
*/
//TODO: write client parser tests to events

func TestServer(t *testing.T) {
	//calling the tests
	t.Run("TestServerClientCreate", testServerClientCreate)
	t.Run("TestServerFunctions", testServerFunctions)
	//os.RemoveAll("./data/")
	os.Remove("../../server.log")
}

/*
client := network.CreateClient(port)
broadcasts := make([]stahl4.BroadcastCommand, 0)

logistic := stahl4.NewLogistic()
state.Logistic = logistic

server.Listen()
client.Start()
// Command test

//buffer := make([]byte, 1024)
*/

func testServerClientCreate(t *testing.T) {
	//We can not run parallel to testServerParse, because we create the dependencies like server and client
	serverPort := "30002"
	// type to logistic
	serverType := uint8(1)
	// private key for the client
	var privK *rsa.PrivateKey
	privK, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	t.Run("CreateClient", func(t *testing.T) {
		clientVal := CreateClient(serverPort, privK)
		if clientVal.serverPort != serverPort {
			t.Errorf("expected port: %s - got: %s", serverPort, client.serverPort)
		}
		client = &clientVal
	})
	t.Run("CreateServer", func(t *testing.T) {
		serverVal := CreateServer(serverType, serverPort)
		if serverVal.port != serverPort {
			t.Errorf("expected port: %s - got: %s", serverPort, server.port)
		}
		server = &serverVal
	})

}

func testServerFunctions(t *testing.T) {

	// for the following tests we need a listener
	go server.Listen()
	// wait 2 seconds for the listener to open
	time.Sleep(2 * time.Second)

	t.Run("AddtoQueue", func(t *testing.T) {
		var command = &commands.GetProductsCommand{
			Amount:  13,
			Message: "test",
			Client:  nil,
			Server:  nil,
		}

		//CommandJson := "{\"asked_id\":\"23\",\"type\":\"GetCommand\"}"
		//command_obj := json.Unmarshal(CommandJson)
		server.AddToQueue(command)
		//var commandRecv command
		commandRecv := server.ReceiveCommand()
		//we marshal the files for comparison and output
		cmdRecv, _ := json.Marshal(commandRecv)
		cmdExp, _ := json.Marshal(command)
		//the output will be marshalled again in Command: {<CMDRecv>}

		if bytes.Compare(cmdExp, cmdRecv) != 0 {
			t.Errorf("expected command: %s - got: %s", string(cmdExp), string(cmdRecv))
		}
	})
	t.Run("ReceiveGetProductsCommand", func(t *testing.T) {
		var command = &commands.GetProductsCommand{
			Amount:  13,
			Message: "test",
			Client:  nil,
			Server:  nil,
		}
		server.queue <- stahl4.CommandItem{Command: command}
		commandRecv := server.ReceiveCommand()
		cmdRecv, _ := json.Marshal(commandRecv)
		cmdExp, _ := json.Marshal(command)
		if bytes.Compare(cmdExp, cmdRecv) != 0 {
			t.Errorf("expected command: %s - got: %s", string(cmdExp), string(cmdRecv))
		}
	})

	t.Run("ReceiveGetMaterialsCommand", func(t *testing.T) {
		var command = &commands.GetMaterialsCommand{
			Amount:  13,
			Message: "test",
			Client:  nil,
			Server:  nil,
		}
		server.queue <- stahl4.CommandItem{Command: command}
		commandRecv := server.ReceiveCommand()
		cmdRecv, _ := json.Marshal(commandRecv)
		cmdExp, _ := json.Marshal(command)
		if bytes.Compare(cmdExp, cmdRecv) != 0 {
			t.Errorf("expected command: %s - got: %s", string(cmdExp), string(cmdRecv))
		}
	})

	t.Run("ReceiveGetCommand", func(t *testing.T) {
		var command = &commands.GetDataCommand{
			AskedId: 13,
			Client:  nil,
			Server:  nil,
		}
		server.queue <- stahl4.CommandItem{Command: command}
		commandRecv := server.ReceiveCommand()
		cmdRecv, _ := json.Marshal(commandRecv)
		cmdExp, _ := json.Marshal(command)
		if bytes.Compare(cmdExp, cmdRecv) != 0 {
			t.Errorf("expected command: %s - got: %s", string(cmdExp), string(cmdRecv))
		}
	})

	t.Run("ReceiveGetMaterialsBroadcast", func(t *testing.T) {
		var command = &commands.GetMaterialsBroadcastCommand{
			ID:     "1337",
			IP:     "10.31.33.7",
			Client: nil,
			Server: nil,
		}
		server.queue <- stahl4.CommandItem{Command: command}
		commandRecv := server.ReceiveCommand()
		cmdRecv, _ := json.Marshal(commandRecv)
		cmdExp, _ := json.Marshal(command)
		if bytes.Compare(cmdExp, cmdRecv) != 0 {
			t.Errorf("expected command: %s - got: %s", string(cmdExp), string(cmdRecv))
		}
	})

	t.Run("ReceiveGetProductsBroadcast", func(t *testing.T) {
		var command = &commands.GetProductsBroadcastCommand{
			ID:     "1337",
			IP:     "10.31.33.7",
			Client: nil,
			Server: nil,
		}
		server.queue <- stahl4.CommandItem{Command: command}
		commandRecv := server.ReceiveCommand()
		cmdRecv, _ := json.Marshal(commandRecv)
		cmdExp, _ := json.Marshal(command)
		if bytes.Compare(cmdExp, cmdRecv) != 0 {
			t.Errorf("expected command: %s - got: %s", string(cmdExp), string(cmdRecv))
		}
	})

	t.Run("Listen_ReceiveCommand", func(t *testing.T) {
		//connection, err := net.DialTimeout("tcp",hostIP()+":"+rem_server.Ip.String(),client.readTimeout)
		remServconnection := connectServer(client, hostIP())
		if server.connections == nil {
			t.Errorf("Server connection is empty after client connect!")
		}
		if remServconnection == nil {
			t.Errorf("error calling connectServer - hostip: %s", hostIP())
		}
		cmdMarshal := "{\"amount\":\"50\",\"message\":\"Fetching new products.\",\"type\":\"GetProductsCommand\"}"
		cmdMarshalBytes := []byte(cmdMarshal)
		var command = &commands.GetProductsCommand{
			Amount:  50,
			Message: "test",
			Client:  nil,
			Server:  remServconnection,
		}
		//send to server
		client.SendLogisticCommand(command)
		time.Sleep(1 * time.Second)
		cmdRecv := server.ReceiveCommand()
		//var cmdExp commands.GetCommand
		//err := json.Unmarshal([]byte(cmdMarshal),&cmdExp)
		cmdRecvBytes, err := json.Marshal(cmdRecv)
		t.Logf("received command is: %s", cmdRecvBytes)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if bytes.Compare(cmdRecvBytes, cmdMarshalBytes) != 0 {
			t.Errorf("expected: %s - got: %s", string(cmdMarshalBytes), string(cmdRecvBytes))
		}
	})

	t.Run("SendEvent", func(t *testing.T) {
		remServconnection := connectServer(client, hostIP())
		var evt = &events.GetEvent{
			Message: "Test",
			Client:  &(*server.connections)[0],
			Server:  remServconnection,
		}
		evtExpect := "{\"message\":\"Test\",\"type\":\"GetEvent\"}"
		evtExpectBytes := []byte(evtExpect)
		server.SendEvent(evt)
		time.Sleep(1 * time.Second)
		evtRecv := client.ReceiveEvent()
		evtRecvMshl, err := json.Marshal(evtRecv)
		if err != nil {
			t.Errorf("there was an error when calling unmarshal: %s", err.Error())
		}
		if bytes.Compare(evtRecvMshl, evtExpectBytes) != 0 {
			t.Errorf("expected: %s- got: %s", string(evtExpectBytes), string(evtRecvMshl))
		}

	})

	t.Run("GetClient", func(t *testing.T) {

		// test a valid value -- the tests before added a single client, which therefore has Id 1
		expectId := int32(1)
		answer := getClient(expectId, *server.connections)
		// we expect nil, because of the invalid id
		if answer == nil {
			t.Errorf("returned nil instead of ClientID")
		}
		if answer.Id != int(expectId) {
			t.Errorf("expected Id: %s - got Id: %s", string(expectId), string(answer.Id))
		}

		// test an invalid value
		invalidId := int32(1337)
		answer = getClient(invalidId, *server.connections)
		if answer != nil {
			t.Errorf("expected nil as answer for invalid value, but got something else back")
		}
	})

	//TODO:
	t.Run("KeepAlive", func(t *testing.T) {
		server.readTimeout = 5 * time.Second // we decrease the server read timeout for our tests
		buffer := make([]byte, 1024)
		// open connection to server manually
		connection, err := net.DialTimeout("tcp", hostIP()+":"+client.serverPort, client.readTimeout)
		if err != nil {
			t.Errorf("Error trying to connect to server.")
		}
		// handshake
		// read challenge
		n, err := readConnectionTimeout(client.readTimeout, connection, buffer)

		var challenge ServerChallenge
		err = json.Unmarshal(buffer[:n], &challenge)
		//log.Print("DBG - got challenge: " + challenge.Nonce)
		if err != nil {
			t.Errorf("Error on parsing server challenge")
		}
		handshake, err := newClientHandshake(hostIP(), challenge.Nonce)
		if err != nil {
			t.Errorf("Error on new handshake")
		}
		err = handshake.sign(client.key)
		if err != nil {
			t.Errorf("Error signing handshake connection to server:")
		}
		hjson, err := json.Marshal(handshake)
		if err != nil {
			t.Errorf("Error marshalling handshake")
		}
		// send handshake
		_, err = writeConnectionTimeout(client.writeTimeout, connection, hjson)

		// read server id / type from connection
		_, err = readConnectionTimeout(client.readTimeout, connection, buffer)

		// end handshake

		time.Sleep(1 * time.Second)
		// write the heartbeat
		var ping = make([]byte, 1)
		ping[0] = 0xF0
		_, err = writeConnectionTimeout(client.writeTimeout, connection, ping)
		time.Sleep(10 * time.Second)
		n, err = readConnectionTimeout(client.readTimeout, connection, buffer)
		//print(n)
		// TODO: seems like we need to read nonce first and append the connection(not sure about that appending is neccessary)
		if n > 0 {
			if buffer[0] != 0xF0 {
				t.Errorf("Expected Keepalive - got: %s", string(buffer))
			}
		} else {
			t.Errorf("Did not receive KeepAlive Packet after sending Ping")
		}

		// negative test
		time.Sleep(6 * time.Second) // server timeout is 5seconds, so this should timeout
		_, err = writeConnectionTimeout(client.writeTimeout, connection, ping)
		time.Sleep(2 * time.Second) // wait for server to eventually answer us (should not, of course)
		n, err = readConnectionTimeout(client.readTimeout, connection, buffer)
		if err == nil {
			t.Errorf("Expected Error after timeouting connection")
		}
	})

	//TODO: Test checkKnownClient when certs are implemented
}
