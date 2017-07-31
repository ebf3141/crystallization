package fastrand

import "math/rand"
import "time"

var Random = rand.New(rand.NewSource(time.Now().UnixNano()))
var r uint64
var bits int

func Rand9() int {
     if (bits < 4) {
         r = Random.Uint64()
         bits = 64
     }

     bits -= 4
     val := int(r & 0x0f)
     r = r >> 4
     return val % 9
}
