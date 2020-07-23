package commands

import (
	"stahl4/service/internal/stahl4"
)

func NewGetDataCommand(client *stahl4.RemoteClient, server *stahl4.RemoteServer, askedid string) stahl4.Command {
	return &GetDataCommand{
		AskedId: askedid,
		Client:  client,
		Server:  server,
	}
}

func NewGetErrorCommand(client *stahl4.RemoteClient, server *stahl4.RemoteServer, askedid string) stahl4.Command {
	return &GetErrorCommand{
		AskedId: askedid,
		Client:  client,
		Server:  server,
	}
}

func NewGetMaterialsCommand(client *stahl4.RemoteClient, server *stahl4.RemoteServer, amount uint8) stahl4.Command {
	return &GetMaterialsCommand{
		Amount:  amount,
		Message: "Getting new materials.",
		Client:  client,
		Server:  server,
	}
}

func NewGetProductsCommand(client *stahl4.RemoteClient, server *stahl4.RemoteServer, amount uint8) stahl4.Command {
	return &GetProductsCommand{
		Amount:  amount,
		Message: "Fetching new products.",
		Client:  client,
		Server:  server,
	}
}

func NewGetMaterialsBroadcastCommand(
	client *stahl4.RemoteClient, server *stahl4.RemoteServer, id string, ip string) stahl4.BroadcastCommand {
	return &GetMaterialsBroadcastCommand{
		//Amount:  amount,
		Message: "Asking for new materials.",
		ID:      id,
		IP:      ip,
		Client:  client,
		Server:  server,
	}
}

func NewGetProductsBroadcastCommand(
	client *stahl4.RemoteClient, server *stahl4.RemoteServer, id string, ip string) stahl4.BroadcastCommand {
	return &GetProductsBroadcastCommand{
		//Amount:  amount,
		Message: "Asking to fetch new products.",
		ID:      id,
		IP:      ip,
		Client:  client,
		Server:  server,
	}
}

/*
func NewSpeedUpCommand(client *stahl4.RemoteClient, server *stahl4.RemoteServer, message string) stahl4.Command {
	return &SpeedUpCommand{
		Message: message,
		Client:  client,
		Server:  server,
	}
}

func NewSlowDownCommand(client *stahl4.RemoteClient, server *stahl4.RemoteServer, message string) stahl4.Command {
	return &SlowDownCommand{
		Client:  client,
		Message: message,
		Server:  server,
	}
}
*/
