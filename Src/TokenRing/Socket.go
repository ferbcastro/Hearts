package TokenRing

import (
	"log"
	"net"
)

const (
	PORT = ":8080"
)

type SockDgram struct {
	localAddr *net.UDPAddr
	destAddr  *net.UDPAddr
	conn      *net.UDPConn
}

func InitSocket(sock *SockDgram, localIp, destIp string) int {
	var err error

	sock.localAddr, err = net.ResolveUDPAddr("udp4", localIp+PORT)
	if err != nil {
		log.Printf("ResolveUDPAddr failed to solve for local ip [%v]\n", err)
		return 1
	}
	sock.destAddr, err = net.ResolveUDPAddr("udp4", destIp+PORT)
	if err != nil {
		log.Printf("ResolveUDPAddr failed to solve for dest ip [%v]\n", err)
		return 1
	}
	sock.conn, err = net.ListenUDP("udp4", sock.localAddr)
	if err != nil {
		log.Printf("ListenUDP failed [%v]\n", err)
	}

	return 0
}

func CloseSocket(sock *SockDgram) int {
	sock.conn.Close()
	log.Println("Closing UDP socket")
	return 0
}

func Recv(sock *SockDgram, arr []byte) int {
	numBytes, _, err := sock.conn.ReadFromUDP(arr)
	if err != nil {
		log.Printf("ReadFromUDP failed [%v] read [%v] bytes\n", err, numBytes)
		return 1
	}
	return 0
}

func Send(sock *SockDgram, arr []byte) int {
	numBytes, err := sock.conn.WriteToUDP(arr, sock.destAddr)
	if err != nil {
		log.Printf("WriteToUDP failed [%v] written [%v] bytes\n", err, numBytes)
		return 1
	}
	return 0
}
