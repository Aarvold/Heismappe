package handleOrders

import (
	def"config"
	"queue"
	"driver"
)

var NumOfOnlineElevs int
var ImConnected bool


func Handle_orders(outgoingMsg, incomingMsg, costMsg, orderIsCompleted chan def.Message) {

	go  assign_external_order(costMsg, outgoingMsg, orderIsCompleted) 

	// Ensures that a order is only registered once for every button push
	var buttonAlreadyPushed [def.NumFloors][def.NumButtons]bool

	for {
		for buttontype := 0; buttontype < 3; buttontype++ {
			for floor := 0; floor < def.NumFloors; floor++ {
				if button_pushed(buttontype,floor) {

					if !buttonAlreadyPushed[floor][buttontype] {
						set_order_light(buttontype, floor)
						handle_new_order(buttontype,floor,outgoingMsg)
					}
					buttonAlreadyPushed[floor][buttontype] = true

				} else {
					buttonAlreadyPushed[floor][buttontype] = false	
				}
			}
		}
	}
}



func set_order_light(buttontype, floor int){
	if !external_order(buttontype) {
		driver.Set_button_lamp(buttontype, floor, 1)
	}else if ImConnected{
		driver.Set_button_lamp(buttontype, floor, 1)
	}
}

func handle_new_order(buttontype,floor int,outgoingMsg chan def.Message){
	if external_order(buttontype) && ImConnected {
		msg := def.Message{Category: def.NewOrder, Floor: floor, Button: buttontype, Cost: -1}
		outgoingMsg <- msg
	} else if !external_order(buttontype){
		// Internal orders: if the desired floor is under the elevator it is set as a order down
		if floor == driver.Get_cur_floor() && driver.Get_dir() == 1{
			queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), -floor))
		}else if floor < driver.Get_cur_floor(){
			queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), -floor))
		}else if floor == def.NumFloors-1 {
			queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), -floor))
		}else{
			queue.Set_Orders(queue.Update_orderlist(queue.Get_Orders(), floor))
		}
		
		queue.Save_backup_to_file()
	}
}

func button_pushed(buttontype,floor int)bool{
	return driver.Get_button_signal(buttontype, floor) == 1
}

func external_order(buttontype int)bool{
	return buttontype == def.BtnUp || buttontype == def.BtnDown
}



