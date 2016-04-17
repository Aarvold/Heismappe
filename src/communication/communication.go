package communication

import (
 	def "config"
 	"driver"
 	"queue"
 	"network"
 	"time"
 	"handleOrders"
 	"fmt"
)

var onlineElevs = make(map[string]network.UdpConnection)

func Handle_msg(incomingMsg, outgoingMsg, costMsg, orderCompleted chan def.Message) {
	
	for {

		msg := <-incomingMsg
		const aliveTimeout = 6 * time.Second

		switch msg.Category {
		case def.Alive:
			//if connection exists
			if connection, exist := onlineElevs[msg.Addr]; exist{
				connection.Timer.Reset(aliveTimeout)
			} else {
				add_new_connection(msg.Addr,aliveTimeout)
			}

		case def.NewOrder:
			//fmt.Printf("%sNew external order recieved to floor %d %s \n", def.ColM, helpFunc.Order_dir(msg.Floor, msg.Button), def.ColN)
			driver.Set_button_lamp(msg.Button, msg.Floor, 1)
			costMsg := def.Message{Category: def.Cost, Floor: msg.Floor, Button: msg.Button, Cost: handleOrders.Cost(queue.Get_Orders(), queue.Order_direction(msg.Floor, msg.Button))}
			outgoingMsg <- costMsg

		case def.CompleteOrder:
			driver.Set_button_lamp(msg.Button, msg.Floor, 0)
			orderCompleted <- msg

		case def.Cost:
			//see handleOrders.assign_external_order
			costMsg <- msg
		}
	}
}

func add_new_connection(addr string, aliveTimeout time.Duration){
	newConnection := network.UdpConnection{addr, time.NewTimer(aliveTimeout)}
	onlineElevs[addr] = newConnection
	handleOrders.NumOfOnlineElevs = len(onlineElevs)
	go connection_timer(&newConnection)
	fmt.Printf("%sConnection to IP %s established!%s\n", def.ColG, addr[0:15], def.ColN)
}

func connection_timer(connection *network.UdpConnection) {
	<-connection.Timer.C
	handleOrders.NumOfOnlineElevs = handleOrders.NumOfOnlineElevs -1
	delete(onlineElevs, connection.Addr)
	fmt.Printf("%sIP %v disconnected %s\n", def.ColR, connection.Addr[0:15], def.ColN)
}