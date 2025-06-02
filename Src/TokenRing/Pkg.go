package TokenRing

import (
  "log"
  "unsafe"
)

type bootData struct{
  id      int    
  newLink string 
  next    string 
}

/*
The Boot Package:
In the data field we will have the following data
  - id - id granted to the receiver
  - newLink - the ip addr of the receiver, boot pkgs dont use dest field
  - next - ip addr of the next computer in the ring
*/
func (client *TokenRingClient) prepareBootPkg(newLink string, next string, id int) int {
  
  msg := bootData{id, newLink, next}

  err := client.encoder.Encode(msg)
  if err != nil {
    log.Printf("Failed to encode data [%v]", err)
    return -1
  }

  client.sendPkg.pkgType = BOOT
  client.sendPkg.size = byte(unsafe.Sizeof(msg))

  client.buffer.Reset()
  client.buffer.Write(client.sendPkg.data[:])

  return 0
}

func (client *TokenRingClient) extractBootData() *bootData {
  client.buffer.Reset()
  client.buffer.Write(client.recvPkg.data[:])
  var data bootData

  err := client.decoder.Decode(&data)
  if err != nil {
    log.Printf("Failed to decode boot data [%v]", err)
    return nil
  }
  return &data
}
