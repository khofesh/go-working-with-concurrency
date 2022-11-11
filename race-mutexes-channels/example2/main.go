package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(s string, m *sync.Mutex) {
	defer wg.Done()

	m.Lock()
	msg = s

	m.Unlock()
}

func main() {
	msg = "hola mundo!"

	var mutex sync.Mutex

	wg.Add(2)
	go updateMessage("hola universo!", &mutex)
	go updateMessage("hola cosmos!", &mutex)
	wg.Wait()

	fmt.Println(msg)
}
