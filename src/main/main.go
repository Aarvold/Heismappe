package main

//export GOPATH=$HOME/Documents/MagAnd/Heismappe/

import (
	"asd"
	"driver"
	//"fmt"
	//"time"
	def "config"
	//"network"
)

func main() {
	driver.Elev_init()
	def.CurFloor = driver.Get_floor_sensor_signal()
	def.CurDir = 0
	go asd.Update_lights()

	//outgoingMsg := make(chan def.Message)
	//incomingMsg := make(chan def.Message)
	//go network.Init(outgoingMsg, incomingMsg)

	var orderlist = []int{5, 4, 2, 3, 4, 6}
	def.CurFloor = 7
	def.CurDir = -1

	print(asd.Cost(orderlist))

	quit := make(chan int)
	go asd.Quit_program(quit)

	<-quit
}
