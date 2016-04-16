package network

import (
	def "config"
	"fmt"
	"net"
	"strconv"
	"time"
	"log"
	"handleOrders"
)

type UdpConnection struct {
	Addr  string
	Timer *time.Timer
}
//Broadcast address
var baddr *net.UDPAddr 

type udpMessage struct {
	raddr  string //if receiving raddr=senders address, if sending raddr should be set to "broadcast" or an ip:port
	data   []byte
	length int //length of received data, in #bytes // N/A for sending
}

var closeConnectionChan = make(chan bool)

func udpInit(localListenPort, broadcastListenPort, message_size int, send_ch, receive_ch chan udpMessage) (err error) {
	//Generating broadcast address
	baddr, err = net.ResolveUDPAddr("udp4", "255.255.255.255:"+strconv.Itoa(broadcastListenPort))
	if err != nil {
		return err
	}
	//Generating local address udpSend
	tempConn, err := net.DialUDP("udp4", nil, baddr)
	if err != nil {
		handleOrders.ImConnected = false
		return err
	}
	defer tempConn.Close()
	tempAddr := tempConn.LocalAddr()
	laddr, err := net.ResolveUDPAddr("udp4", tempAddr.String())
	laddr.Port = localListenPort
	def.Laddr = laddr.String()

	//Creating local listening connections
	localListenConn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return err
	}

	//Creating listener on broadcast connection
	broadcastListenConn, err := net.ListenUDP("udp", baddr)
	if err != nil {
		localListenConn.Close()
		return err
	}
	go udp_receive_server(localListenConn, broadcastListenConn, message_size, receive_ch)
	go udp_transmit_server(localListenConn, broadcastListenConn, send_ch)
	go udp_connection_closer(localListenConn, broadcastListenConn)

	fmt.Printf("Generating local address: \t Network(): %s \t String(): %s \n", laddr.Network(), laddr.String())
	fmt.Printf("Generating broadcast address: \t Network(): %s \t String(): %s \n", baddr.Network(), baddr.String())
	return err
}

func udp_transmit_server(lconn, bconn *net.UDPConn, send_ch <-chan udpMessage) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR in udp_transmit_server: %s \n Closing connection.", r)
			closeConnectionChan <- true
		}
	}()

	var err error
	var n int

	for {
		//fmt.Printf("udp_transmit_server: waiting on new value on Global_Send_ch \n")
		msg := <-send_ch
		if msg.raddr == "broadcast" {
			n, err = lconn.WriteToUDP(msg.data, baddr)
		} 
		if err != nil || n < 0 {
			fmt.Printf("Error: udp_transmit_server: writing\n")
			panic(err)
		}
		//fmt.Printf("udp_transmit_server: Sent %s to %s \n", msg.data, msg.raddr)
	}
}

func udp_receive_server(lconn, bconn *net.UDPConn, message_size int, receive_ch chan<- udpMessage) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR in udp_receive_server: %s \n Closing connection.", r)
			closeConnectionChan <- true
		}
	}()

	bconn_rcv_ch := make(chan udpMessage)
	lconn_rcv_ch := make(chan udpMessage)

	go udp_connection_reader(lconn, message_size, lconn_rcv_ch)
	go udp_connection_reader(bconn, message_size, bconn_rcv_ch)

	for {
		select {

		case buf := <-bconn_rcv_ch:
			receive_ch <- buf

		case buf := <-lconn_rcv_ch:
			receive_ch <- buf
		}
	}
}

func udp_connection_reader(conn *net.UDPConn, message_size int, rcv_ch chan<- udpMessage) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR in udp_connection_reader: %s \n Closing connection.", r)
			conn.Close()
		}
	}()

	for {
		buf := make([]byte, message_size)
		//fmt.Printf("udp_connection_reader: Waiting on data from UDPConn\n")
		n, raddr, err := conn.ReadFromUDP(buf)
		//fmt.Printf("udp_connection_reader: Received %s from %s \n", string(buf), raddr.String())
		if err != nil || n < 0 {
			fmt.Printf("Error: udp_connection_reader: reading\n")
			panic(err)
		}
		rcv_ch <- udpMessage{raddr: raddr.String(), data: buf, length: n}
	}
}

func udp_connection_closer(lconn, bconn *net.UDPConn) {
	<-closeConnectionChan
	handleOrders.ImConnected = false
	fmt.Printf("%sDisconnected \n%s",def.ColR,def.ColN)
	lconn.Close()
	bconn.Close()
}
