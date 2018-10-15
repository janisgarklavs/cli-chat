package main

import "fmt"

type Server struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	connect    chan *Client
	disconnect chan *Client
}

func (s *Server) listen() {
	for {
		select {
		case client := <-s.connect:
			fmt.Println("connected")
			s.clients[client] = true
		case client := <-s.disconnect:
			fmt.Println("disconnected")
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
		case message := <-s.broadcast:
			fmt.Println("message")
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
		}
	}
}

func newServer() *Server {
	return &Server{
		broadcast:  make(chan []byte),
		connect:    make(chan *Client),
		disconnect: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}
