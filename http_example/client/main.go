package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	ConnAddress = flag.String("addr", "0.0.0.0:17180", "http service address")
	// Clients and Transports are safe for concurrent use by multiple goroutines
	// and for efficiency should only be created once and re-used.
	Client = &http.Client{
		Timeout: 4 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       http.DefaultMaxIdleConnsPerHost << 1,
			IdleConnTimeout:    time.Minute,
			DisableCompression: true,
			// Disable HTTP/2
			TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
		},
	}
)

func main() {
	req, err := http.NewRequest(http.MethodGet, *ConnAddress+"/", nil)
	if err != nil {
		panic(err)
	}
	//
	// customize your request here, such as its header
	//
	resp, err := Client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	sb := strings.Builder{}
	resp.Write(&sb)
	fmt.Print(sb.String())
}
