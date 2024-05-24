package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	fmt.Println("Start")
	// get_rid_of_share_loop_variables()
	// range_integer()
	defer_time()
	fmt.Println("End")
}
func defer_time() {
	t := time.Now()

	defer log.Println("defer1", time.Since(t)) // non-deferred call to time.Since
	tmp := time.Since(t)

	defer log.Println("defer2", tmp) // equivalent to the previous defer

	defer func() {
		log.Println("defer3", time.Since(t)) // a correctly deferred call to time.Since
	}()
}
func range_integer() {
	for i := range 10 {
		fmt.Println(10 - i)
	}
	fmt.Println("go1.22 has lift-off!")
}
func get_rid_of_share_loop_variables() {
	// go 1.22
	values := []int{1, 2, 3, 4, 5}
	for _, val := range values {
		go func() {
			fmt.Printf("%d ", val)
		}()
	}
	// Result: 1 2 3 4 5
	time.Sleep(10 * time.Second)

}
