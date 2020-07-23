package stahl4

import (
	"crypto/rsa"
	"net"
	"sync"
	"time"
)

type SteelServer interface {
	AddToQueue(Command)
	ReceiveCommand() Command
	SendEvent(Event) error
	Listen()
}

type SteelClient interface {
	AddToQueue(Event)
	ReceiveEvent() Event
	SendCommand(Command) error
	SendLogisticCommand(Command) error
	SendLogisticsBroadcast(Command) error
	SendBroadcast(Command) error
	SendEventToIP(Event, string) error
	Start()
}

type RemoteClient struct {
	// mutex for thread safety
	Mux sync.Mutex
	// the id that we give that client (hash of the PubKey)
	Id string
	// the ip of the client
	Ip net.Addr
	// the connection object of the client
	Connection net.Conn
	// the pub key of the client
	PubKey rsa.PublicKey
}

// the type that we use to push commands onto the queue (to be able to expand the data if needed)
type CommandItem struct {
	// the actual command that has to be handled
	Command Command
}

// the type that we use to push events onto the queue (to be able to expand the data if needed)
type EventItem struct {
	// the actual event that has to be handled
	Event Event
}

type RemoteServer struct {
	// mutex for thread safety
	Mux sync.Mutex
	// the ip of the server
	Ip net.Addr
	// the connection object of the server
	Connection net.Conn
	// the timestamp of the last send message
	Timestamp time.Time
	// the type of the server
	Typ uint8
}
