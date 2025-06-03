package TokenRing

import (
    "bytes"
    "encoding/gob"
    "log"
)

/* Package types */
const (
    DATA = iota
    BOOT // Pkg used to bootstrap the ring
    FORWARD // Pkg used to test the ring 
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
    id       int
    ipaddr   string

    sock     SockDgram

    ipAddrs  []string

    sendPkg  TokenRingPackage
    recvPkg  TokenRingPackage
}

// init initializes the client socket, read buffer, and gob encoder.
// ipaddr: the local IP address of this machine.
func (client *TokenRingClient) init(ipaddr string) int {
    client.ipaddr = ipaddr

    if err := client.sock.InitSocket(ipaddr); err != 0 {
        log.Printf("Failed to initialize socket on [%s]", ipaddr)
        return -1
    }


    // Register all types that may be transmitted
    gob.Register(TokenRingPackage{})
    gob.Register(bootData{})

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

    client.id = 0 // Starter node always gets ID 0

    return 0
}

// setupNext prepares and sends a boot package to initialize a new node in the ring.
// - newLink: the IP address of the new node being added
// - next: the IP address of the next node in the ring
// - id: the ID to assign to the new node
func (client *TokenRingClient) setupNext(newLink string, next string, id int) int {
    // Prepare the boot package with the given information
    //err := client.sendPkg.prepareBootPkg(newLink, next, id)
    err := client.Send(0, BOOT,  bootData{ newLink, next, id})
    if err <= 0 {
        log.Printf("Failed to prepare boot package for %s -> %s (ID %d)", newLink, next, id)
        return -1
    }
/*
    // Send the boot package over the socket
    err = client.send()
    if err <= 0 {
        log.Printf("Failed to send boot package to %s", newLink)
        return -1
    }
*/
    return 0
}

// CreateRing initializes the Token Ring network by connecting all clients in a circular list.
// The caller is assumed to be the starter and assigned ID 0.
func (client *TokenRingClient) CreateRing(ipAddrs []string) int {
    client.ipAddrs = ipAddrs

    // Initialize the starter (caller) with ID 0 and bind to ipAddrs[0]
    err := client.initAsStarter(ipAddrs[0])
    if err != 0 {
        log.Printf("Failed to init Token Ring starter: %v\n", err)
        return -1
    }

    for {
        // Setup the ring by bootstrapping each client with its ID and next IP
        for i := 1; i < len(ipAddrs); i++ {
            nextIdx := (i + 1) % len(ipAddrs)
            log.Printf("Setting up link: %s (ID %d) -> %s\n", ipAddrs[i], i, ipAddrs[nextIdx])

            err := client.setupNext(ipAddrs[i], ipAddrs[nextIdx], i)
            if err != 0 {
                log.Printf("Failed to setup next link: %v\n", err)
                return -1
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

    return 0
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
            return client.id
        }
    }
}

// recv reads a package from the socket and decodes it into recvPkg.
func (client *TokenRingClient) recv() int {
    buffer := make([]byte, 1024)
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

func (client *TokenRingClient) Send(dest int, msgType int, data any) int {

  client.sendPkg.PkgType = msgType 
  client.sendPkg.Dest = dest
  err := client.sendPkg.encodeIntoDataField(data)
  if err != 0 {
    log.Printf("Failed to encode data into dataField")
    return -1
  }

  err = client.send()
  if err <= 0 {
    log.Printf("Failed to encode data into dataField")
  }
  
  return err
}

/* Block until valid pkg for the calling machine arrives */
func Recv(out any) {

}


