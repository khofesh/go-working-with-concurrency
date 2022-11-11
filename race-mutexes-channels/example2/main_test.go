package main

import "testing"

func Test_updateMessage(t *testing.T) {
	msg = "hola mundo!"

	wg.Add(2)
	go updateMessage("x")
	go updateMessage("goodbye, cruel world!")
	wg.Wait()

	if msg != "goodbye, cruel world!" {
		t.Error("incorrect value in msg")
	}
}
