package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"

	"github.com/Alienero/goproxy/goproxy"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	psw := flag.String("password", "", "")
	addr := flag.String("listen", ":808", "-listen 127.0.0.1:8080")
	flag.Parse()
	if *psw == "" {
		panic("Password is nil")
	}
	err := http.ListenAndServeTLS(*addr, "cert.pem", "key.pem", goproxy.NewServer(*psw))
	if err != nil {
		log.Fatal(err)
	}
}
