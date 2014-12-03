package goproxy

import (
	"bufio"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	// "net/http/httputil"
	// "net/url"
	"sync"
	"sync/atomic"
	"time"
)

func copyHeaders(dst, src http.Header) {
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

type Client struct {
	remote_addr string
	psw         string

	upload   uint64
	download uint64

	errCounter uint64
	lock       *sync.RWMutex
}

func NewClient(remote, psw string) *Client {
	if remote == "" {
		panic("rempte addr is nil")
	}
	if psw == "" {
		panic("password is nil")
	}
	return &Client{
		remote_addr: remote,
		psw:         psw,
		lock:        new(sync.RWMutex),
	}
}

// var remoteAddr = "jcode.name:808"
// var psw = ""

func (client *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Proxy-Authenticate", client.psw)
	c, err := net.Dial("tcp", client.remote_addr)
	if err != nil {
		//should not print
		log.Println(GetColorError(err))
		return
	}
	defer c.Close()
	if tc, ok := c.(*net.TCPConn); ok {
		if tc.SetReadBuffer(4096*32) != nil {
			log.Println(GetColorError(err))
			return
		}
		if tc.SetKeepAlive(true) != nil {
			log.Println(GetColorError(err))
			return
		}
		if tc.SetKeepAlivePeriod(1*time.Minute) != nil {
			log.Println(GetColorError(err))
			return
		}
	}
	conn := NewRecord(tls.Client(c, &tls.Config{InsecureSkipVerify: true}), client.lock, &client.upload, &client.download)
	//Send the requset to own proxy server
	read := bufio.NewReader(conn)
	if err = r.WriteProxy(conn); err != nil {
		log.Println(GetColorError(err))
		atomic.AddUint64(&client.errCounter, 1)
		return
	}
	if r.Method == "CONNECT" {
		//send the https request
		resp, err := http.ReadResponse(read, r)
		if err != nil {
			log.Println(GetColorError(err))
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return
		}
		hij, ok := w.(http.Hijacker)
		if !ok {
			log.Println("httpserver does not support hijacking")
			return
		}
		proxyClient, _, e := hij.Hijack()
		if e != nil {
			log.Println("Cannot hijack connection " + e.Error())
			return
		}
		//write the 200 ok to the client
		proxyClient.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		//connect keep the connect alive
		go copyAndClose(proxyClient, read)
		copyAndClose(conn, proxyClient)
		return
	}
	//read the response
	resp, err := http.ReadResponse(read, r)
	if err != nil {
		log.Println(GetColorError(err))
		atomic.AddUint64(&client.errCounter, 1)
		return
	}
	defer resp.Body.Close()
	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

type Load struct {
	Download uint64
	Upload   uint64
}

// TODO
func (c *Client) GetTrafficChan() chan *Load {
	ch := make(chan *Load)
	return ch
}

// Print the traffic load per second.
func (c *Client) PrintTraffic() {
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			c.lock.Lock()
			// Print
			up := c.upload
			down := c.download
			// Reset
			c.upload = 0
			c.download = 0
			c.lock.Unlock()
			if up == 0 && down == 0 {
				continue
			}
			s := Getload(up, down)
			log.Println(s)
		}
	}
}

// Record the traffic
type Record struct {
	c net.Conn

	lock *sync.RWMutex

	upload   *uint64
	download *uint64
}

func NewRecord(c net.Conn, lock *sync.RWMutex, upload *uint64, download *uint64) *Record {
	return &Record{
		c:        c,
		lock:     lock,
		upload:   upload,
		download: download,
	}
}

func (record *Record) Write(p []byte) (n int, err error) {
	n, err = record.c.Write(p)
	record.lock.RLock()
	atomic.AddUint64(record.upload, uint64(n))
	record.lock.RUnlock()
	return
}

func (record *Record) Read(p []byte) (n int, err error) {
	n, err = record.c.Read(p)
	record.lock.RLock()
	atomic.AddUint64(record.download, uint64(n))
	record.lock.RUnlock()
	return
}

func (record *Record) Close() error {
	return record.c.Close()
}
