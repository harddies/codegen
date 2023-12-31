package main

import (
	"fmt"

	"github.com/harddies/codegen/cmd"
)

func main() {
	defer recoverPanic()

	err := cmd.Execute()
	if err != nil {
		fmt.Println("cmd.Execute", err)
	}
}

func recoverPanic() {
	if err := recover(); err != nil {
		println("panic", err)
	}
}
