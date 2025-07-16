package main

import (
	"log"
	"time"
)

func init() {
	time.Local = time.FixedZone("Asia/Jakarta", 7*60*60)
}

func main() {
	ws, err := NewWebServer()
	if err != nil {
		log.Fatal(err)
	}
	ws.Start()
}
