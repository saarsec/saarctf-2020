package network

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"stahl4/service/internal/commands"
	"stahl4/service/internal/events"
	"stahl4/service/internal/stahl4"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	// mutex for thread safety
	mux sync.Mutex
	// the event queue (FIFO)
	queue chan stahl4.EventItem
	// the client connections
	connections *[]stahl4.RemoteServer
	// the read timeout
	readTimeout time.Duration
	// the write timeout
	writeTimeout time.Duration
	// port
	serverPort string
	// our ip
	Ip string
	// private key
	key *rsa.PrivateKey
}

func CreateClient(port string, key *rsa.PrivateKey) Client {
	// the queue for the events that the servers send us
	// this creates a buffered channel (see https://tour.golang.org/concurrency/3)
	queue := make(chan stahl4.EventItem, 1000)
	// slice for the storage of the server connections
	var remoteClients []stahl4.RemoteServer
	// set the timeouts for the client
	// set this lower than on the server because we initiate the heartbeats and therefore do not always expect a message
	readTimeout := 60 * time.Second
	writeTimeout := 20 * time.Second
	return Client{
		queue:        queue,
		connections:  &remoteClients,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		serverPort:   port,
		Ip:           hostIP(),
		key:          key,
	}
}

func (client *Client) Start() {
	log.Print("starting connectPeers function")
	connectPeers(client)
}

func (client *Client) AddToQueue(event stahl4.Event) {
	select {
	// push the command onto the queue
	case client.queue <- stahl4.EventItem{Event: event}:
		return
	default:
		log.Print("Queue already full, command could not be pushed")
	}
}

func (client *Client) ReceiveEvent() stahl4.Event {
	select {
	// get the first CommandItem that was added to the queue (FIFO)
	case event := <-client.queue:
		// return the Command of the CommandItem
		return event.Event
		// nothing on the queue
	default:
		return nil
	}
}

func (client *Client) SendGetCommand(command commands.GetDataCommand) error {
	for command.GetServer() == nil {
		var random int
		if len(*client.connections) > 0 {
			random = rand.Int() % len(*client.connections)
		} else {
			return errors.New("no connection to choose from")
		}
		remoteServer := getRemoteServer(random, *client)
		command.SetServer(remoteServer)
	}
	// Convert the Event to its JSON representation
	data, err := json.Marshal(command)
	// Check if no error occurred
	if err != nil {
		// save the error into the log
		// log.Print(err)
		return err
	} else {
		// send the message to the server
		err := sendMessage(&command, client, data)
		// if there is an error
		if err != nil {
			// close the connection
			closeServerConnErr(client, command.GetServer(), command.GetServer().Connection, err)
			return err
		}
		log.Print("Send GetCommand")
		// set the timestamp the thread safe way
		remoteServer := command.GetServer()
		remoteServer.Mux.Lock()
		remoteServer.Timestamp = time.Now()
		remoteServer.Mux.Unlock()
		return nil
	}
}

func (client *Client) SendLogisticCommand(command stahl4.Command) error {
	if command.GetServer() == nil {
		if len(*client.connections) != 0 {
			counter := len(*client.connections) * 2
			for command.GetServer() == nil && counter > 0 {
				var random int
				if len(*client.connections) > 0 {
					random = rand.Int() % len(*client.connections)
				} else {
					return errors.New("no connection to choose from")
				}
				remoteServer := getRemoteServer(random, *client)
				if remoteServer.Typ == 1 {
					command.SetServer(remoteServer)
				}
				counter--
			}
			if counter <= 0 {
				return errors.New("no connection was found")
			}
		} else {
			return errors.New("there are no connections to choose from")
		}
	}
	// Convert the Event to its JSON representation
	data, err := json.Marshal(command)
	// Check if no error occurred
	if err != nil {
		// save the error into the log
		// log.Print(err)
		return err
	} else {
		// send the message to the server
		err := sendMessage(command, client, data)
		// if there is an error
		if err != nil {
			// close the connection
			closeServerConnErr(client, command.GetServer(), command.GetServer().Connection, err)
			return err
		}
		log.Print("Send LogisticCommand")
		// set the timestamp the thread safe way
		remoteServer := command.GetServer()
		remoteServer.Mux.Lock()
		remoteServer.Timestamp = time.Now()
		remoteServer.Mux.Unlock()
		return nil
	}
}

func (client *Client) SendLogisticsBroadcast(command stahl4.Command) error {
	// Convert the Event to its JSON representation
	data, err := json.Marshal(command)
	// Check if no error occurred
	if err != nil {
		// save the error into the log
		// log.Print(err)
		return err
	} else {
		send := false
		// iterate over all servers
		for _, server := range *client.connections {
			// check if the server is logistic
			if server.Typ == 1 {
				// set server
				command.SetServer(&server)
				// send the message to the server
				err := sendMessage(command, client, data)
				// if there is an error
				if err != nil {
					// close the connection
					closeServerConnErr(client, command.GetServer(), server.Connection, err)
					return err
				}
				log.Print("Send LogisticBroadcast")
				// set the timestamp the thread safe way
				server.Mux.Lock()
				server.Timestamp = time.Now()
				server.Mux.Unlock()
				send = true
			}
		}
		if send {
			return nil
		} else {
			return errors.New("no message send")
		}
	}
}

func (client *Client) SendBroadcast(command stahl4.Command) error {
	// Convert the Event to its JSON representation
	data, err := json.Marshal(command)
	// Check if no error occurred
	if err != nil {
		// save the error into the log
		// log.Print(err)
		return err
	} else {
		send := false
		// iterate over all servers
		for _, server := range *client.connections {
			// set server
			command.SetServer(&server)
			// send the message to the server
			err := sendMessage(command, client, data)
			// if there is an error
			if err != nil {
				// close the connection
				closeServerConnErr(client, command.GetServer(), server.Connection, err)
				return err
			}
			log.Print("Send Broadcast")
			// set the timestamp the thread safe way
			server.Mux.Lock()
			server.Timestamp = time.Now()
			server.Mux.Unlock()
			send = true
		}
		if send {
			return nil
		} else {
			return errors.New("no message send")
		}
	}
}
func (client *Client) SendEventToIP(event stahl4.Event, ip string) error {
	remoteServer := connectServer(client, ip)
	if remoteServer == nil {
		return errors.New("couldn't connect to remote server")
	}
	// Convert the Event to its JSON representation
	data, err := json.Marshal(event)
	// Check if no error occurred
	if err != nil {
		// save the error into the log
		// log.Print(err)
		return err
	} else {

		// lock for validity of nil check
		remoteServer.Mux.Lock()
		if remoteServer.Connection == nil {
			remoteServer.Mux.Unlock()
			return errors.New("no server connection available")
		}
		// send the message to the server
		_, err := writeConnectionTimeout(client.writeTimeout, remoteServer.Connection, data)
		remoteServer.Mux.Unlock()

		// if there is an error
		if err != nil {
			// close the connection
			closeServerConnErr(client, remoteServer, remoteServer.Connection, err)
			return err
		}
		// set the timestamp the thread safe way
		remoteServer.Mux.Lock()
		remoteServer.Timestamp = time.Now()
		remoteServer.Mux.Unlock()
		return nil
	}
}

func handleServerConnection(connection net.Conn, client *Client) {
	// create buffer to read into
	buffer := make([]byte, 1024)
	// get a pointer to our remote server
	server := getServer(connection, client.connections)
	// if we were not able to find the server return to terminate the function
	if server == nil {
		return
	}
	// do as long as there is no new connection of this server
	for server.Connection == connection {
		timestamp := time.Now()
		// check if we need to send a heartbeat because we did not exchange messages for at least 20 seconds
		if (timestamp.Add(-20 * time.Second)).Sub(server.Timestamp) > 0*time.Second {
			// create our heartbeat
			var ping = make([]byte, 1)
			ping[0] = 0xF0
			// send the heartbeat
			_, err := writeConnectionTimeout(client.writeTimeout, connection, ping)
			// if there is an error during the heartbeat
			if err != nil {
				// close the connection
				closeServerConnErr(client, server, connection, err)
				// let this handle the termination of the loop
				continue
			}
			// set the timestamp the thread safe way
			server.Mux.Lock()
			server.Timestamp = time.Now()
			server.Mux.Unlock()
		}
		// read from the connection with the timeouts that were set
		n, _ := readConnectionTimeout(client.readTimeout, connection, buffer)
		// if we read more than 0 bytes
		if n > 0 {
			// set the timestamp the thread safe way
			server.Mux.Lock()
			server.Timestamp = time.Now()
			server.Mux.Unlock()
			// if we receive a keep alive
			if buffer[0] == 0xF0 {
				// don't do anything
				continue
			} else { // if we receive something else
				// try to parse the received message
				event := events.ParseEvent(server, buffer, n)
				if event != nil {
					client.AddToQueue(event)
				}
			}
		}
	}
}

// function to get a pointer to a remote server based on the connection
func getServer(conn net.Conn, connections *[]stahl4.RemoteServer) *stahl4.RemoteServer {
	// iterate over all known remote servers
	for _, server := range *connections {
		// if we find a match
		if conn == server.Connection {
			// return the match
			return &server
		}
	}
	// if we found nothing return nil
	return nil
}

// function to find and connect to new peers
func connectPeers(client *Client) {
	// get a list of valid ips
	hosts := hosts(client.Ip)
	//log.Print("connectPeers: got hosts, starting loop")
	for {
		// as long as we have less then 10 active connections
		log.Print("Currently connected to ", countConnected(client.connections), " clients.")
		fmt.Println("Currently connected to ", countConnected(client.connections), " clients.")
		if countConnected(client.connections) < 10 {
			// iterate over the ips
			for _, host := range hosts {
				// log.Printf("Trying to connect to %s", host)
				res := connectServer(client, host)
				// as long as we have less then 10 active connections
				//log.Print("Currently connected to ", countConnected(client.connections), " clients.")
				//fmt.Println("Currently connected to ", countConnected(client.connections), " clients.")
				if res != nil {
					log.Printf("Established connection to %s", host)
					if countConnected(client.connections) >= 10 {
						// we found enough peers
						break
					}
				}
			}
			// sleep to not generate to much traffic
			time.Sleep(1 * time.Second)
		} else { // don't search for new peers
			time.Sleep(60 * time.Second)
		}
	}
}

func checkAlreadyConnected(host string, connections *[]stahl4.RemoteServer, port string) bool {
	// iterate over all known remote servers
	for _, server := range *connections {
		// if we find a match
		// lock for validity of nil check
		server.Mux.Lock()
		if server.Connection != nil {
			// check if connection still exists
			if host+":"+port == server.Connection.RemoteAddr().String() {
				// return that we are already connected
				server.Mux.Unlock()
				return true
			}
		}
		server.Mux.Unlock()
	}
	// we are not connected to the server
	return false
}

func connectServer(client *Client, host string) *stahl4.RemoteServer {
	// check if we are already connected to this host
	if checkAlreadyConnected(host, client.connections, client.serverPort) {
		return nil
	}
	// create a buffer to hold an answer
	buffer := make([]byte, 1024)
	// try to connect to the ip and ignore the ip if there is an error
	connection, err := net.DialTimeout("tcp", host+":"+client.serverPort, client.readTimeout)
	if err != nil {
		// log.Print("Couldn't establish connection to server peer " + host + " because of " + err.Error())
		return nil
	}
	// handle handshake
	n, err := readConnectionTimeout(client.readTimeout, connection, buffer)
	// if there is an error while reading close this connection
	if err != nil {
		// log.Printf("Timeout while recv message from server peer: %s", connection.RemoteAddr().String())
		// close the connection
		_ = connection.Close()
		// let this handle the end of the loop
		return nil
	}
	var remoteType int
	i := 0
	for i < maxChallenges {
		i++
		var challenge ServerChallenge
		err = json.Unmarshal(buffer[:n], &challenge)
		if err != nil {
			log.Print("Couldnt not parse challenge. Closing connection to server peer " + connection.RemoteAddr().String() + " because of " + err.Error())
			_ = connection.Close()
			return nil
		}
		log.Print("Client: Received Server nonce ", challenge.Nonce)
		handshake, err := newClientHandshake(host, challenge.Nonce)
		if err != nil {
			log.Print("Couldnt not create the handshake. Closing connection to server peer " + connection.RemoteAddr().String() + " because of " + err.Error())
			_ = connection.Close()
			return nil
		}
		err = handshake.sign(client.key)
		if err != nil {
			log.Print("Couldnt sign handshake. Closing connection to server peer " + connection.RemoteAddr().String() + " because of " + err.Error())
			_ = connection.Close()
			return nil
		}
		hjson, err := json.Marshal(handshake)
		if err != nil {
			log.Print("Handshake marshal failed. Closing connection to server peer " + connection.RemoteAddr().String() + " because of " + err.Error())
			_ = connection.Close()
			return nil
		}
		log.Print("Client: Send Client Handshake ", string(hjson))
		_, err = writeConnectionTimeout(client.writeTimeout, connection, hjson)
		if err != nil {
			// log.Print("Closing connection to server peer " + connection.RemoteAddr().String() + " because of " + err.Error())
			_ = connection.Close()
			return nil
		}
		// try to get an answer and log any error
		n, err = readConnectionTimeout(client.readTimeout, connection, buffer)
		if err != nil {
			// log.Print("Closing connection to server peer " + connection.RemoteAddr().String() + " because of " + err.Error())
			_ = connection.Close()
			return nil
		}

		if strings.Contains(string(buffer[:n]), "type") {
			break
		}
	}
	// get the type of the server
	// try to parse the first answer
	var data map[string]string
	// decode the outer json layer and save the object into the variable
	err = json.Unmarshal(json.RawMessage(buffer[:n]), &data)
	// if there was an error unmarshalling the message
	if err != nil {
		log.Print("Parsing of handshake response failed. Closing connection to server peer " + connection.RemoteAddr().String() + " because of " + err.Error())
		_ = connection.Close()
		return nil
	}
	remoteType, err = strconv.Atoi(data["type"])
	if err != nil {
		// log.Print("type error of " + connection.RemoteAddr().String() + " because of " + err.Error())
		log.Print(err.Error())
		return nil
	}
	// create the server object and append it to our known servers
	server := stahl4.RemoteServer{Ip: connection.RemoteAddr(), Connection: connection, Timestamp: time.Now(), Typ: uint8(remoteType)}
	// lock the client for thread safety
	client.mux.Lock()
	*client.connections = append(*client.connections, server)
	// unlock the client
	client.mux.Unlock()
	log.Printf("client: Finished handshake with %s", client.Ip)
	// start handling the Connection (e.g. keep alive)
	go handleServerConnection(connection, client)
	return &server
}

// function to count the servers with active connections
func countConnected(remoteServers *[]stahl4.RemoteServer) int {
	// create the count variable
	var count int
	// iterate over all servers we know of
	for _, server := range *remoteServers {
		// if there is a connection
		// lock for validity of nil check
		server.Mux.Lock()
		if server.Connection != nil {
			// increase the count
			count++
		}
		server.Mux.Unlock()
	}
	// return the count
	return count
}

// function to close a connection if an error occurred
func closeServerConnErr(client *Client, server *stahl4.RemoteServer, connection net.Conn, err error) {
	// write the event into the log
	// log.Print("Closing connection to a server peer because of " + err.Error())
	// lock client for thread-safety
	client.mux.Lock()
	// delete the connection from the client
	var newConns []stahl4.RemoteServer
	for _, c := range *client.connections {
		// append only if we are not the deleted connection
		if c.Connection != connection {
			newConns = append(newConns, c)
		}
	}
	// update client connections
	*client.connections = newConns
	// close the connection if not already closed
	if connection != nil {
		_ = connection.Close()
	}
	client.mux.Unlock()
	// lock the server for thread safety
	server.Mux.Lock()
	// set the connection for the server to nil to indicate it (important so that e.g. handleServerConnection will
	// terminate)
	if server.Connection != nil {
		server.Connection = nil
	}
	// unlock the server
	server.Mux.Unlock()
}

// function which creates a list of ips which could be peers
func hosts(ip string) []string {
	// create an array to save the ips
	var ips []string
	// get our own ip and save each single octet
	self := strings.Split(ip, ".")
	// if we do not have a valid ip
	if len(self) < 3 {
		// write a fatal log entry (i.e. write the log entry and exit with code 1)
		log.Fatal("Couldn't resolve hosts IP address, only found: " + strings.Join(self, ""))
	}
	// third octet
	thirdOctet, err := strconv.Atoi(self[2])
	// if the conversion fails log the error
	if err != nil {
		log.Print(err.Error())
	}
	// iterate over all possible values
	for i := 0; i < 256; i++ {
		// use the iterator as offset of our ip to build a proper network
		// i.e. such that not all services connect to ips *.*.*.2 - *.*.*.11
		j := uint8(thirdOctet + i)
		// create the ips and append them to the list excluding our own and reserved addresses
		if (j >= 1) && (j <= 254) && int(j) != thirdOctet {
			ips = append(ips, strings.Join(self[0:2], ".")+"."+strconv.Itoa(int(j))+"."+self[3])
		}
	}
	//return the list of ips
	return ips
}

// function which tries to return our ip
func hostIP() string {
	// create the variable to save the ip
	var ip string
	// the subnets we want to look for if we find multiple ips
	subnet1 := "10.32"
	subnet2 := "10.33"
	bad_subnet := "127."
	// get our hostname and log the error if there is one
	name, err := os.Hostname()
	if err != nil {
		log.Print(err.Error())
	}
	// try to resolve the hostname and log the error if there is one
	addresses, err := net.LookupHost(name)
	if err != nil {
		log.Print(err.Error())
	}
	// iterate over all ips we found
	for _, address := range addresses {
		// look for the (first) one which contains our predefined part
		if strings.Contains(address, subnet1) || strings.Contains(address, subnet2) {
			ip = address
			break
		}
	}
	// if we did not find an ip yet
	if len(ip) == 0 {
		// fall back to iterate over our interfaces and log all errors
		addresses, err := net.InterfaceAddrs()
		if err != nil {
			log.Print(err.Error())
		}
		// iterate over the ips we found on the interfaces
		for _, address := range addresses {
			// ignore ipv6 ips
			if !strings.Contains(address.String(), ":") {
				// look for the (first) one which contains our predefined part
				if strings.Contains(address.String(), subnet1) || strings.Contains(address.String(), subnet2) {
					ip = address.String()
					break
				}
			}
		}
	}
	// use fallback only if not forbidden
	if !(len(os.Args) == 2 && strings.Compare(os.Args[1], "game-network-only") == 0) {
		// if we are not running in one of the two game subnets -> just use the first ip we find that is not
		// localhost
		if len(ip) == 0 {
			// fall back to iterate over our interfaces and log all errors
			addresses, err := net.InterfaceAddrs()
			if err != nil {
				log.Print(err.Error())
			}
			// iterate over the ips we found on the interfaces
			for _, address := range addresses {
				// ignore ipv6 ips
				if !strings.Contains(address.String(), ":") {
					// look for the (first) one which contains our predefined part
					if !strings.Contains(address.String(), bad_subnet) {
						ip = address.String()
						break
					}
				}
			}
		}
	}
	if len(ip) == 0 {
		panic("Could not find valid host IP.")
	}
	// remove the netmask if there is one
	ip = strings.Split(ip, "/")[0]
	// return the ip
	fmt.Printf("Using IP %s \n", ip)
	//return "127.0.0.1"
	return ip
}

func getRemoteServer(index int, client Client) *stahl4.RemoteServer {
	return &(*client.connections)[index]
}

func sendMessage(command stahl4.Command, client *Client, data []byte) error {
	// lock for validity of nil check
	server := command.GetServer()
	server.Mux.Lock()
	if server.Connection == nil {
		server.Mux.Unlock()
		log.Print("No Server connection available!")
		return errors.New("no server connection available")
	}
	_, err := writeConnectionTimeout(client.writeTimeout, server.Connection, data)
	server.Mux.Unlock()
	return err
}
