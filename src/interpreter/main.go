package main

import (
	"fmt"
	"os"

	"smalltalklsp/interpreter/demo"
	"smalltalklsp/interpreter/vm"
)

func main() {
	fmt.Println("SmalltalkLSP Bytecode Interpreter")

	if len(os.Args) > 1 {
		if os.Args[1] == "demo" {
			// Run the factorial demo
			demo.RunFactorialDemo()
		} else {
			imagePath := os.Args[1]
			fmt.Printf("Loading image from: %s\n", imagePath)

			// Load and execute the image
			virtualMachine := vm.NewVM()
			if err := virtualMachine.LoadImage(imagePath); err != nil {
				fmt.Printf("Error loading image: %s\n", err)
				os.Exit(1)
			}

			result, err := virtualMachine.Execute()
			if err != nil {
				fmt.Printf("Error executing image: %s\n", err)
				os.Exit(1)
			}

			fmt.Printf("Final result: %s\n", result)
		}
	} else {
		fmt.Println("Usage: interpreter <image-path|demo>")
		os.Exit(1)
	}
}
