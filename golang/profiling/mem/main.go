package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	// Create a CPU profile file
	cpuProfileFile, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer cpuProfileFile.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(cpuProfileFile); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	// Create a memory profile file
	memProfileFile, err := os.Create("mem.prof")
	if err != nil {
		panic(err)
	}
	defer memProfileFile.Close()

	// Generate some random numbers and calculate their squares
	for i := 0; i < 10; i++ {
		n := rand.Intn(100)
		s := square(n)
		fmt.Printf("%d^2 = %d\n", n, s)
	}

	// Write memory profile to file
	if err := pprof.WriteHeapProfile(memProfileFile); err != nil {
		panic(err)
	}
	fmt.Println("Memory profile written to mem.prof")

	time.Sleep(5 * time.Second)
}

func square(n int) int {
	return n * n
}
