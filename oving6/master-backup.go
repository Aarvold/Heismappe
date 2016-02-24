package main

import (
	"fmt"
	//"log"

	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const port = "20008"
const localIP = "129.241.187.161"
const broadcastIP = "129.241.187.255"

func main() {

	/*	err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Waiting for command to finish...")
		err = cmd.Wait()
		log.Printf("Command finished with error: %v", err)
	*/
	iterChan := make(chan string)
	quitServer := make(chan int, 2)
	quitServer1 := make(chan int)

	go server(quitServer, iterChan)

	go stimer(quitServer, quitServer1, iterChan)

	last_value := <-quitServer1

	go client(last_value)

	time.Sleep(10 * time.Second)
}

func stimer(quitServer, quitServer1 chan int, iterChan chan string) {
	var last_value int
	cmd := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run master-backup.go %d", strconv.Itoa(last_value))

	last_value = 0

	if len(os.Args) > 1 {
		last_value, _ = strconv.Atoi(os.Args[1])
		fmt.Printf("Initiated with value = %d \n", last_value)
	}
	for {

		timer := time.NewTimer(time.Second * 4)
		select {
		case <-timer.C:

			quitServer <- 1
			quitServer1 <- last_value
			fmt.Println("new Master1")

			cmd.Run()

		case cur_value_string := <-iterChan:

			cur_value, _ := strconv.Atoi(cur_value_string)

			if (cur_value - 1) != (last_value) {
				quitServer <- 1
				quitServer1 <- last_value
				fmt.Println("new Master2")
				cmd.Run()
			} else {
				last_value = cur_value
			}
			//fmt.Println("En iterasjon i stimer")
		}
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func client(last_value int) {
	BroadcastAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(broadcastIP, port))
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(localIP, port))
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, BroadcastAddr)
	CheckError(err)

	defer Conn.Close()
	i := last_value
	for {
		fmt.Printf("Sent %d \n", i)
		msg := strconv.Itoa(i)
		i++
		buf := []byte(msg)
		_, err := Conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}

		time.Sleep(time.Second * 1)
	}
}

func server(quitServer chan int, iterChan chan string) {
	BroadcastAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(broadcastIP, port))
	CheckError(err)

	BroadcastConn, err := net.ListenUDP("udp", BroadcastAddr)
	CheckError(err)
	defer BroadcastConn.Close()

	for {
		select {
		case <-quitServer:
			return
		default:
			buff := make([]byte, 1024)

			n, addr, err := BroadcastConn.ReadFromUDP(buff)
			fmt.Println("Received ", string(buff[0:n]), " from ", addr)
			val_recv := string(buff[0:n])

			iterChan <- val_recv

			time.Sleep(500 * time.Millisecond)

			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
	}
}
