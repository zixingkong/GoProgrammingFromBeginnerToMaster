package main

import "fmt"

func main() {
	sum := 0
	for i := 0; i < 5; i++ {
		sum += i
		fmt.Println("当前值:", sum)
	}
}
