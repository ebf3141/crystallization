package fastrand

import "time"
import "math/rand"

var seed uint64 = uint64(time.Now().UnixNano())
var Random = rand.New(rand.NewSource(time.Now().UnixNano()))
/*
 * This is the SplitMix64 algorithm used by the SpittableRandom class
 * as modified by https://github.com/vpxyz/xorshift
 * It generates 64 random bits on each call. We want an integer
 * between 0 and 8, but instead of using modulo 9, we instead take
 * the low 32 bits, multiply by 9 and then shift right 32, which
 * is faster than the modulo operation.
 */
func Rand9() int {
    seed = seed + uint64(0x9E3779B97F4A7C15)
    z := seed
    z = (z ^ (z >> 30)) * uint64(0xBF58476D1CE4E5B9)
    z = (z ^ (z >> 27)) * uint64(0x94D049BB133111EB)
    return int(((z ^ (z >> 31)) & 0xFFFFFFFF) * 9 >> 32)
}
