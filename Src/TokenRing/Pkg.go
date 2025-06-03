package TokenRing

import (
    "log"
    "encoding/gob"
)

type bootData struct{
    Id      int    
    NewLink string 
    Next    string 
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

    client.buffer.Reset()
    encoder := gob.NewEncoder(&client.buffer)
    err := encoder.Encode(&msg)
    if err != nil {
        log.Printf("Failed to encode data [%v]", err)
        return -1
    }

    client.sendPkg.PkgType = BOOT
    client.sendPkg.Size = byte(len(client.buffer.Bytes()))

    copy(client.sendPkg.Data[:], client.buffer.Bytes())

    return 0
}

func (client *TokenRingClient) extractBootData() *bootData {
    client.buffer.Reset()
    client.buffer.Write(client.recvPkg.Data[:])
    var data bootData

    decoder := gob.NewDecoder(&client.buffer)
    err := decoder.Decode(&data)
    if err != nil {
        log.Printf("Failed to decode boot data [%v]", err)
        return nil
    }
    return &data
}
