package main

import (
	"./chat"
	"net/http"
)


func main() {
	server := chat.NewServer()
	go server.Start()

	http.HandleFunc("/chat", server.ServeWebSocket)
	http.Handle("/", http.FileServer(http.Dir("./")))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
