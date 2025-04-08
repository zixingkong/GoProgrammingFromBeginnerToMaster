package main

import (
	"io"
	method_set_utils "github.com/method_set_utils"
)

func main() {
	method_set_utils.DumpMethodSet((*io.Writer)(nil))
	method_set_utils.DumpMethodSet((*io.Reader)(nil))
	method_set_utils.DumpMethodSet((*io.Closer)(nil))
	method_set_utils.DumpMethodSet((*io.ReadWriter)(nil))
	method_set_utils.DumpMethodSet((*io.ReadWriteCloser)(nil))
}
