package main

//export GOPATH=$HOME/Documents/MagAnd/Heismappe/

import (
	"asd"
	"driver"
	//"fmt"
	//"time"
	def "config"
	//"network"
	//"queue"
)

var OnlineElevs = make(map[string]network.UdpConnection)

var outgoingMsg = make(chan def.Message, 10)

func main() {
	driver.Elev_init()
	def.CurFloor = driver.Get_floor_sensor_signal()
	def.CurDir = 0
	go asd.Update_lights(outgoingMsg)

	//outgoingMsg := make(chan def.Message)
	//incomingMsg := make(chan def.Message)
	//go network.Init(outgoingMsg, incomingMsg)

	//def.Orders = []int{-3, -4, 2, 3, -6}
	//def.CurFloor = 3
	//def.CurDir = -1
	go asd.Go_to_floor()
	//fmt.Printf("%v", queue.Update_orderlist(orderlist, 4))
	//fmt.Print(queue.Cost(orderlist, 4))

	quit := make(chan int)
	go asd.Quit_program(quit)

	<-quit
}


func Handle_msg(msg def.Message){
	const aliveTimeout = 2 * time.Second

	switch msg.Category {
	case def.Alive:
		if connection, exist := def.onlineElevs[msg.Addr]; exist {
			connection.Timer.Reset(aliveTimeout)
		} else {
			newConnection := network.UdpConnection{msg.Addr, time.NewTimer(aliveTimeout)}
			def.onlineElevs[msg.Addr] = newConnection
			numOnline = len(def.onlineElevs)
			go connectionTimer(&newConnection)
			log.Printf("%sConnection to IP %s established!%s", def.ColG, msg.Addr[0:15], def.ColN)
		}
	case def.NewOrder:

	case def.CompleteOrder:
		
	case def.Cost:
		
	}
}
