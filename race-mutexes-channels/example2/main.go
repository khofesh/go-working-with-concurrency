package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(s string) {
	defer wg.Done()
	msg = s
}

func main() {
	msg = "hola mundo!"

	wg.Add(2)
	go updateMessage("hola universo!")
	go updateMessage("hola cosmos!")
	wg.Wait()

	fmt.Println(msg)
}
