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
	dest      byte
	pkgType   byte
	data      [DATA_SIZE]byte
}

var recvPkg TokenRingPackage
var sendPkg TokenRingPackage

// Who calls this automatically has id 0 ??
func CreateRing(ipArr []string) {

}

// Return id of the caller machine ??
func EnterRing(ip string) int {

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
