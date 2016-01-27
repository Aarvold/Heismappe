package main


import (
    "fmt"
)

type State int

const (
    IDLE State = iota
    RUN 
    DOOR_OPEN
)

func main(){
	
	state_machine(2)

}

func state_machine(newState State){

	// Have a global varible for current_floor. current_floor == last floor 
	// Have a global varible for next_floor. Next floor will be the closest ordered floor

	switch {
    case newState == IDLE:
        fmt.Println("IDLE")
        // Set motor speed to zero
        // Wait for order
        // If order at current floor -> go to DOOR_OPEN
        // Else -> go to RUN
        return

    case newState == RUN:
        fmt.Println("RUNing")
        //
        // Set motor speed and direction 
        // If at ordered floor -> stop elev and go to DOOR_OPEN
        return

    case newState == DOOR_OPEN:
        fmt.Println("Door is OPEN")
        // Open door, wait 3 sec and then close
        // If more orders -> go to RUN
        // Else -> go to IDLE
        return 

    }
    return 
}

