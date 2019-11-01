package main

import (
	"flag"
	"gpu-demonstration-api/device-monitor"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":9002", "http service address")
var addrTLS = flag.String("addrTLS", ":9003", "https service address")

func main() {
	DeviceMonitor.Init()
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
