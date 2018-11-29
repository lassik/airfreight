# Airfreight

## Utterly trivial static file embedding via `go generate`

[![Go Report Card](https://goreportcard.com/badge/github.com/lassik/airfreight)](https://goreportcard.com/report/github.com/lassik/airfreight)

There's a plethora of libraries for embedding arbitrary files inside a
Go executable. I looked at many of them, tried out two popular
choices, and was repeatedly surprised by magic behavior such as failed
builds and missing assets. Even after using those libraries fairly
extensively I still don't understand how their builds really work and
why they fail seemingly at random.

So I wrote Airfreight. This bikeshed is just the right tint of blue.
Like the dodgy pun in the name says, it's about the most lightweight
thing that could possibly work. This one is for people who are tired
of magic and complexity in the name of convenience. There are no
configuration options. Everything is bog standard and fully
composable.

## Features

- Put your static files in any directories you want.

- Describe them in a trivial [`staticgen/main.go`](example/staticgen/main.go) file.

- Put `//go:generate go run staticgen/main.go` in your real `main.go`.

- Run `go generate && go build` to build your project. No magic build
  steps. No magic options to make different kinds of builds.

- Generated code goes into any `.go` file and Go package of your
  choice. It contains one or more maps of static files. Those are
  ordinary Go maps. File contents are just strings.

- Generated code conforms to `gofmt` format.

- Generated `.go` file can easily be read to see exactly what it does
  (hardly anything!)

- Generated `.go` file can be committed to version control or not,
  according to taste. Both choices work fine.

- It's easy to generate multiple `.go` files if you need to.

- File modification times are included for caching and the like.

- Use the optional
  [`http.FileSystem`](https://golang.org/pkg/net/http/#FileSystem)
  shim for web server support.

- No live-reload "developer mode" that you may enable by accident when
  doing production builds. If you want to live-reload your assets from
  disk during development, you need to add your own
  [`http.FileServer()`](https://golang.org/pkg/net/http/#FileServer)
  as an alternative to Airfreight's shim.

## Continuous integration

### Travis

The default Travis Go build doesn't run `go generate`, tries to run
`go get` in the `install` step, and doesn't assume the new Go module
system. Put this in your `.travis.yml` to fix it:

    language: go
    go:
      - "1.11"
    env:
      - GO111MODULE=on
    install: []
    script:
      - go generate
      - go build
      - go test
