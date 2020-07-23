package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"math/rand"
	"stahl4/service/internal/commands"
	"stahl4/service/internal/events"
	"stahl4/service/internal/network"
	"stahl4/service/internal/stahl4"
	"strconv"
	"time"
)

func startMachine(typ uint8) {
	rand.Seed(time.Now().UTC().UnixNano())
	// println("Parse Public Key")
	idHash := sha256.New()
	idHash.Write(x509.MarshalPKCS1PublicKey(&privK.PublicKey))
	id := hex.EncodeToString(idHash.Sum(nil))
	// println("Creating Server")
	server := network.CreateServer(typ, port)
	// println("Creating Client")
	client := network.CreateClient(port, privK)
	// println("Created Server and Client.")
	broadcasts := make([]stahl4.BroadcastCommand, 0)
	state := &stahl4.State{}
	switch typ {
	case logistic:
		logistic := stahl4.NewLogistic()
		state.Logistic = logistic
	case production:
		production := stahl4.NewProduction()
		state.Production = production
		go func(state *stahl4.State) {
			if state.Production.MaterialAmount > 0 && state.Production.ProductAmount < state.Production.ProductCapacity {
				state.Mux.Lock()
				state.Production.ProductAmount++
				state.Production.MaterialAmount--
				state.Mux.Unlock()
			}
			time.Sleep(1 * time.Second)
		}(state)
	}
	go server.Listen()
	// println("Started Server.")
	go client.Start()
	// println("Started Client.")

	sendMaterialBroadcastTime := time.Now()
	sendProductBroadcastTime := time.Now()

	for {
		command := server.ReceiveCommand()
		if command != nil {
			bc, ok := command.(stahl4.BroadcastCommand)
			if ok {
				known := false
				for _, knownBC := range broadcasts {
					if stahl4.CompareBroadcasts(bc, knownBC) {
						known = true
						break
					}
				}
				if !known {
					broadcasts = appendBroadcast(broadcasts, bc)
					_ = client.SendBroadcast(command)
					if typ == logistic {
						event := command.ToEvent(state)
						if event != nil {
							_ = client.SendEventToIP(event, bc.GetIP())
						}
					}
				}
			} else {
				event := command.ToEvent(state)
				if event != nil {
					_ = server.SendEvent(event)
				}
			}
		}
		event := client.ReceiveEvent()
		if event != nil {
			event.UpdateState(state)
			switch event.(type) {
			case *events.GetMaterialsBroadcastAcceptEvent:
				if typ == production {
					state.Mux.Lock()
					materialAmount := state.Production.MaterialAmount
					materialCapacity := state.Production.MaterialCapacity
					askedForMaterials := state.Production.AskedForMaterials
					state.Mux.Unlock()

					if askedForMaterials == false {
						amount := materialCapacity - materialAmount
						if amount >= 150 {
							amount = 150
						}
						accept := commands.NewGetMaterialsCommand(event.GetClient(), event.GetServer(), uint8(amount))
						if accept != nil {
							_ = client.SendLogisticCommand(accept)
							state.Mux.Lock()
							state.Production.AskedForMaterials = true
							state.Mux.Unlock()
						}
					}
				}
			case *events.GetProductsBroadcastAcceptEvent:
				if typ == production {
					state.Mux.Lock()
					productCapacity := state.Production.ProductCapacity
					productAmount := state.Production.ProductAmount
					askedToGetProducts := state.Production.AskedToGetProducts
					state.Mux.Unlock()

					if askedToGetProducts == false {
						amount := productCapacity - productAmount
						if amount >= 150 {
							amount = 150
						}
						accept := commands.NewGetProductsCommand(event.GetClient(), event.GetServer(), uint8(amount))
						_ = client.SendLogisticCommand(accept)
						state.Mux.Lock()
						state.Production.AskedToGetProducts = true
						state.Mux.Unlock()
					}
				}
			}
		}
		// println("Starting nice to have comm logic")
		switch typ {
		case logistic:
			//println("We are logistic, matSt:", state.Logistic.MaterialsStored, "prodSt:", state.Logistic.ProductsStored)
			if float32(state.Logistic.MaterialsStored) < 0.15*float32(state.Logistic.TotalMaterialCapacity) ||
				state.Logistic.NumberOfNoMaterialAnswers > 10 {
				// refill materials from the central storage
				println("Refilling materials")
				state.Mux.Lock()
				state.Logistic.Busy = true
				state.Mux.Unlock()

				// sleep to simulate material refilling
				stahl4.MovementSleep()

				state.Mux.Lock()
				state.Logistic.MaterialsStored = state.Logistic.TotalMaterialCapacity
				state.Logistic.NumberOfNoMaterialAnswers = 0
				state.Logistic.Busy = false
				state.Mux.Unlock()
			}
			if float32(state.Logistic.ProductsStored) > 0.32*float32(state.Logistic.TotalProductCapacity) {
				// drop of products in central storage
				state.Logistic.Busy = true
				stahl4.MovementSleep()
				state.Mux.Lock()
				state.Logistic.ProductsStored = 0
				state.Logistic.Busy = false
				state.Mux.Unlock()
			}
		case production:
			state.Mux.Lock()
			materialAmount := float32(state.Production.MaterialAmount)
			materialCapacity := float32(state.Production.MaterialCapacity)
			askedForMaterials := state.Production.AskedForMaterials
			state.Mux.Unlock()
			//println("We are prod, matAm:", materialAmount, "prodAsk:", askedForMaterials)
			if materialAmount <= 0.05*materialCapacity &&
				!askedForMaterials && time.Now().Add(-30*time.Second).Sub(sendMaterialBroadcastTime) >= 0*time.Second {
				hash := sha512.New().Sum([]byte(client.Ip + time.Now().String() + strconv.Itoa(rand.Int())))
				broadcast := commands.NewGetMaterialsBroadcastCommand(nil, nil, string(hash), client.Ip)
				err := client.SendBroadcast(broadcast.(stahl4.Command))
				if err == nil {
					broadcasts = appendBroadcast(broadcasts, broadcast)
					sendMaterialBroadcastTime = time.Now()
					time.Sleep(5 * time.Second)
				}
			} else if materialAmount <= 0.25*materialCapacity &&
				!askedForMaterials {
				err := client.SendLogisticsBroadcast(
					commands.NewGetMaterialsCommand(nil, nil, uint8(0.75*materialCapacity/10)))
				if err == nil {
					state.Mux.Lock()
					state.Production.AskedForMaterials = true
					state.Mux.Unlock()
					time.Sleep(5 * time.Second)
				}
			} else if materialAmount <= 0.5*materialCapacity &&
				!askedForMaterials {
				err := client.SendLogisticCommand(commands.NewGetMaterialsCommand(nil, nil, 150))
				if err == nil {
					state.Mux.Lock()
					state.Production.AskedForMaterials = true
					state.Mux.Unlock()
					time.Sleep(5 * time.Second)
				}
			} else {
				if !askedForMaterials {
					err := client.SendLogisticCommand(commands.NewGetMaterialsCommand(nil, nil, uint8(state.Production.MaterialCapacity-state.Production.MaterialAmount)))
					if err == nil {
						state.Mux.Lock()
						state.Production.AskedForMaterials = true
						state.Mux.Unlock()
						time.Sleep(5 * time.Second)
					}
				}
			}

			state.Mux.Lock()
			productAmount := float32(state.Production.ProductAmount)
			productCapacity := float32(state.Production.ProductCapacity)
			askedToGetProducts := state.Production.AskedToGetProducts
			state.Mux.Unlock()
			// println("We are prod, prodAm:", productAmount, "prodAsk:", askedToGetProducts)
			if productAmount >= 0.95*productCapacity &&
				!askedToGetProducts && time.Now().Add(-30*time.Second).Sub(sendProductBroadcastTime) >= 0*time.Second {
				hash := sha512.New().Sum([]byte(client.Ip + time.Now().String() + strconv.Itoa(rand.Int())))
				broadcast := commands.NewGetProductsBroadcastCommand(nil, nil, string(hash), client.Ip)
				err := client.SendBroadcast(broadcast.(stahl4.Command))
				if err == nil {
					sendProductBroadcastTime = time.Now()
					broadcasts = appendBroadcast(broadcasts, broadcast)
					time.Sleep(5 * time.Second)
				}
			} else if productAmount >= 0.75*productCapacity &&
				!askedToGetProducts {
				err := client.SendLogisticsBroadcast(
					commands.NewGetProductsCommand(nil, nil, uint8(0.75*productCapacity/10)))
				if err == nil {
					state.Mux.Lock()
					state.Production.AskedToGetProducts = true
					state.Mux.Unlock()
					time.Sleep(5 * time.Second)
				}
			} else if productAmount >= 0.5*productCapacity &&
				!askedToGetProducts {
				err := client.SendLogisticCommand(commands.NewGetProductsCommand(nil, nil, 150))
				if err == nil {
					state.Mux.Lock()
					state.Production.AskedToGetProducts = true
					state.Mux.Unlock()
					time.Sleep(5 * time.Second)
				}
			} else {
				if !askedToGetProducts {
					err := client.SendLogisticCommand(commands.NewGetProductsCommand(nil, nil, uint8(state.Production.ProductCapacity-state.Production.ProductAmount)))
					if err == nil {
						state.Mux.Lock()
						state.Production.AskedToGetProducts = true
						state.Mux.Unlock()
						time.Sleep(5 * time.Second)
					}
				}
			}
		}

		rInt := rand.Int() % 1000000 // seems like a good production value
		if rInt < 7 {
			_ = client.SendGetCommand(commands.GetDataCommand{
				AskedId: id,
				Client:  nil,
				Server:  nil,
			})
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func appendBroadcast(broadcasts []stahl4.BroadcastCommand, broadcast stahl4.BroadcastCommand) []stahl4.BroadcastCommand {
	if len(broadcasts) > 512 {
		copy(broadcasts, broadcasts[25:])
	}
	return append(broadcasts, broadcast)
}
