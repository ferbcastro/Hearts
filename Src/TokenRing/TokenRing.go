package TokenRing

/* Package types */
const ()

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
	src       byte
	dest      byte
	pkgType   byte
	data      [DATA_SIZE]byte
}

var recvPkg TokenRingPackage
var sendPkg TokenRingPackage

func CreateRing(ipArr []string) {

}

func EnterRing(ip string) {

}

func Send(dest int, msgType int, bytes []byte) {

}

/* Block until valid pkg for the calling machine arrives */
func Recv(bytes []byte, size uint) {

}

/* */
func recv() int {

}

/* */
func send() {

}
