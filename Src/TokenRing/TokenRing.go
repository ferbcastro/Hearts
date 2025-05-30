package TokenRing

/* Package types */
const ()

/* Other constants */
const (
	DATA_SIZE = 128 // dunno
)

/* This struct must be aligned */
type TokenRingPackage struct {
	tokenBusy byte
	dest      byte
	pkgType   byte
	data      [DATA_SIZE]byte
}

// change if needed
var recvPkg TokenRingPackage
var sendPkg TokenRingPackage

func CreateRing(ipArr []string) {

}

func EnterRing(ip string) {

}

func Send(dest int, msgType int, bytes []byte) {

}

/* Block until valid pkg for the calling machine arrives. Always forward pkg */
func Recv(bytes []byte, size uint) {

}
