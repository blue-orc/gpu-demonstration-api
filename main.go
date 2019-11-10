package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gpu-demonstration-api/controllers"
	"gpu-demonstration-api/database"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	database.Initialize()
	go DeviceMonitor.Init()
	fmt.Println("Device monitor started")

	mux1 := mux.NewRouter()
	initializeControllers(mux1)
	go func() {
		fmt.Println("API starting")
		err := http.ListenAndServe(":80",
			handlers.CORS(
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
				handlers.AllowedOrigins([]string{"*"}))(mux1))
		if err != nil {
			log.Fatal("API ListenAndServe: ", err)
		}
	}()

	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	fmt.Println("Starting websocket server")
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("WS ListenAndServe: ", err)
	}
	fmt.Println("closing")
}

func initializeControllers(r *mux.Router) {
	controllers.InitStatusController(r)
	controllers.InitPythonJobRunnerController(r)
}
