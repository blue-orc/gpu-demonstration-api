package main

import (
	"fmt"
	"gpu-demonstration-api/device-monitor"
	"gpu-demonstration-api/python-job-runner"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) retrieveAndPushData() {
	for {
		ds := []byte("DeviceStatus")
		ds = append(ds, byte('\u0017'))
		dsBytes, err := DeviceMonitor.GetCurrentStatusJSON()
		if err != nil {
			fmt.Println("Error getting Device Status: " + err.Error())
			continue
		}
		ds = append(ds, dsBytes...)
		ds = append(ds, byte('\u00DE'))
		h.broadcast <- ds

		js := []byte("JobStatus")
		js = append(js, byte('\u0017'))
		jsBytes, err := PythonJobRunner.GetStatusJSON()
		if err != nil {
			fmt.Println("Get Python Job Status: " + err.Error())
			continue
		}
		js = append(js, jsBytes...)
		js = append(js, byte('\u00DE'))
		h.broadcast <- js

		time.Sleep(1 * time.Second)
	}
}

func (h *Hub) run() {
	go h.retrieveAndPushData()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
