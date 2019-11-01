package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gpu-demonstration-api/controllers"
	"gpu-demonstration-api/device-monitor"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":9002", "http service address")
var addrTLS = flag.String("addrTLS", ":9003", "https service address")

func main() {
	go DeviceMonitor.Init()
	fmt.Println("Device monitor started")

	r := mux.NewRouter()
	initializeControllers(r)
	go http.ListenAndServe(":80",
		handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(r))

	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Hit")
		serveWs(hub, w, r)
	})

	fmt.Println("Starting websocket server")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initializeControllers(r *mux.Router) {
	controllers.InitStatusController(r)
	controllers.InitPythonJobRunnerController(r)
}
