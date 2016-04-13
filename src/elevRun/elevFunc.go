package elevRun

import (
	def "config"
	"driver"
	"fmt"
	"helpFunc"
	"math"
	"queue"
	"time"
)

func Run_elev(outgoingMsg chan def.Message) {
	var dir float64
	driver.Elev_init()

	for {
		floorSensor := driver.Get_floor_sensor_signal()
		if (floorSensor == 0) || (floorSensor == def.NumFloors-1) {
			def.CurDir = -def.CurDir
		}
		if (len(def.Orders) > 0) && (floorSensor != -1) {
			if (math.Abs(float64(def.Orders[0])) - float64(def.CurFloor)) != 0 {
				//difference divided with abs(difference) gives direction = -1 or 1
				dir = (math.Abs(float64(def.Orders[0])) - float64(def.CurFloor)) / (float64(helpFunc.Difference_abs(def.Orders[0], def.CurFloor)))
				driver.Set_motor_dir(int(dir))
				def.CurDir = int(dir)
			}

			//def.CurFloor = driver.Get_floor_sensor_signal()
			//Hvis heisen er på vei opp/ned og er i en etasje, så fucker den opp
			if float64(floorSensor) == math.Abs(float64(def.Orders[0])) {
				//fmt.Printf("%sFloat cur floor = %v float def.Orders[0] = %v %s \n", def.ColY, float64(floorSensor), math.Abs(float64(def.Orders[0])), def.ColN)
				driver.Set_motor_dir(0)
				driver.Set_door_open_lamp(1)
				time.Sleep(2 * time.Second)
				driver.Set_door_open_lamp(0)
				//Kjøre en func som setter motorDir her?
				if def.Orders[0] < 0 {
					msg := def.Message{Category: def.CompleteOrder, Floor: def.CurFloor, Button: def.BtnDown, Cost: -1}
					outgoingMsg <- msg
				} else {
					msg := def.Message{Category: def.CompleteOrder, Floor: def.CurFloor, Button: def.BtnUp, Cost: -1}
					outgoingMsg <- msg
				}

				//Removes the first element in def.Orders
				//fmt.Printf("%sCurrent floor %d and orders[0] %d and orders = %v %s\n", def.Col0, def.CurFloor, def.Orders[0], def.Orders, def.ColN)
				def.Orders = append(def.Orders[:0], def.Orders[1:]...)

				driver.Set_button_lamp(def.BtnInside, def.CurFloor, 0)

				fmt.Printf("%sOrder deleted updated orderlist is %v %s \n", def.ColB, def.Orders, def.ColN)

			}
		}
	}
}

func Update_lights_orders(outgoingMsg chan def.Message) {
	var floorSensorSignal int
	var buttonState [def.NumFloors][def.NumButtons]bool
	for {
		floorSensorSignal = driver.Get_floor_sensor_signal()
		for buttontype := 0; buttontype < 3; buttontype++ {
			for floor := 0; floor < def.NumFloors; floor++ {
				if driver.Get_button_signal(buttontype, floor) == 1 {
					driver.Set_button_lamp(buttontype, floor, 1)

					if !buttonState[floor][buttontype] {

						if buttontype == def.BtnUp || buttontype == def.BtnDown {
							msg := def.Message{Category: def.NewOrder, Floor: floor, Button: buttontype, Cost: -1}
							outgoingMsg <- msg
						} else {
							//Internal orders: if the desired floor is under the elevator it is set as a order down
							//fmt.Printf("Intern ordre mottat \n")
							if floor < def.CurFloor {
								def.Orders = queue.Update_orderlist(def.Orders, -floor, false)
							} else {
								def.Orders = queue.Update_orderlist(def.Orders, floor, false)
							}
							fmt.Printf("%sInternal: Order list is updated to %v %s \n", def.ColR, def.Orders, def.ColN)
						}
					}
					buttonState[floor][buttontype] = true
				} else {
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
