package main

import "github.com/lassik/airfreight/packer"

func main() {
	packer.Package("main").Map("Static", "static").WriteFile("static.go")
}
