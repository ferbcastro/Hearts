package TokenRing

/* Package types */
const ()

/* Other constants */
const (
	DATA_SIZE = 32
)

/* This struct must be aligned */
type TokenRingPackage struct {
	tokenBusy byte
	dest      byte
	pkgType   byte
	data      [DATA_SIZE]byte
}

func CreateRing() {

}

func EnterRing() {

}

func Send() {

}

/* Block until valid pkg for the calling machine arrives. Forward pkg with
 * different destination */
func Recv() {

}
