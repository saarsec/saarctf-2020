package commands

/*

import (
	"encoding/json"
	"stahl4/service/internal/stahl4"
	"stahl4/service/internal/events"
	"log"
)
type SpeedUpCommand struct {
	Message string               `json:"message"`
	Client  *stahl4.RemoteClient `json:"-"` //Ignore in Marshal
	Server  *stahl4.RemoteServer `json:"-"` //Ignore in Marshal
}

func (speedUpCommand *SpeedUpCommand) ToEvent(state *stahl4.State) stahl4.Event {
	// NOT YET TESTED
	var answer stahl4.Event
	if state.Production.Halted {
		msg := "Production is currently halted. Ignoring your message."

		//option 2 for "halted"
		speedUpCommand.UpdateState(state, msg, 2)

		answer = events.NewHaltedEvent(msg, speedUpCommand.Client)
	} else if state.Production.SpedUp {
		msg := "Already at highest speed."

		//option 1 for "highest speed"
		speedUpCommand.UpdateState(state, msg, 1)

		answer = events.NewHighestSpeedEvent(msg, speedUpCommand.Client)
	} else {
		msg := "Speeding up production."

		//option 0 for "Speeding up"
		speedUpCommand.UpdateState(state, msg, 0)

		answer = events.NewSpeedingUpEvent(msg, speedUpCommand.Client)
	}
	return answer
}

/*
	if state.Production.SlowedDown {
		//we are slowed down, so we stop it
		state.Production.SlowedDown = false
		//logging
		log.Println("Sped up to normal speed, as requested from " + speedingUpEvent.Client.Ip.String())
	} else {
		//we are not slowed down, so we just speed down
		state.Production.SpedUp = true
		//logging
		log.Println("Sped up to highest speed, as requested from " + speedingUpEvent.Client.Ip.String())
	}

func (speedUpCommand *SpeedUpCommand) UpdateState(state *stahl4.State, message string, option int) {
	//NOT TESTED
	if option == 0 {
		//increase Speed
		if state.Production.SlowedDown {
			//we are slowed down, so we stop it
			state.Production.SlowedDown = false
			//logging
			log.Println("Sped up to normal speed, as requested from " + speedUpCommand.Client.Ip.String())
		} else {
			//we are not slowed down, so we just speed down
			state.Production.SpedUp = true
			//logging
			log.Println("Sped up to highest speed, as requested from " + speedUpCommand.Client.Ip.String())
		}
	} else if option == 1 {
		//we are at highest speed ; no state changes
		log.Println("Already at highest speed. Ignoring request of " + speedUpCommand.Client.Ip.String())

	} else {
		//we are halted
		log.Println("Ignoring speed changing command, because server is halted. Request came from " +
			speedUpCommand.Client.Ip.String())
	}
	state.Production.LastMessage = message
}

func (speedUpCommand *SpeedUpCommand) ToByteArray() []byte {
	return []byte(speedUpCommand.Message)
}

func (speedUpCommand *SpeedUpCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":  	"SpeedUpCommand",
		"message":	speedUpCommand.Message,
	})
}

func (speedUpCommand *SpeedUpCommand) GetClient() *stahl4.RemoteClient {
	return speedUpCommand.Client
}

func (speedUpCommand *SpeedUpCommand) SetClient(client *stahl4.RemoteClient) {
	speedUpCommand.Client = client
}

func (speedUpCommand *SpeedUpCommand) GetServer() *stahl4.RemoteServer {
	return speedUpCommand.Server
}

func (speedUpCommand *SpeedUpCommand) SetServer(server *stahl4.RemoteServer) {
	speedUpCommand.Server = server
}

type SlowDownCommand struct {
	Message string               `json:"message"`
	Client  *stahl4.RemoteClient `json:"-"` //Ignore in Marshal
	Server  *stahl4.RemoteServer `json:"-"` //Ignore in Marshal
}

func (slowDownCommand *SlowDownCommand) ToEvent(state *stahl4.State) stahl4.Event {
	//NOT TESTED
	var answer stahl4.Event
	if state.Production.Halted {
		msg := "Production is currently halted ignoring your message."

		//option 2 for "halted"
		slowDownCommand.UpdateState(state, msg, 2)

		answer = events.NewHaltedEvent(msg, slowDownCommand.Client)
	} else if state.Production.SlowedDown {
		msg := "Already at lowest speed."
		slowDownCommand.UpdateState(state, msg, 1)
		answer = events.NewLowestSpeedEvent(msg, slowDownCommand.Client)
	} else {
		msg := "Slowing down production."
		slowDownCommand.UpdateState(state, msg, 0)
		answer = events.NewSlowingDownEvent(msg, slowDownCommand.Client)
	}
	return answer
}

func (slowDownCommand *SlowDownCommand) UpdateState(state *stahl4.State, message string, option int) {
	// NOT TESTED
	if option == 0 {
		//decrease Speed
		if state.Production.SpedUp {
			//we are sped up, so we stop it
			state.Production.SpedUp = false
			//logging
			log.Println("Slowed down to normal speed, as requested from " + slowDownCommand.Client.Ip.String())
		} else {
			//we are not sped up, so we just slow down
			state.Production.SlowedDown = true
			//logging
			log.Println("Slowed down to slowest speed, as requested from " + slowDownCommand.Client.Ip.String())
		}
	} else if option == 1 {
		//we are at lowest speed ; no state changes
		log.Println("Already at lowest speed. Ignoring request of " + slowDownCommand.Client.Ip.String())
	} else {
		//we are halted
		log.Println("Ignoring speed changing command, because server is halted. Request came from " +
			slowDownCommand.Client.Ip.String())
	}
	state.Production.LastMessage = message
}
func (slowDownCommand *SlowDownCommand) ToByteArray() []byte {
	return []byte(slowDownCommand.Message)
}

func (slowDownCommand *SlowDownCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type":  	"SlowDownCommand",
		"message":	slowDownCommand.Message,
	})
}

func (slowDownCommand *SlowDownCommand) GetClient() *stahl4.RemoteClient {
	return slowDownCommand.Client
}

func (slowDownCommand *SlowDownCommand) SetClient(client *stahl4.RemoteClient) {
	slowDownCommand.Client = client
}

func (slowDownCommand *SlowDownCommand) GetServer() *stahl4.RemoteServer {
	return slowDownCommand.Server
}

func (slowDownCommand *SlowDownCommand) SetServer(server *stahl4.RemoteServer) {
	slowDownCommand.Server = server
}
*/
