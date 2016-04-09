package asd

import (
	"driver"
	"fmt"
	def "config"
	"math"
	"time"
	"queue"
)

func Go_to_floor() {
	var dir float64
	for {
		if len(def.Orders)>0{
			if ((def.Orders[0] - def.CurFloor)!=0) {
				//difference divided with abs(difference) gives direction = -1 or 1
				dir = (float64(math.Abs(float64(def.Orders[0])) - float64(def.CurFloor))) / math.Abs(float64(math.Abs(float64(def.Orders[0])) - float64(def.CurFloor)))
				driver.Set_motor_dir(int(dir))
				def.CurDir = int(dir)

			}

			//def.CurFloor = driver.Get_floor_sensor_signal()
			if float64(def.CurFloor) == math.Abs(float64(def.Orders[0])) {
				driver.Set_motor_dir(0)
				driver.Set_door_open_lamp(1)
				time.Sleep(2 * time.Second)
				driver.Set_door_open_lamp(0)
				//Removes the first element in def.Orders
				fmt.Printf("%v",def.Orders)
				def.Orders=append(def.Orders[:0],def.Orders[1:]...)
				driver.Set_button_lamp(def.BtnUp, def.CurFloor, 0)
				driver.Set_button_lamp(def.BtnDown, def.CurFloor, 0)
				driver.Set_button_lamp(def.BtnInside, def.CurFloor, 0)
				fmt.Printf("%v \n",def.Orders)
			}
		}
	}
}

func Update_lights(outgoingMsg <-chan def.Message) {
	var floorSensorSignal int
	for {
		floorSensorSignal = driver.Get_floor_sensor_signal()
		for buttontype := 0; buttontype < 3; buttontype++ {
			for floor := 0; floor < def.NumFloors; floor++ {
				if driver.Get_button_signal(buttontype, floor) == 1 {
					driver.Set_button_lamp(buttontype, floor, 1)
					if(buttontype == def.BtnUp || buttontype == def.BtnDown){
						msg := def.Message{Category: def.NewOrder, Floor: floor, Button: buttontype, Cost: queue.Cost(def.Orders,floor)}
						//outgoingMsg <- msg

						fmt.Print("aefaef")
					} else {

						//If the desired floor is under the elevator it is set as a order down
						if floor < def.CurFloor{
							def.Orders = queue.Update_orderlist(def.Orders,-floor)
						}else {
							def.Orders = queue.Update_orderlist(def.Orders,floor)
						}
						
					}

				}
				if floorSensorSignal != -1 {
					def.CurFloor = floorSensorSignal
					driver.Set_floor_indicator(def.CurFloor)
					
				}
				time.Sleep(1 * time.Millisecond)
			}
		}
	}
}


func Quit_program(quit chan int) {
	for {
		time.Sleep(time.Second)
		if driver.Get_stop_signal() == 1 {
			quit <- 1
		}
	}
}
/*
func Run_elev(){
	for{
		if len(def.Orders)>0{
			Go_to_floor(def.Orders[0])

		}
	}
}
*/