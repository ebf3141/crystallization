package fastrand

import "math/rand"
import "time"

var Random = rand.New(rand.NewSource(time.Now().UnixNano()))
var r uint64
var bits int

func Rand9() int {
	if (bits < 6) {
		r = Random.Uint64()
		bits = 64
	}
	
	bits -= 6
	val := int(r & 0x3f)
	r = r >> 6
	if val != 63{
		return val % 9
	}
	
	return Rand9()
}
