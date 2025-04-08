package main

import method_set_utils "github.com/method_set_utils"

type Interface1 interface {
	M1()
}

type Interface2 interface {
	M1()
	M2()
}

type Interface3 interface {
	Interface1
	Interface2
}

type Interface4 interface {
	Interface2
	M2()
}

func main() {
	method_set_utils.DumpMethodSet((*Interface3)(nil))
}
