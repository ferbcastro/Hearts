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

/* This struct must be aligned */
type TokenRingPackage struct {
	tokenBusy byte
	dest      byte
	pkgType   byte
  size      byte
	data      [DATA_SIZE]byte
}

type TokenRingClient struct {
  id       int
  ipaddr   string

  sock     SockDgram

  ipAddrs  []string

  buffer   bytes.Buffer

  encoder  *gob.Encoder
  decoder  *gob.Decoder

  sendPkg  TokenRingPackage
  recvPkg  TokenRingPackage
}


/*
init token client:
  - local - the machinhe ip addr
  - dest - the ip addr of the next machine 
  - id - id granted to this client 
*/
func (client *TokenRingClient) init(ipaddr string) int {

  client.ipaddr = ipaddr 
  err := client.sock.InitSocket(client.ipaddr)
  if err != 0 {
    log.Printf("Failed to initialize socket [%v]", err)
    return -1
  }

  client.encoder = gob.NewEncoder(&client.buffer)
  client.decoder = gob.NewDecoder(&client.buffer)

  return 0
}

func (client *TokenRingClient) connect(dest string) int {

  err := client.sock.SetDest(dest)
  if err != 0 {
    log.Printf("Failed to set socket destination [%v]", err)
    return -1
  }
  return 0
}

func (client *TokenRingClient) initAsStarter(ipaddr string) int {
  err := client.init(ipaddr)
  if err != 0 {
    log.Printf("Failed to Initialize client")
    return -1
  }

  err = client.connect(client.ipAddrs[1])
  if err != 0 {
    log.Printf("Failed to connect client")
    return -1
  }

  client.id = 0

  return 0
}

func (client *TokenRingClient ) setupNext(newLink string, next string, id int) int {

  err := client.prepareBootPkg(newLink, next, id) 
  if err != 0 {
    log.Printf("Failed to prepare boot pkg")
    return -1
  }

  err = client.send()
  if err > 0 {
    log.Printf("Failed to send boot pkg")
    return err
  }

  return 0
}

// Who calls this automatically has id 0 ??
func (client *TokenRingClient) CreateRing(ipAddrs []string) int {

  client.ipAddrs = ipAddrs
  err := client.initAsStarter(ipAddrs[0])
  if err != 0 {
    log.Printf("Failed to init Token Ring Starter[%v]\n", err)
    return -1
  }

  for {
    for i := 1; i < len(ipAddrs); i++ {
      client.setupNext(ipAddrs[i], ipAddrs[i+1],i)
    }

    client.sendPkg.pkgType = FORWARD
    client.send()
    client.recv()
   
  }
}

// Return id of the caller machine ??
func (client *TokenRingClient) EnterRing(ip string) int {
  var isRingComplete bool

  client.init(ip)

  for {
    if isRingComplete {
      break
    }

    err := client.recv()
    if err == -1 {
      continue 
    }

    switch client.recvPkg.pkgType {
    case BOOT:
      bootd := client.extractBootData()

      if bootd == nil {
        log.Printf("Failed to extract data")
        continue
      }

      err = client.connect(bootd.next)
      if err != 0 {
        log.Printf("Failed to connect to the next link")
        continue
      }

      client.id = bootd.id
    case RING_COMPLETE:
      isRingComplete = true
      break
    }

    client.forward() 
  }

  return client.id
}

/* */
func (client *TokenRingClient) recv() int {
  ret := client.sock.Recv(client.buffer.Bytes())
  if ret <=0 {
    log.Printf("Failed to recv a pkg ")
    return -1
  }

  err := client.decoder.Decode(&client.recvPkg)
  if err != nil {
    log.Printf("Failed to decode pkg [%v]", err)
    return -1
  }

  return ret
}

/* */
func (client *TokenRingClient) send() int {
  err := client.encoder.Encode(client.sendPkg)
  if err != nil {
    log.Printf("Failed to encode pkg [%v]\n", err)
    return -1
  }

  ret := client.sock.Send(client.buffer.Bytes())
  if ret <= 0 {
    log.Printf("Failed to send package")
    return -1
  }

  return ret
}

func (client *TokenRingClient) forward() int {
  client.sendPkg = client.recvPkg
  return client.send()
}

func Send(dest int, msgType int, data any) {

}

/* Block until valid pkg for the calling machine arrives */
func Recv(out any) {

}


