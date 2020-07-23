package events

import (
	"stahl4/service/internal/stahl4"
)

func NewGetEvent(message string, client *stahl4.RemoteClient, server *stahl4.RemoteServer) stahl4.Event {
	return &GetEvent{
		Message: message,
		Client:  client,
		Server:  server,
	}
}

func NewBusyEvent(msg string, option uint8, client *stahl4.RemoteClient, server *stahl4.RemoteServer) stahl4.Event {
	return &BusyEvent{
		//Message: "The server is busy and can not process your request.",
		Message: msg,
		Client:  client,
		Option:  option,
		Server:  server,
	}
}

func NewNoProductSpaceEvent(msg string, client *stahl4.RemoteClient, server *stahl4.RemoteServer) stahl4.Event {
	return &NoProductSpaceEvent{
		//Message: "There is not enough empty product space.",
		Message: msg,
		Client:  client,
		Server:  server,
	}
}

func NewNoMaterialsEvent(msg string, client *stahl4.RemoteClient, server *stahl4.RemoteServer) stahl4.Event {
	return &NoMaterialsEvent{
		//Message: "There are not enough materials to fulfil your request.",
		Message: msg,
		Client:  client,
		Server:  server,
	}
}

func NewGetMaterialsBroadcastAcceptEvent(message string, id string, client *stahl4.RemoteClient, server *stahl4.RemoteServer) stahl4.Event {
	return &GetMaterialsBroadcastAcceptEvent{
		Message: message,
		//Amount:  amount,
		ID:     id,
		Client: client,
		Server: server,
	}
}

func NewGetProductsBroadcastAcceptEvent(message string, id string, client *stahl4.RemoteClient, server *stahl4.RemoteServer) stahl4.Event {
	return &GetProductsBroadcastAcceptEvent{
		Message: message,
		//Amount:  amount,
		ID:     id,
		Client: client,
		Server: server,
	}
}

func NewDeliveredMaterialsEvent(message string, client *stahl4.RemoteClient, server *stahl4.RemoteServer, amount uint8) stahl4.Event {
	return &DeliveredMaterialsEvent{
		Message: message,
		Amount:  amount,
		Client:  client,
		Server:  server,
	}
}

func NewFetchedProductsEvent(message string, client *stahl4.RemoteClient, server *stahl4.RemoteServer, amount uint8) stahl4.Event {
	return &FetchedProductsEvent{
		Message: message,
		Amount:  amount,
		Client:  client,
		Server:  server,
	}
}

/*
func NewLowestSpeedEvent(message string, client *stahl4.RemoteClient, server *stahl4.RemoteServer) stahl4.Event {
	return &LowestSpeedEvent{
		Message: message,
		Client:  client,
		Server:  server,
	}
}

func NewHighestSpeedEvent(message string, client *stahl4.RemoteClient, server *stahl4.RemoteServer) stahl4.Event {
	return &HighestSpeedEvent{
		Message: message,
		Client:  client,
		Server:  server,
	}
}

func NewSlowingDownEvent(message string, client *stahl4.RemoteClient, server *stahl4.RemoteServer) stahl4.Event {
	return &SlowingDownEvent{
		Message: message,
		Client:  client,
		Server:  server,
	}
}

func NewSpeedingUpEvent(message string, client *stahl4.RemoteClient, server *stahl4.RemoteServer) stahl4.Event {
	return &SpeedingUpEvent{
		Message: message,
		Client:  client,
		Server:  server,
	}
}

func NewHaltedEvent(message string, client *stahl4.RemoteClient, server *stahl4.RemoteServer) stahl4.Event {
	return &HaltedEvent{
		Message: message,
		Client:  client,
		Server:  server,
	}
}
*/
