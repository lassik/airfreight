package main

import "github.com/lassik/airfreight/packer"

func main() {
	packer.Package("main").Map("static", "static").WriteFile("static.go")
}
