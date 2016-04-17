package elevRun

import (
	def "config"
	"driver"
	"fmt"
	"helpFunc"
	"math"
	"queue"
	"time"
	"handleOrders"

)


func Elev_init(){
	fmt.Printf("%sInitialising elevator, please wait...%s\n", def.ColG, def.ColN)
	driver.Elev_init()
	driver.Set_motor_dir(-1)
	defer driver.Set_motor_dir(0)

	queue.Get_backup_from_file()

	for {
		if driver.Get_floor_sensor_signal() != -1 {
			fmt.Printf("%sInitialising was successful%s\n", def.ColG, def.ColN)
			return
		}
	}
}

func Run_elev(outgoingMsg chan def.Message) {
	for {

		update_current_floor()

		//If we have orders and we are not between floors
		if (len(queue.Get_Orders()) > 0) && (driver.Get_cur_floor() != -1) {
			set_direction(queue.Get_Orders()[0],driver.Get_cur_floor())

			if arrived_at_destination() {
				driver.Set_motor_dir(0)

				if handleOrders.ImConnected{
					send_order_complete_msg(queue.Get_Orders()[0],outgoingMsg)
				}else{
					turn_off_lights()
				}
				//fmt.Printf("%sOrder to floor %d deleted, current floor = %d\n",def.ColB,queue.Get_Orders()[0],driver.Get_cur_floor())
				queue.Remove_first_element_in_orders()
				queue.Save_backup_to_file()
				driver.Set_button_lamp(def.BtnInside, driver.Get_cur_floor(), 0)
				open_and_close_door()
			}
		}
	}
}

func turn_off_lights(){
	if queue.Get_Orders()[0] < 0 {
		driver.Set_button_lamp(def.BtnDown, driver.Get_cur_floor(), 0)
	} else {
		driver.Set_button_lamp(def.BtnUp, driver.Get_cur_floor(), 0)
    }
}

func open_and_close_door(){
	driver.Set_door_open_lamp(1)
	time.Sleep(3 * time.Second)
	driver.Set_door_open_lamp(0)
}

func set_direction(order,curFloor int){
	if (math.Abs(float64(order)) - float64(curFloor)) != 0 {
		//difference between cur floor and desired floor divided with abs(difference) gives direction = -1 or 1
		dir := (math.Abs(float64(order)) - float64(curFloor)) / (float64(helpFunc.Difference_abs(order, curFloor)))
		driver.Set_motor_dir(int(dir))
		driver.Set_dir(int(dir))
	}
}

func arrived_at_destination()bool{
	return driver.Get_cur_floor() == int(math.Abs(float64(queue.Get_Orders()[0])))	
}

func send_order_complete_msg(order int,outgoingMsg chan def.Message){
	//if order down
	if order < 0 {
		msg := def.Message{Category: def.CompleteOrder, Floor: driver.Get_cur_floor(), Button: def.BtnDown, Cost: -1}
		outgoingMsg <- msg
	} else {
		msg := def.Message{Category: def.CompleteOrder, Floor: driver.Get_cur_floor(), Button: def.BtnUp, Cost: -1}
		outgoingMsg <- msg
	}
}

func update_current_floor(){
	floorSensorSignal := driver.Get_floor_sensor_signal()
	if floorSensorSignal != -1 {
		driver.Set_cur_floor(floorSensorSignal)
		driver.Set_floor_indicator(floorSensorSignal)
	}
}


