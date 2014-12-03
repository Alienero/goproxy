package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/Alienero/goproxy/goproxy"
)

var (
	psw    string
	addr   string
	remote string

	f_psw     = flag.String("password", "", "-password xxxx")
	f_addr    = flag.String("listen", "", "-listen")
	f_remote  = flag.String("remote", "", "-remote")
	f_isColor = flag.Bool("iscolor", false, "-iscolor true")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	var (
		bio  = bufio.NewReader(os.Stdin)
		line []byte
		err  error
	)

	if runtime.GOOS == "windows" && *f_isColor {
		log.Println("Open the color")
		goproxy.SetColor()
	}

	if *f_psw != "" {
		log.Println("Use the Arg[password]")
		psw = *f_psw
	} else {
		psw = GetPass()
	}

	if *f_addr != "" {
		log.Println("Use the Arg[listen]")
		addr = *f_addr
	} else {
		print("Please input the listen address<entenr to use default>>")
		line, _, err = bio.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		addr = string(line)
	}

	if *f_remote != "" {
		log.Println("Use the Arg[remote]")
		remote = *f_remote
	} else {
		print("Please input the remote address<entenr to use default>>")
		line, _, err = bio.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		remote = string(line)
	}

	// Set default
	if addr == "" {
		addr = "127.0.0.1:808"
	}
	if remote == "" {
		remote = "yim.so:808"
	}
	log.Println("Listen in", addr)
	log.Println("Rmote address is", remote)

	c := goproxy.NewClient(remote, psw)
	go c.PrintTraffic()
	err = http.ListenAndServe(addr, c)
	if err != nil {
		log.Fatal(err)
	}
}
