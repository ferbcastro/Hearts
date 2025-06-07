package TokenRing

import (
	"bytes"
	"encoding/gob"
	"log"
)

/* This struct must be aligned */
type TokenRingPackage struct {
	TokenBusy byte
	Ack       byte
	Src       byte
	Dest      byte
	Serial    byte
	PkgType   byte
	Size      byte
	Data      [DATA_SIZE]byte
	buffer    bytes.Buffer // unexported field to use for enconding and decoding operations
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

func (client *TokenRingClient) prepareSendPkg(dest byte, msgType int, data any) int {
	client.sendPkg.Src = client.id
	client.sendPkg.TokenBusy = 1 
	client.sendPkg.Ack = 0 
	client.sendPkg.Dest = dest
	client.sendPkg.PkgType = byte(msgType)
	client.serial++
	client.sendPkg.Serial = client.serial
	err := client.sendPkg.encodeIntoDataField(data)
	if err != 0 {
		log.Printf("Failed to encode data into dataField")
		return -1
	}
	return 0
}
