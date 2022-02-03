package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const prompt = "and press ENTER when ready."

func main() {
	// random seed
	rand.Seed(time.Now().UnixNano())

	var firstNumber = rand.Intn(8) + 2
	var secondNumber = rand.Intn(8) + 2
	var subtraction = rand.Intn(8) + 2
	var answer = firstNumber*secondNumber - subtraction

	game(firstNumber, secondNumber, subtraction, answer)
}

func game(firstNumber, secondNumber, subtraction, answer int) {
	reader := bufio.NewReader(os.Stdin)

	// display a welcome/instructions
	fmt.Println("Guess the Number Game")
	fmt.Println("_____________________")
	fmt.Println("")

	fmt.Println("Think of a number between 1 and 10", prompt)
	reader.ReadString('\n')

	//take them through the games
	fmt.Println("Multiply your number by", firstNumber, prompt)
	reader.ReadString('\n')

	fmt.Println("Now multiply the result by", secondNumber, prompt)
	reader.ReadString('\n')

	fmt.Println("Divide the result by the number you originally thought of", prompt)
	reader.ReadString('\n')

	fmt.Println("Now subtract", subtraction, prompt)
	reader.ReadString('\n')

	//give them the answer
	fmt.Println("The answer is", answer)
}
