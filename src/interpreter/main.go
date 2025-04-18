package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("SmalltalkLSP Bytecode Interpreter")
	
	if len(os.Args) > 1 {
		imagePath := os.Args[1]
		fmt.Printf("Loading image from: %s\n", imagePath)
		
		// Load and execute the image
		vm := NewVM()
		if err := vm.LoadImage(imagePath); err != nil {
			fmt.Printf("Error loading image: %s\n", err)
			os.Exit(1)
		}
		
		if err := vm.Execute(); err != nil {
			fmt.Printf("Error executing image: %s\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Usage: interpreter <image-path>")
		os.Exit(1)
	}
}
