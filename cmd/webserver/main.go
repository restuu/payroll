package main

import "log"

func main() {
	ws, err := NewWebServer()
	if err != nil {
		log.Fatal(err)
	}
	ws.Start()
}
