package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func Test_updateMessage(t *testing.T) {
	wg.Add(1)

	go updateMessage("epsilon", &wg)

	wg.Wait()

	if msg != "epsilon" {
		t.Error("expected to find epsilon, but it's not there")
	}
}

func Test_printMessage(t *testing.T) {
	stdOut := os.Stdout

	read, write, _ := os.Pipe()
	os.Stdout = write

	msg = "epsilon"
	printMessage()

	_ = write.Close()

	result, _ := io.ReadAll(read)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "epsilon") {
		t.Error("expected to find epsilon, but it's not there")
	}
}

func Test_main(t *testing.T) {
	stdOut := os.Stdout

	read, write, _ := os.Pipe()
	os.Stdout = write

	main()

	_ = write.Close()

	result, _ := io.ReadAll(read)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "Hello, universe!") {
		t.Error("expected to find Hello, universe!, but it's not there")
	}

	if !strings.Contains(output, "Hello, cosmos!") {
		t.Error("expected to find Hello, cosmos!, but it's not there")
	}

	if !strings.Contains(output, "Hello, world!") {
		t.Error("expected to find Hello, world!, but it's not there")
	}
}
