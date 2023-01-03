package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jbert/gof/interpreter"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}

func run() error {
	i := interpreter.New()
	err := i.Run("1 2 3 dup . 4 + DUP .")
	if err != nil {
		return fmt.Errorf("Failed to run: %s\n", err)
	}
	i.DumpStack(os.Stdout)
	return nil
}
