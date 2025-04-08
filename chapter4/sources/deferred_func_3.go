package main

import "fmt"

func bar() {
	fmt.Println("raise a panic")
	panic(-1)
}

func foo() {
	bar()
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("recovered from a panic")
		}
	}()
}

func main() {
	foo()
	fmt.Println("main exit normally")
}
