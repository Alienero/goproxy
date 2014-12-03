package goproxy

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

var hasPort = regexp.MustCompile(`:\d+$`)

type ConnectActionLiteral int

const (
	ConnectAccept = iota
	ConnectReject
	ConnectMitm
)

var (
	OkConnect = &ConnectAction{Action: ConnectAccept}
)

type ConnectAction struct {
	Action ConnectActionLiteral
}

func handleHttps(w http.ResponseWriter, r *http.Request) {

	hij, ok := w.(http.Hijacker)
	if !ok {
		log.Println("httpserver does not support hijacking")
	}

	proxyClient, _, e := hij.Hijack()
	if e != nil {
		log.Println("Cannot hijack connection " + e.Error())
		return
	}

	todo, host := OkConnect, r.URL.Host

	switch todo.Action {
	case ConnectAccept:
		if !hasPort.MatchString(host) {
			host += ":80"
		}
		https_proxy := os.Getenv("https_proxy")
		if https_proxy == "" {
			https_proxy = os.Getenv("HTTPS_PROXY")
		}
		var targetSiteCon net.Conn
		var e error
		if https_proxy != "" {
			targetSiteCon, e = net.Dial("tcp", https_proxy)
		} else {
			targetSiteCon, e = net.Dial("tcp", host)
		}
		if e != nil {
			// trying to mimic the behaviour of the offending website
			// don't answer at all
			log.Println(e)
			return
		}
		if https_proxy != "" {
			connectReq := &http.Request{
				Method: "CONNECT",
				URL:    &url.URL{Opaque: host},
				Host:   host,
				Header: make(http.Header),
			}
			e = connectReq.Write(targetSiteCon)
			if e != nil {
				targetSiteCon.Close()
				log.Println(e)
				return
			}

			// Read response.
			// Okay to use and discard buffered reader here, because
			// TLS server will not speak until spoken to.
			br := bufio.NewReader(targetSiteCon)
			resp, err := http.ReadResponse(br, connectReq)
			if err != nil {
				targetSiteCon.Close()
				w.WriteHeader(500)
				return
			}
			if resp.StatusCode != 200 {
				targetSiteCon.Close()
				w.WriteHeader(resp.StatusCode)
				io.Copy(w, resp.Body)
				resp.Body.Close()
				return
			}
		}

		proxyClient.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		go copyAndClose(targetSiteCon, proxyClient)
		copyAndClose(proxyClient, targetSiteCon)

	default:
		proxyClient.Close()
	}
}
func copyAndClose(w io.WriteCloser, r io.Reader) {
	io.Copy(w, r)
	if err := w.Close(); err != nil {
		log.Println("Error closing", err)
	}
}
