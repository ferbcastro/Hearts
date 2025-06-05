package TokenRing

import (
	"bytes"
	"encoding/gob"
	"log"
)

/* Package types */
const (
	DATA          = iota
	BOOT          // Pkg used to bootstrap the ring
	FORWARD       // Pkg used to test the ring
	RING_COMPLETE // Pkg used to communicate the ring is complete
)

/* Other constants */
const (
	DATA_SIZE   = 128 // dunno
	TOKEN_FREE  = 0
	VALID_PKG   = 1
	FORWARD_PKG = 2
)

type TokenRingClient struct {
	id     byte
	serial byte
	ipaddr string

	sock SockDgram

	ipAddrs []string

	hasToForward bool
	waitForToken bool

	sendPkg TokenRingPackage
	recvPkg TokenRingPackage
}

// recv reads a package from the socket and decodes it into recvPkg.
func (client *TokenRingClient) recv() int {
	buffer := make([]byte, 1024)
	client.recvPkg = TokenRingPackage{}
	ret := client.sock.Recv(buffer)
	if ret <= 0 {
		log.Println("TokenRingClient: failed to receive package")
		return -1
	}

	decoder := gob.NewDecoder(bytes.NewReader(buffer[:ret]))
	if err := decoder.Decode(&client.recvPkg); err != nil {
		log.Printf("TokenRingClient: failed to decode package: %v", err)
		return -1
	}

	return ret
}

// send encodes sendPkg and writes it to the socket.
func (client *TokenRingClient) send() int {
	client.sendPkg.buffer.Reset()

	encoder := gob.NewEncoder(&client.sendPkg.buffer)
	if err := encoder.Encode(&client.sendPkg); err != nil {
		log.Printf("TokenRingClient: failed to encode package: %v", err)
		return -1
	}

	ret := client.sock.Send(client.sendPkg.buffer.Bytes())
	if ret <= 0 {
		log.Println("TokenRingClient: failed to send package")
		return -1
	}

	return ret
}

func (client *TokenRingClient) forward() int {
	client.sendPkg = client.recvPkg
	return client.send()
}

/*
TODO Implement token logic on Send and Recv
*/

func (client *TokenRingClient) Send(dest byte, data any) int {

	if client.waitForToken {
		client.Recv(nil)
	} else {
		client.waitForToken = true

	}

	client.prepareSendPkg(dest, DATA, data)

	var err int
	for {
		err = client.send()
		if err <= 0 {
			log.Printf("Failed to send data ")
			return -1
		}

		// wait for the pkg to comeback
		err = client.recv()
		if err <= 0 {
			log.Printf("Failed to recv pkg ")
			return err
		}

		// check pkg and free token
		if client.recvPkg.Src == client.id && client.recvPkg.Serial == client.serial {
			client.sendPkg.TokenBusy = 0
			err = client.send()
			if err <= 0 {
				log.Printf("Failed to send data ")
				return -1
			}
			break
		}
	}
	return err
}

/* Block until valid pkg for the calling machine arrives */
/* if passed nil recv waits for the token */
func (client *TokenRingClient) Recv(out any) {

	if client.hasToForward {
		client.forward()
		client.hasToForward = false
	}

	for {
		err := client.recv()
		if err <= 0 {
			log.Printf("Failed to receive a pk\n")
			continue
		}

		//log.Printf("pkg received: %+v\n", client.recvPkg)

		if client.recvPkg.TokenBusy == 0 {
			if out == nil {
				client.sendPkg.TokenBusy = 1
				return
			}
		} else {
			if client.recvPkg.Dest == client.id {
				err = client.recvPkg.decodeFromDataField(out)
				if err != 0 {
					log.Printf("Failed to decode pkg data\n")
					continue
				}
				client.hasToForward = true
				return
			}
		}

		client.forward()
	}
}
