package util

import (
	"math/rand"
)

// Function to implement the Miller-Rabin Primality Test
func IsPrime(n int, k int) bool {
	if n == 2 || n == 3 {
		return true
	}

	if n <= 1 || n%2 == 0 {
		return false
	}

	// Find r and d such that n = 2^r * d + 1
	r, d := 0, n-1
	for d%2 == 0 {
		r++
		d /= 2
	}

	// Repeat k times
	for i := 0; i < k; i++ {
		a := rand.Intn(n-2) + 2
		x := PowMod(a, d, n)

		if x == 1 || x == n-1 {
			continue
		}

		j := 0
		for ; j < r; j++ {
			x = PowMod(x, 2, n)
			if x == n-1 {
				break
			}
		}

		if j == r {
			return false
		}
	}

	return true
}

func PowMod(x, y, m int) int {
	r := 1
	x = x % m
	for y > 0 {
		if y%2 == 1 {
			r = (r * x) % m
		}
		y = y / 2
		x = (x * x) % m
	}

	return r
}
