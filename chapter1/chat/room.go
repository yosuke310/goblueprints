package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yosuke310/goblueprints/chapter1/trace"
)

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("Client joined.")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left.")
		case msg := <-r.forward:
			r.tracer.Trace("New message : ", string(msg))
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace("Message sent.")
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("Message send failed, cleanup client.")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}
