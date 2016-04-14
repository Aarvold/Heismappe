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
		if atTopOrBottom(floorSensor) {
			def.CurDir = -def.CurDir
		}
		if (len(def.Orders) > 0) && (floorSensor != -1) {
			
			setDirection(def.Orders[0],def.CurFloor)

			//def.CurFloor = driver.Get_floor_sensor_signal()
			//Hvis heisen er på vei opp/ned og er i en etasje, så fucker den opp
			if arrivedAtDestination() {
				//fmt.Printf("%sFloat cur floor = %v float def.Orders[0] = %v %s \n", def.ColY, float64(floorSensor), math.Abs(float64(def.Orders[0])), def.ColN)
				driver.Set_motor_dir(0)
				driver.Set_door_open_lamp(1)
				time.Sleep(2 * time.Second)
				driver.Set_door_open_lamp(0)
				//Kjøre en func som setter motorDir her?
				
				sendCompleteMsg(def.Orders[0])

				//Removes the first element in def.Orders
				//fmt.Printf("%sCurrent floor %d and orders[0] %d and orders = %v %s\n", def.Col0, def.CurFloor, def.Orders[0], def.Orders, def.ColN)
				def.Orders = append(def.Orders[:0], def.Orders[1:]...)
				driver.Set_button_lamp(def.BtnInside, def.CurFloor, 0)

				fmt.Printf("%sOrder deleted updated orderlist is %v %s \n", def.ColB, def.Orders, def.ColN)

			}
		}
	}
}

func setDirection(orders,curFloor int){
	if (math.Abs(float64(orders[0])) - float64(curFloor)) != 0 {
		//difference divided with abs(difference) gives direction = -1 or 1
		dir = (math.Abs(float64(orders[0])) - float64(curFloor)) / (float64(helpFunc.Difference_abs(orders[0], curFloor)))
		driver.Set_motor_dir(int(dir))
		def.CurDir = int(dir)
	}
}

func arrivedAtDestination()bool{
	return float64(floorSensor) == math.Abs(float64(def.Orders[0]))	
}

func atTopOrBottom(floorSensor int)bool{
	return (floorSensor == 0) || (floorSensor == def.NumFloors-1)
}

func sendCompleteMsg(){
	if def.Orders[0] < 0 {
		msg := def.Message{Category: def.CompleteOrder, Floor: def.CurFloor, Button: def.BtnDown, Cost: -1}
		outgoingMsg <- msg
	} else {
		msg := def.Message{Category: def.CompleteOrder, Floor: def.CurFloor, Button: def.BtnUp, Cost: -1}
		outgoingMsg <- msg
	}
}


func Update_lights_orders(outgoingMsg chan def.Message) {
	var floorSensorSignal int
	var alreadyPushed [def.NumFloors][def.NumButtons]bool
	for {
		floorSensorSignal = driver.Get_floor_sensor_signal()
		for buttontype := 0; buttontype < 3; buttontype++ {
			for floor := 0; floor < def.NumFloors; floor++ {
				if buttonPushed(buttontype,floor) {
					driver.Set_button_lamp(buttontype, floor, 1)

					if !alreadyPushed[floor][buttontype] {
						handleNewOrder(buttontype,floor,outgoingMsg)
					}
					alreadyPushed[floor][buttontype] = true
				} else {
					alreadyPushed[floor][buttontype] = false
				}
				//denne if-en hører ikke hjemme her
				if floorSensorSignal != -1 {
					def.CurFloor = floorSensorSignal
					driver.Set_floor_indicator(def.CurFloor)

				}
			}
		}
	}
}

func handleNewOrder(buttontype,floor int,outgoingMsg chan def.Message){
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

func buttonPushed(buttontype,floor int)bool{
	return driver.Get_button_signal(buttontype, floor) == 1
}

