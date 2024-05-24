package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
)

func main() {
	// Create a CPU profile file
	f, err := os.Create("profile.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	// Generate some random numbers and calculate their squares
	for i := 0; i < 1000000; i++ {
		n := rand.Intn(100)
		s := square(n)
		fmt.Printf("%d^2 = %d\n", n, s)
	}
}

func square(n int) int {
	return n * n
}
