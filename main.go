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

type HandleHTTP struct {
	http *http.Server
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Path is %s", r.URL.Path[1:])

}

func main() {
	go DeviceMonitor.Init()
	fmt.Println("Device monitor started")

	mux1 := mux.NewRouter()
	initializeControllers(mux1)
	go func() {
		err := http.ListenAndServe(":80",
			handlers.CORS(
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
				handlers.AllowedOrigins([]string{"*"}))(mux1))
		fmt.Println("API started")
		if err != nil {
			log.Fatal("API ListenAndServe: ", err)
		}
	}()

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
		log.Fatal("WS ListenAndServe: ", err)
	}
	fmt.Println("closing")
}

func initializeControllers(r *mux.Router) {
	controllers.InitStatusController(r)
	controllers.InitPythonJobRunnerController(r)
}
