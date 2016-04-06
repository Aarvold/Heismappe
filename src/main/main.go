package main

//export GOPATH=$HOME/Documents/MagAnd/Heismappe/

import (
	"asd"
	"driver"
	"fmt"
	//"time"
	def "config"
	//"network"
	"queue"
)

func main() {
	driver.Elev_init()
	def.CurFloor = driver.Get_floor_sensor_signal()
	def.CurDir = 0
	go asd.Update_lights()

	//outgoingMsg := make(chan def.Message)
	//incomingMsg := make(chan def.Message)
	//go network.Init(outgoingMsg, incomingMsg)

	var orderlist = []int{-5, -4, 2, 3, -6}
	def.CurFloor = 3
	def.CurDir = -1

	fmt.Printf("%v", queue.Update_orderlist(orderlist, 4))
	fmt.Print(queue.Cost(orderlist, 4))

	quit := make(chan int)
	go asd.Quit_program(quit)

	<-quit
}
