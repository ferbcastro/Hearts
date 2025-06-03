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
  id       byte
  ipaddr   string

  sock     SockDgram

  ipAddrs  []string

  sendPkg  TokenRingPackage
  recvPkg  TokenRingPackage
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

  client.sendPkg.PkgType = byte(msgType)
  client.sendPkg.Dest = byte(dest)
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
func (client *TokenRingClient) Recv(out any) {
  for {
    err := client.recv()
    if err <= 0 {
      log.Printf("Failed to receive a pk\n") 
      continue
    }

    if client.recvPkg.Dest == client.id {
      err = client.recvPkg.decodeFromDataField(out)  
      if err != 0 {
        log.Printf("Failed to decode pkg data\n") 
        continue
      }

      return 
    }
    client.forward()
  }
}


