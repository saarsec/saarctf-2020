package network

import (
	"stahl4/service/internal/commands"
	"stahl4/service/internal/events"
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

var client_cl *Client
var server_cl *Server
func TestClient(t *testing.T) {
	//calling the tests
	t.Run("TestServerClientCreateCl", testServerClientCreateCl)
	t.Run("TestClientFunctions", testClientFunctions)
	// os.RemoveAll("./data/")
}

func testServerClientCreateCl(t *testing.T) {
	//We can not run parallel to testServerParse, because we create the dependencies like server and client
	serverPort := "30003"
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
		client_cl = &clientVal
	})
	t.Run("CreateServer", func(t *testing.T) {
		serverVal := CreateServer(serverType, serverPort)
		if serverVal.port != serverPort {
			t.Errorf("expected port: %s - got: %s", serverPort, server.port)
		}
		server_cl = &serverVal
		go server_cl.Listen()
	})
}

func testClientFunctions(t *testing.T) {
	t.Run("AddToQueue", func(t *testing.T) {
		var ev01 = &events.GetEvent{
			Message: "Get_Event_test",
			Client:  nil,
			Server:  nil,
		}
		client_cl.AddToQueue(ev01)
		var ev01Recv = client_cl.ReceiveEvent()
		if ev01Recv == nil {
			t.Errorf("no event returned after AddToQueue + ReceiveEvent")
		}
	})

	t.Run("SendGetCommand", func(t *testing.T) {
		var remServer = connectServer(client_cl, hostIP())
		if remServer == nil {
			t.Errorf("Could not connect to server")
			t.FailNow()
		}
		var cmd = commands.GetDataCommand{
			AskedId: 1337,
			Client:  nil,
			Server:  remServer,
		}
		var err = client_cl.SendGetCommand(cmd)
		if err != nil {
			t.Errorf("Client error sending getCommand")
		}
	})
}
