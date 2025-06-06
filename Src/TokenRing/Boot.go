package TokenRing

import (
	"encoding/gob"
	"log"
)

// Structure used to send information needed to setup links
type bootData struct {
	Id      byte
	NewLink string
	Next    string
}

// init initializes the client socket, read buffer, and gob encoder.
// ipaddr: the local IP address of this machine.
func (client *TokenRingClient) init(ipaddr string) int {
	client.ipaddr = ipaddr

	if err := client.sock.InitSocket(ipaddr); err != 0 {
		log.Printf("Failed to initialize socket on [%s]", ipaddr)
		return -1
	}

	client.waitForToken = true

	// Register all types that may be transmitted
	gob.Register(TokenRingPackage{})
	gob.Register(bootData{})
	gob.Register(Message{})

	return 0
}

// connect sets the destination IP address for outgoing packets.
func (client *TokenRingClient) connect(dest string) int {
	if err := client.sock.SetDest(dest); err != 0 {
		log.Printf("Failed to set socket destination to [%s]", dest)
		return -1
	}
	return 0
}

// initAsStarter sets up the current client as the starter node in the ring.
// It initializes local state and connects to the second node in the list.
func (client *TokenRingClient) initAsStarter(ipaddr string) int {
	// Initialize local socket and encoder
	if err := client.init(ipaddr); err != 0 {
		log.Println("Failed to initialize starter client")
		return -1
	}

	// Connect to the next node in the ring (second address in the list)
	if err := client.connect(client.ipAddrs[1]); err != 0 {
		log.Println("Failed to connect to the first link in the ring")
		return -1
	}

	client.id = 0
	client.waitForToken = false

	return 0
}

// CreateRing initializes the Token Ring network by connecting all clients in a circular list.
// The caller is assumed to be the starter and assigned ID 0.
func (client *TokenRingClient) CreateRing(ipAddrs []string) []byte {
	client.ipAddrs = ipAddrs

	ids := make([]byte, len(ipAddrs))
	ids[0] = 0
	// Initialize the starter (caller) with ID 0 and bind to ipAddrs[0]
	err := client.initAsStarter(ipAddrs[0])
	if err != 0 {
		log.Printf("Failed to init Token Ring starter: %v\n", err)
		return nil
	}

	for {
		// Setup the ring by bootstrapping each client with its ID and next IP
		for i := 1; i < len(ipAddrs); i++ {
			ids[i] = byte(i)

			nextIdx := (i + 1) % len(ipAddrs)
			log.Printf("Setting up link: %s (ID %d) -> %s\n", ipAddrs[i], i, ipAddrs[nextIdx])

			data := bootData{byte(i), ipAddrs[i], ipAddrs[nextIdx]}
			client.prepareSendPkg(0, BOOT, data)
			err := client.send()
			if err <= 0 {
				log.Printf("Failed to prepare boot package for %s -> %s (ID %d)", ipAddrs[i], ipAddrs[nextIdx], i)
				return nil
			}
		}

		// Verify the ring is complete via a FORWARD token loop
		// Send FORWARD token into the ring
		client.sendPkg.PkgType = FORWARD
		client.send()

		// Wait to receive it back
		client.recv()
		if client.recvPkg.PkgType == FORWARD {
			log.Println("Ring verified successfully.")
			break
		}
	}

	// Signal that the ring is now complete
	client.sendPkg.PkgType = RING_COMPLETE
	client.send()
	client.recv()

	return ids
}

// EnterRing allows a node to join an existing token ring.
// The node waits to receive a BOOT package containing its ID and next hop,
// connects to the next node, and waits for the ring completion signal.
// Return caller id
func (client *TokenRingClient) EnterRing(ip string) int {
	// Initialize the client
	err := client.init(ip)
	if err != 0 {
		log.Printf("Failed to initialize client at IP [%s]", ip)
		return -1
	}

	for {
		// Wait for a valid package
		err = client.recv()
		if err == -1 {
			continue
		}

		switch client.recvPkg.PkgType {
		case BOOT:
			var bootd bootData
			err = client.recvPkg.decodeFromDataField(&bootd)
			if err != 0 {
				log.Println("Failed to extract boot data")
				continue
			}

			log.Printf("Received BOOT: %+v", bootd)

			// If this boot package is not for this client, forward it
			if bootd.NewLink != client.ipaddr {
				client.forward()
				continue
			}

			// Connect to the next node in the ring
			err = client.connect(bootd.Next)
			if err != 0 {
				log.Println("Failed to connect to next node in ring")
				continue
			}

			// Save the assigned ID
			client.id = bootd.Id

		case FORWARD:
			// Forward the token as part of ring verification
			client.forward()

		case RING_COMPLETE:
			// Ring setup is complete
			log.Println("Received RING_COMPLETE, entering ring")
			client.forward()
			return int(client.id)
		}
	}
}
