package main

import (
	"fmt"
	"myapp/packageone"
)

// Package level var
var one = "One"

func main() {
	// Block level var
	var one = "main One"
	fmt.Println(one)
	newFunc()

	var newString = packageone.PublicVar
	fmt.Println(newString)

	packageone.Exported()
}

func newFunc() {
	// var one = "The number one!"
	fmt.Println(one)
}
