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
	url := "http://" + listener.Addr().String() + "/hello.html"
	fmt.Println("By visiting " + url + " you will get this file:")
	fmt.Println("")
	fmt.Println(static["/hello.html"].Contents)
	http.Handle("/", http.FileServer(airfreight.HTTPFileSystem(static)))
	panic(http.Serve(listener, nil))
}
