package TokenRing

import (
    "log"
    "bytes"
    "encoding/gob"
)

/* This struct must be aligned */
type TokenRingPackage struct {
    TokenBusy byte
    Dest      byte
    PkgType   byte
    Size      byte
    Data      [DATA_SIZE]byte
    buffer    bytes.Buffer // unexported field to use for enconding and decoding operations
}

type bootData struct{
    Id      int    
    NewLink string 
    Next    string 
}

func (pkg *TokenRingPackage) encodeIntoDataField(s any) int {

    pkg.buffer.Reset()
    encoder := gob.NewEncoder(&pkg.buffer)
    err := encoder.Encode(s)
    if err != nil {
        log.Printf("Failed to encode data [%v]", err)
        return -1
    }

    pkg.Size = byte(len(pkg.buffer.Bytes()))
    copy(pkg.Data[:], pkg.buffer.Bytes())

    return 0
}

func (pkg *TokenRingPackage) decodeFromDataField(s any) int {
    pkg.buffer.Reset()
    pkg.buffer.Write(pkg.Data[:])

    decoder := gob.NewDecoder(&pkg.buffer)
    err := decoder.Decode(s)
    if err != nil {
        log.Printf("Failed decode data field [%v]", err)
        return -1
    }

    return 0
}


/*
The Boot Package:
In the data field we will have the following data
- id - id granted to the receiver
- newLink - the ip addr of the receiver, boot pkgs dont use dest field
- next - ip addr of the next computer in the ring
*/
func (pkg *TokenRingPackage) prepareBootPkg(newLink string, next string, id int) int {

    pkg.PkgType = BOOT
    msg := bootData{id, newLink, next}

    pkg.encodeIntoDataField(msg) 

    return 0
}

