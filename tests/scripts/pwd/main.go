package main

import (
	"fmt"
	"os"
)

func main() {
	path, _ := os.Getwd()
	fmt.Println("Current path:", path)
	err := os.WriteFile("pwd.txt", []byte(path), 0600)
	if err != nil {
		fmt.Println("Could not write to pwd.txt:", err)
		return
	}
	fmt.Println("Wrote to pwd.txt")
}
