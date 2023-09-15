package main

import "codegen/cmd"

func main() {
	defer recoverPanic()

	cmd.Execute()
}

func recoverPanic() {
	if err := recover(); err != nil {
		println("panic", err)
	}
}
