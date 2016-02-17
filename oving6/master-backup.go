package main

import (
	"fmt"
	//"log"
	"net"
	//"os/exec"
	"strconv"
	"time"
)

const port = "20008"
const localIP = "129.241.187.161"
const broadcastIP = "129.241.187.255"

func main() {
	/*	cmd := exec.Command("sleep", "1")
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Waiting for command to finish...")
		err = cmd.Wait()
		log.Printf("Command finished with error: %v", err)
	*/
	iterChan := make(chan int)
	quitServer := make(chan int, 2)
	quitServer1 := make(chan int)

	go server(quitServer, iterChan)

	go stimer(quitServer, quitServer1, iterChan)

	<-quitServer1

	go client()

	time.Sleep(30 * time.Second)
}

func stimer(quitServer, quitServer1, iterChan chan int) {
	for {
		//fmt.Println("stimer kjører")
		timer := time.NewTimer(time.Second * 4)
		select {
		case <-timer.C:

			fmt.Println("new Master")
			quitServer <- 1
			fmt.Println("new Master1")
			quitServer1 <- 1
			fmt.Println("new Master2")

		case <-iterChan:
			//fmt.Println("En iterasjon i stimer")
		}
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func client() {
	BroadcastAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(broadcastIP, port))
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(localIP, port))
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, BroadcastAddr)
	CheckError(err)

	defer Conn.Close()
	i := 0
	for {

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

func server(quitServer, iterChan chan int) {
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
			fmt.Println("fooooor")
			n, addr, err := BroadcastConn.ReadFromUDP(buff)
			fmt.Println("etttter")
			fmt.Println("Received ", string(buff[0:n]), " from ", addr)
			iterChan <- n

			time.Sleep(500 * time.Millisecond)

			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
	}
}
