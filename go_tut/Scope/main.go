package main

import "fmt"

func main() {
	var one = "One"
	fmt.Println(one)
	newFunc()
}

func newFunc() {
	var one = "The number one!"
	fmt.Println(one)
}
