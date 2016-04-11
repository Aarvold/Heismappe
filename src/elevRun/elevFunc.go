package elevRun

import (
	"driver"
	"fmt"
	def "config"
	"math"
	"time"
	"helpFunc"
	"queue"
)


func Run_elev() {
	var dir float64
	driver.Elev_init()

	for {
		if len(def.Orders)>0{
			if ((def.Orders[0] - def.CurFloor)!=0) {
				//difference divided with abs(difference) gives direction = -1 or 1
				dir = (math.Abs(float64(def.Orders[0])) - float64(def.CurFloor)) / float64(helpFunc.DifferenceAbs(def.Orders[0],def.CurFloor))
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

func Update_lights_orders(outgoingMsg chan def.Message) {
	var floorSensorSignal int
	var buttonState[def.NumFloors][def.NumButtons] bool
	for {
		floorSensorSignal = driver.Get_floor_sensor_signal()
		for buttontype := 0; buttontype < 3; buttontype++ {
			for floor := 0; floor < def.NumFloors; floor++ {
				if driver.Get_button_signal(buttontype, floor) == 1 {
					driver.Set_button_lamp(buttontype, floor, 1)

					if !buttonState[floor][buttontype]{

						if(buttontype == def.BtnUp || buttontype == def.BtnDown){
							msg := def.Message{Category: def.NewOrder, Floor: floor, Button: buttontype, Cost: queue.Cost(def.Orders,floor)}
							outgoingMsg <- msg
						} else {
							//If the desired floor is under the elevator it is set as a order down
							if floor < def.CurFloor{
								def.Orders = queue.Update_orderlist(def.Orders,-floor)
							}else {
								def.Orders = queue.Update_orderlist(def.Orders,floor)
							}
						}
					}
					buttonState[floor][buttontype] = true
				} else{
					buttonState[floor][buttontype] = false
				}

				if floorSensorSignal != -1 {
					def.CurFloor = floorSensorSignal
					driver.Set_floor_indicator(def.CurFloor)
					
				}
			}
		}
	}
}






/*

func DifferenceAbs(val1 ,val2 int) int{
	return int(math.Abs(math.Abs(float64(val1))-math.Abs(float64(val2))))
}
*/