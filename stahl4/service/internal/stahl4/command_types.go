package stahl4

type Command interface {
	ToEvent(*State) Event
	ToByteArray() []byte
	MarshalJSON() ([]byte, error)
	GetClient() *RemoteClient
	SetClient(*RemoteClient)
	GetServer() *RemoteServer
	SetServer(*RemoteServer)
}

type BroadcastCommand interface {
	Command
	GetID() string
	SetID(string)
	GetIP() string
	SetIP(string)
}
