package main

import (
	"log"

	"./demo/printer"
)

func main() {
	proxy, err := printer.NewProxy("./config.client", "SimplePrinter:default -p 10000", false)
	if err != nil {
		log.Printf("init|%v", err)
		return
	}

	if proxy == nil {
		log.Printf("proxy is nil")
		return
	}

	defer proxy.Close()

	log.Printf("proxy is %v", proxy)

	data, err := proxy.Echo([]byte("hello world"))

	if err != nil {
		log.Printf("DO|%v", err)
	}

	log.Printf("%s", data)
}
