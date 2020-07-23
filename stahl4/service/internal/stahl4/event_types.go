package stahl4

type Event interface {
	UpdateState(*State)
	ToByteArray() []byte
	MarshalJSON() ([]byte, error)
	GetClient() *RemoteClient
	SetClient(*RemoteClient)
	GetServer() *RemoteServer
	SetServer(*RemoteServer)
}

type BroadcastEvent interface {
	GetID() string
	SetID(string)
}
