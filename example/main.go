package main

//go:generate go run staticgen/main.go

import (
	"fmt"
	"net"
	"net/http"

	"github.com/lassik/airfreight"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	fmt.Println("http://" + listener.Addr().String() + "/hello.html")
	http.Handle("/", http.FileServer(airfreight.MapFileSystem(Static)))
	panic(http.Serve(listener, nil))
}
