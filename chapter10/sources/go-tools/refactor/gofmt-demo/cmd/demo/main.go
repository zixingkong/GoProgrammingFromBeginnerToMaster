package main

import (
	"github.com/bigwhite/gofmt-demo/pkg/bar"
	"github.com/bigwhite/gofmt-demo/pkg/foo"
	"github.com/rs/zerolog/log"
)

func main() {
	foo.Foo()
	bar.Bar()
	log.Print("gofmt-demo")
}
