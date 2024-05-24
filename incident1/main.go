package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type RaftBlock struct {
	ID        uint64
	Message   string
	CreatedAt time.Time
}

type Value struct {
	Timestamp time.Time
	Value     float64
}

func averageOfChan(in chan *Value) float64 {
	var sum float64
	var count int
	for v := range in {
		sum += v.Value
		count++
	}
	return sum / float64(count)
}

func averageOfSlice(in []*Value) float64 {
	var sum float64
	var count int
	for _, v := range in {
		sum += v.Value
		count++
	}
	return sum / float64(count)
}

func prepareChan() chan int {
	var count int = 10000000

	c := make(chan int, count)

	for i := 0; i < count; i++ {
		c <- i
	}
	close(c)
	return c
}
func oneChan() int64 {
	c := prepareChan()

	foundVal := true
	start := time.Now()
	for {
		select {
		case _, foundVal = <-c:
			break
		}
		if !foundVal {
			break
		}
	}
	ms := time.Since(start).Milliseconds()
	fmt.Printf("1 Chan - Standard: %dms\n", ms)
	return ms
}
func accountState() {
	ch := make(chan *RaftBlock, 1)
	ch1 := make(chan *RaftBlock, 1)
	go func() {
		for msg := range ch1 {
			fmt.Printf("Start sleep %d \n", msg.ID)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("End sleep %d \n", msg.ID)
			ch <- msg
		}
	}()
	go func() {
		for {
			time.Sleep(10 * time.Second)
			fmt.Printf("live\n")

		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	count := uint64(0)
	for {
		select {
		case block := <-ch:
			d := time.Since(block.CreatedAt)

			fmt.Printf("Finish Block %d - %d\n", block.ID, d.Microseconds())

			if block.ID == 5 {
				// n := rand.IntN(40)
				// ts := int32(n)
				// fmt.Printf("Sleep %d \n", ts)

				st := time.Duration(30) * time.Millisecond
				time.Sleep(st)

				fmt.Printf("Wake up... %d \n", block.ID)
			}

		case <-time.After(20 * time.Millisecond):
			count++

			id := count
			msg := fmt.Sprintf("Block %d", id)
			fmt.Printf("Create Block %d \n", id)

			block := &RaftBlock{
				ID:        id,
				Message:   msg,
				CreatedAt: time.Now(),
			}
			ch <- block
			fmt.Printf("Sent Block %d \n", id)

		case <-ctx.Done():
			close(ch)

			// default:
			// fmt.Printf("call runtime.Gosched \n")

			// runtime.Gosched()
		}
	}

	cancel()
}
func CompareChannelSlice() {
	// Create a large array of random numbers
	input := make([]*Value, 1e7)
	for i := 0; i < 1e7; i++ {
		input[i] = &Value{time.Unix(int64(i), 0), rand.Float64()}
	}

	func() {
		st := time.Now()
		in := make(chan *Value, 1e4)
		go func() {
			defer close(in)
			for _, v := range input {
				in <- v
			}
		}()
		value := averageOfChan(in)
		fmt.Println("Channel version took", value, time.Since(st))
	}()

	func() {
		st := time.Now()
		value := averageOfSlice(input)
		fmt.Println("Slice version took", value, time.Since(st))
	}()
}
func main() {
	accountState()
	// oneChan()
	// CompareChannelSlice()
}
