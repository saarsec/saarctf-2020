package network

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"log"
	"net"
	"stahl4/service/internal/commands"
	"stahl4/service/internal/stahl4"
	"strings"
	"sync"
	"time"
)

type Server struct {
	// mutex for thread safety
	mux sync.Mutex
	// the command queue (FIFO)
	queue chan stahl4.CommandItem
	// the client connections
	connections *[]stahl4.RemoteClient
	// the read timeout
	readTimeout time.Duration
	// the write timeout
	writeTimeout time.Duration
	// type of server we are
	typ uint8
	// our ip
	ip string
	// port
	port string
}

func CreateServer(typ uint8, port string) Server {
	// the queue for the commands that the client sends us
	// this creates a buffered channel (see https://tour.golang.org/concurrency/3)
	queue := make(chan stahl4.CommandItem, 1000)
	// slice for the storage of the client connections
	var remoteClients []stahl4.RemoteClient
	// set the timeouts for the server
	readTimeout := 60 * time.Second
	writeTimeout := 20 * time.Second
	return Server{
		queue:        queue,
		connections:  &remoteClients,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		typ:          typ,
		ip:           hostIP(),
		port:         port,
	}
}

func (server *Server) AddToQueue(command stahl4.Command) {
	select {
	// push the command onto the queue
	case server.queue <- stahl4.CommandItem{Command: command}:
		return
	default:
		log.Print("Queue already full, command could not be pushed")
	}
}

func (server *Server) ReceiveCommand() stahl4.Command {
	select {
	// get the first CommandItem that was added to the queue (FIFO)
	case cmd := <-server.queue:
		// return the Command of the CommandItem
		return cmd.Command
	// nothing on the queue
	default:
		return nil
	}
}

// the actual server function which binds to the port and handles incoming connections
func (server *Server) Listen() {
	// a counter which we use to choose id's for new clients
	// idCounter := int32(0)
	// bind to the port and listen
	l, err := net.Listen("tcp4", "0.0.0.0:"+server.port) // l, err := net.Listen("tcp", HostIP()+":21485")
	// check if there was no error while binding to the port and kill the program if there is an error
	if err != nil {
		panic(err)
	}
	// defer the closure of the socket
	defer l.Close()
	// do forever
	for {
		// accept new connections
		c, err := l.Accept()
		// if there is an error log it and continue
		if err != nil {
			log.Print(err.Error())
			continue
		}
		// handle each connection in it's own goroutine
		log.Printf("Received new connection from: %s", c.RemoteAddr().String())
		go handleClientConnection(c, server)
	}
}

func (server *Server) SendEvent(event stahl4.Event) error {
	// Convert the Event to its JSON representation
	data, err := json.Marshal(event)
	// Check if no error occurred
	if err != nil {
		// save the error into the log
		// log.Print(err)
		return err
	} else {
		// send the message to the client
		// fmt.Printf("answering client: %s with data: %s \n", event.GetClient().Ip, data)
		_, err := writeConnectionTimeout(server.writeTimeout, event.GetClient().Connection, data)
		// if there is an error
		if err != nil {
			// close the connection
			closeClientConnErr(event.GetClient(), event.GetClient().Connection, err)
			return err
		}
		return nil
	}
}

func handleClientConnection(connection net.Conn, server *Server) {
	// create buffer to read into
	buffer := make([]byte, 4096)
	var chandshake ClientHandshake
	i := 0
	for i < maxChallenges {
		// handle handshake
		i++
		challenge, err := newServerChallenge()
		if err != nil {
			log.Printf("Could not create server challenge")
			// close the connection
			_ = connection.Close()
			// let this handle the end of the loop
			return
		}
		cjson, err := json.Marshal(challenge)
		if err != nil {
			log.Printf("Could not marshal server challenge")
			// close the connection
			_ = connection.Close()
			// let this handle the end of the loop
			return
		}
		log.Print("Server: Send Server Challenge ", string(cjson))
		_, err = writeConnectionTimeout(server.writeTimeout, connection, cjson)
		if err != nil {
			log.Printf("Timeout while sending challenge to client: %s", connection.RemoteAddr().String())
			// close the connection
			_ = connection.Close()
			// let this handle the end of the loop
			return
		}
		n, err := readConnectionTimeout(server.readTimeout, connection, buffer)
		// if there is an error while reading close this connection
		if err != nil {
			log.Printf("Timeout while recv message from client: %s", connection.RemoteAddr().String())
			// close the connection
			_ = connection.Close()
			// let this handle the end of the loop
			return
		}
		log.Print("Server: Received ClientHandshake ", string(buffer[:n]))
		var data map[string]string
		// decode the outer json layer and save the object into the variable
		err = json.Unmarshal(json.RawMessage(buffer[:n]), &data)
		// if there was an error unmarshalling the message
		if err != nil {
			log.Printf("Couldnt parse message from client: %s - %s", connection.RemoteAddr().String(), err.Error())
			// close the connection
			_ = connection.Close()
			// let this handle the end of the loop
			return
		}
		err = chandshake.fromMap(data)
		if err != nil {
			log.Printf("Couldnt rebuild handshake from client: %s - %s", connection.RemoteAddr().String(), err.Error())
			// close the connection
			_ = connection.Close()
			// let this handle the end of the loop
			return
		}
		if !chandshake.verify(server, challenge.Nonce) {
			log.Printf("Invalid handshake from client: %s", connection.RemoteAddr().String())
			_ = connection.Close()
			return
		} else {
			break
		}
	}
	log.Printf("Checking if we know client: %s", connection.RemoteAddr().String())
	// check if we already know the client to save id's
	client, check := checkKnownClient(chandshake, *server.connections)
	// if we don't know the cliencreateclientt
	if !check {
		log.Printf("We do not know client: %s", connection.RemoteAddr().String())
		// get and set the client id to the next id (atomic to do it thread-safe)
		// id := atomic.AddInt32(idCounter, 1)
		// answer to the client connect such that they get to know their id and our type
		log.Printf("Sending handshake response to client: %s", connection.RemoteAddr().String())
		res, err := json.Marshal(ServerHandshake{
			Type: int(server.typ),
		})
		if err != nil {
			log.Printf("Could not generate Handshake response: %s", err.Error())
		}
		_, _ = writeConnectionTimeout(server.writeTimeout, connection, []byte(res))
		// lock the server for thread safety
		server.mux.Lock()
		// new ID handling
		idHash := sha256.New()
		//idHash.Write(x509.MarshalPKCS1PublicKey(&chandshake.PubKey))
		idHash.Write(x509.MarshalPKCS1PublicKey(&chandshake.PubKey))
		id := hex.EncodeToString(idHash.Sum(nil))
		//id := base64.StdEncoding.EncodeToString(idHash.Sum(nil))
		// append the new connection to the known connections
		*server.connections = append(*server.connections, stahl4.RemoteClient{
			Mux:        sync.Mutex{},
			Id:         id,
			Ip:         connection.RemoteAddr(),
			Connection: connection,
			PubKey:     chandshake.PubKey,
		})
		// unlock the server
		server.mux.Unlock()
		// get a reference to our new client
		client = getClient(id, *server.connections)
		// if we were not able to find our client
		if client == nil {
			// let the function return to terminate
			return
		}
		// add a log entry for the new client
		log.Print("New peer connected: " + connection.RemoteAddr().String())
	} else { // if we already know the client
		log.Printf("We do know client: %s", connection.RemoteAddr().String())
		res, err := json.Marshal(ServerHandshake{
			Type: int(server.typ),
		})
		if err != nil {
			log.Printf("Could not generate Handshake response: %s", err.Error())
		}
		// lock the client for thread safety
		client.Mux.Lock()
		// set the connection of the client to this connection
		client.Connection = connection
		_, _ = writeConnectionTimeout(server.writeTimeout, connection, []byte(res))
		// unlock the client
		client.Mux.Unlock()
	}
	// do as long as there is no new connection of this client
	for client.Connection == connection {
		// read from the connection with the timeouts that were set
		n, err := readConnectionTimeout(server.readTimeout, connection, buffer)
		// if there is an error while reading close this connection
		if err != nil {
			log.Printf("Timeout while recv message from client: %s", connection.RemoteAddr().String())
			closeClientConnErr(client, connection, err)
			// let this handle the end of the loop
			continue
		}
		// if we read more than 0 bytes
		if n > 0 {
			// log.Printf("Received message: %s \n from client: %s", buffer[:n], connection.RemoteAddr().String())
			// We received a heartbeat message
			if buffer[0] == 0xF0 {
				// create our heartbeat answer
				var pong = make([]byte, 1)
				pong[0] = 0xF0
				// answer the heartbeat with a heartbeat of ourselves
				_, err := writeConnectionTimeout(server.writeTimeout, client.Connection, pong)
				// if there is an error during the heartbeat
				if err != nil {
					// close the connection
					closeClientConnErr(client, client.Connection, err)
					// let this handle the termination of the loop
					continue
				}
			} else { // if we receive a different message
				// try to parse the received message
				command := commands.ParseCommand(client, buffer, n)
				if command != nil {
					server.AddToQueue(command)
					// log.Printf("Added command of client to queue: %s", connection.RemoteAddr().String())
				}
			}
		}
	}
}

// function to close a connection if an error occurred
func closeClientConnErr(client *stahl4.RemoteClient, connection net.Conn, err error) {
	// write the event into the log
	log.Print("Closing connection to client peer " + client.Ip.String() + " because of " + err.Error())
	// lock client for thread safety
	client.Mux.Lock()
	// close the connection
	if connection != nil {
		_ = connection.Close()
	}
	// set the connection for the client to nil to indicate it (important so that e.g. handleClientConnection will
	// terminate)
	client.Connection = nil
	// unlock the client
	client.Mux.Unlock()
}

// function to find and return the client with a given id
func getClient(id string, connections []stahl4.RemoteClient) *stahl4.RemoteClient {
	// iterate over all connections
	for _, cl := range connections {
		// if the connection id is equal to the id that we search
		if strings.Compare(cl.Id, id) == 0 {
			// return a pointer to the client
			return &cl
		}
	}
	// nothing was found so we return nil
	return nil
}

// function to check if we know the client with a given ip address
func checkKnownClient(chandshake ClientHandshake, connections []stahl4.RemoteClient) (*stahl4.RemoteClient, bool) {
	// iterate over all connections
	for _, cl := range connections {
		if cl.PubKey.E == chandshake.PubKey.E && cl.PubKey.N.Cmp(chandshake.PubKey.N) == 0 {
			// return a pointer to the client and a bool to indicate that we have found the client
			return &cl, true
		}
	}
	// no client was found, return nil and false
	return nil, false
}
