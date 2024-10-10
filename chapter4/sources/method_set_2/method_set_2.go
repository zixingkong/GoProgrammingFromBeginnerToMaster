package main

import method_set_utils "github.com/method_set_utils"

type Interface interface {
	M1()
	M2()
}

type T struct{}

func (t T) M1()  {}
func (t *T) M2() {}

func main() {
	var t T
	var pt *T
	method_set_utils.DumpMethodSet(&t)
	method_set_utils.DumpMethodSet(&pt)
	method_set_utils.DumpMethodSet((*Interface)(nil))
}
