package main

import (
	"fmt"
	"time"
)

func listenToChan(ch chan int) {
	for {
		// print the data you got
		i := <-ch
		fmt.Println("got ", i, "from channel")

		// busy
		time.Sleep(1 * time.Second)
	}
}

func main() {
	ch := make(chan int, 10)

	go listenToChan(ch)

	for i := 0; i <= 100; i++ {
		fmt.Println("sending ", i, "to channel...")
		ch <- i

		fmt.Println("sent ", i, "to channel")
	}

	fmt.Println("done")
	close(ch)
}
