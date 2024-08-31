//go:generate cp $GOROOT/misc/wasm/wasm_exec.js .
//go:generate env GOOS=js GOARCH=wasm go build -o jsonnet.wasm ./jsonnet

package main

import (
	"embed"
	"net"
	"net/http"
)

//go:embed wasm_exec.js jsonnet.wasm index.html
var content embed.FS

func main() {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	println("http://" + l.Addr().String())
	// panic(http.Serve(l, http.FileServer(http.FS(content))))
	panic(http.Serve(l, http.FileServer(http.Dir("."))))
}
