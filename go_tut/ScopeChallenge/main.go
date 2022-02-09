package main

import (
	"myapp/packageone"
)

var myVar = "myVar"

func main() {
	var blockVar = "blockVar"
	var PackageVar = packageone.PackageVar

	packageone.PrintMe(myVar, blockVar, PackageVar)
}
