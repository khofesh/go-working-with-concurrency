package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_printSomething(t *testing.T) {
	stdOut := os.Stdout

	read, write, _ := os.Pipe()

	os.Stdout = write

	var wg sync.WaitGroup
	wg.Add(1)

	go printSomething("something", &wg)

	wg.Wait()

	_ = write.Close()

	result, _ := io.ReadAll(read)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "something") {
		t.Errorf("Expected to find something, but it's not there")
	}
}
