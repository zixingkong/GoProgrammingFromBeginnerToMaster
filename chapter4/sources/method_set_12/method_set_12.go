package main

import method_set_utils 		"github.com/method_set_utils"

type T struct{}

func (T) M1()  {}
func (*T) M2() {}

type Interface interface {
	M1()
	M2()
}

type T1 = T
type Interface1 = Interface

func main() {
	var t T
	var pt *T
	var t1 T1
	var pt1 *T1

	method_set_utils.DumpMethodSet(&t)
	method_set_utils.DumpMethodSet(&t1)

	method_set_utils.DumpMethodSet(&pt)
	method_set_utils.DumpMethodSet(&pt1)

	method_set_utils.DumpMethodSet((*Interface)(nil))
	method_set_utils.DumpMethodSet((*Interface1)(nil))
}
