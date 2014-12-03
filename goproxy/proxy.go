package goproxy

import (
	"io"
	"net/http"
	"strings"
	"time"
)

var transport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	ResponseHeaderTimeout: 30 * time.Second,
}

type Server struct {
	psw string
}

func NewServer(psw string) *Server {
	return &Server{
		psw: psw,
	}
}

//a normally server proxy
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//TODO you can do something what limit the user's request
	seis := strings.Split(r.Header.Get("Proxy-Authenticate"), ";")
	if seis[len(seis)-1] != s.psw {
		http.Error(w, "error", 407)
		return
	}
	r.Header.Del("Proxy-Authenticate")
	if r.Method == "CONNECT" {
		handleHttps(w, r)
		return
	}
	r.RequestURI = ""
	//TODO remove the unuse head
	resp, err := transport.RoundTrip(r)
	if err != nil {
		//TODO write the error to the log file
		println(err.Error())
		return
	}
	defer resp.Body.Close()
	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
