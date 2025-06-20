package main

import (
	"log"
	"net"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("error listen:", err)
		return
	}
	defer l.Close()
	log.Println("listen ok")

	var i int
	for {
		log.Println("Sleeping for 10 seconds...")
		time.Sleep(time.Second * 10)
		log.Println("Woke up, attempting to accept a connection...")
		if _, err := l.Accept(); err != nil {
			log.Println("accept error:", err)
			break
		}
		i++
		log.Printf("%d: accept a new connection at %s\n", i, time.Now().Format(time.RFC3339))
	}
}
