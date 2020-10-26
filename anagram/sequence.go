package anagram

import (
	"math/big"
)

// Send the sequence 2, 3, 4, ... to channel 'ch'.
func generate(ch chan<- int64) {
	var i int64 = 2

	for ; ; i++ {
		if big.NewInt(i).ProbablyPrime(0) {
			ch <- i
		}
	}
}
