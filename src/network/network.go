package network

import (
	def "config"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"handleOrders"
)

func Init(outgoingMsg, incomingMsg chan def.Message) {
	// Ports randomly chosen to reduce likelihood of port collision.
	const localListenPort = 37203
	const broadcastListenPort = 37204
	const messageSize = 1024

	var udpSend = make(chan udpMessage)
	var udpReceive = make(chan udpMessage, 10)
	err := udpInit(localListenPort, broadcastListenPort, messageSize, udpSend, udpReceive)
	if err != nil {
		log.Print("UdpInit() error: %v \n", err)
	}

	handleOrders.ImConnected = true

	go aliveSpammer(outgoingMsg)
	go forwardOutgoing(outgoingMsg, udpSend)
	go forwardIncoming(incomingMsg, udpReceive)

	fmt.Println(def.ColG, "Network successfully initialised", def.ColN)
}

// Periodically notifyes other elevators that this is elevator is alive
func aliveSpammer(outgoingMsg chan<- def.Message) {
	const spamInterval = 2000 * time.Millisecond
	alive := def.Message{Category: def.Alive, Floor: -1, Button: -1, Cost: -1}
	for {
		outgoingMsg <- alive
		time.Sleep(spamInterval)
	}
}

// ForwardOutgoing continuosly checks for messages to be sent by reading the OutgoingMsg channel.
// Sends messages as JSON 
func forwardOutgoing(outgoingMsg <-chan def.Message, udpSend chan<- udpMessage) {
	for {
		msg := <-outgoingMsg

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("%sjson.Marshal error: %v\n%s", def.ColR, err, def.ColN)
		}

		udpSend <- udpMessage{raddr: "broadcast", data: jsonMsg, length: len(jsonMsg)}
	}
}


func forwardIncoming(incomingMsg chan<- def.Message, udpReceive <-chan udpMessage) {
	for {
		udpMessage := <-udpReceive
		var message def.Message

		if err := json.Unmarshal(udpMessage.data[:udpMessage.length], &message); err != nil {
			fmt.Printf("json.Unmarshal error: %s\n", err)
		}

		message.Addr = udpMessage.raddr
		incomingMsg <- message
	}
}
