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


func Elev_init(){
	driver.Elev_init()
	driver.Set_motor_dir(-1)
	defer driver.Set_motor_dir(0)

	queue.Get_backup_from_file()

	for {
		if driver.Get_floor_sensor_signal() != -1 {
			return
		}
	}

}

func Run_elev(outgoingMsg chan def.Message) {
	for {
		update_current_floor()
		floorSensor := driver.Get_floor_sensor_signal()
		if at_top_or_bottom(floorSensor) {
			def.CurDir = -def.CurDir
		}
		if (len(queue.Get_Orders()) > 0) && (floorSensor != -1) {
			//var orders = queue.Get_Orders()
			//fmt.Printf("Orders = %v \n",orders)
			set_direction(queue.Get_Orders()[0],def.CurFloor)

			//def.CurFloor = driver.Get_floor_sensor_signal()
			//Hvis heisen er på vei opp/ned og er i en etasje, så fucker den opp
			if arrived_at_destination() {
				//fmt.Printf("%sFloat cur floor = %v float orders[0] = %v %s \n", def.ColY, float64(floorSensor), math.Abs(float64(def.Orders[0])), def.ColN)
				
				send_complete_msg(queue.Get_Orders()[0],outgoingMsg)
				//Removes the first element in def.Orders
				fmt.Printf("%sOrder to floor %d deleted, current floor = %d\n",def.ColB,queue.Get_Orders()[0],def.CurFloor)
				queue.Set_Orders(append(queue.Get_Orders()[:0], queue.Get_Orders()[1:]...))
				queue.Save_backup_to_file()
				driver.Set_button_lamp(def.BtnInside, def.CurFloor, 0)


				driver.Set_motor_dir(0)
				driver.Set_door_open_lamp(1)
				time.Sleep(2 * time.Second)
				driver.Set_door_open_lamp(0)
				//Kjøre en func som setter motorDir her?
				

			}
		}
	}
}

func set_direction(order,curFloor int){
	//var dir float64
	if (math.Abs(float64(order)) - float64(curFloor)) != 0 {
		//difference divided with abs(difference) gives direction = -1 or 1
		dir := (math.Abs(float64(order)) - float64(curFloor)) / (float64(helpFunc.Difference_abs(order, curFloor)))
		driver.Set_motor_dir(int(dir))
		def.CurDir = int(dir)
		//fmt.Printf(def.ColY,"Current dir = %d \n",def.CurDir,def.ColN)
	}
}

func arrived_at_destination()bool{
	return def.CurFloor == int(math.Abs(float64(queue.Get_Orders()[0])))	
}

func at_top_or_bottom(floorSensor int)bool{
	return (floorSensor == 0) || (floorSensor == def.NumFloors-1)
}

func send_complete_msg(order int,outgoingMsg chan def.Message){
	if order < 0 {
		msg := def.Message{Category: def.CompleteOrder, Floor: def.CurFloor, Button: def.BtnDown, Cost: -1}
		outgoingMsg <- msg
	} else {
		msg := def.Message{Category: def.CompleteOrder, Floor: def.CurFloor, Button: def.BtnUp, Cost: -1}
		outgoingMsg <- msg
	}
}



func Update_lights_orders(outgoingMsg chan def.Message) {
	var alreadyPushed [def.NumFloors][def.NumButtons]bool

	for {
		for buttontype := 0; buttontype < 3; buttontype++ {
			for floor := 0; floor < def.NumFloors; floor++ {
				if button_pushed(buttontype,floor) {
					driver.Set_button_lamp(buttontype, floor, 1)

					if !alreadyPushed[floor][buttontype] {
						handle_new_order(buttontype,floor,outgoingMsg)
					}
					alreadyPushed[floor][buttontype] = true
				} else {
					alreadyPushed[floor][buttontype] = false
				}
			}
		}
	}
}

func update_current_floor(){
	var floorSensorSignal = driver.Get_floor_sensor_signal()
	if floorSensorSignal != -1 {
		def.CurFloor = floorSensorSignal
		driver.Set_floor_indicator(floorSensorSignal)
	}
}

func handle_new_order(buttontype,floor int,outgoingMsg chan def.Message){
	if buttontype == def.BtnUp || buttontype == def.BtnDown {
		msg := def.Message{Category: def.NewOrder, Floor: floor, Button: buttontype, Cost: -1}
		outgoingMsg <- msg
	} else {
		//Internal orders: if the desired floor is under the elevator it is set as a order down
		//fmt.Printf("Intern ordre mottat \n")
		if floor < def.CurFloor {
			queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), -floor, false))
		}else if floor == def.NumFloors-1 {
			queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), -floor, false))
		}else{
			queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), floor, false))
		}

		queue.Save_backup_to_file()
	}
}

func button_pushed(buttontype,floor int)bool{
	return driver.Get_button_signal(buttontype, floor) == 1
}

