package main

import (
"log"
"github.com/sacOO7/gowebsocket"
"os"
"os/signal"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New("ws://localhost/")

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Fatal("Received connect error - ", err)
	}

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server");
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Received message - " + message)
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Received ping - " + data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Received pong - " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		return
	}

	socket.Connect()

	socket.SendText("Thoughtworks guys are awesome !!!!")

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}
}