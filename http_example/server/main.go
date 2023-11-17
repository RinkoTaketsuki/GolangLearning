package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

// sorted by initialization order
var (
	ListeningAddress = flag.String("addr", "0.0.0.0:17180", "http service address")
	Middlewares      = []Middleware{
		WithPanicRecovery,
		WithLogger,
		WithMetrics,
	}
	Logger       *log.Logger
	FullServeMux *ServeMux
	Server       *http.Server
)

func init() {
	flag.Parse()
	logFile, err := os.OpenFile("./server.log", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("OpenFile:", err)
	}
	Logger = log.New(logFile, "", log.LstdFlags)
	FullServeMux = NewServeMux()
	FullServeMux.ApplyMiddlewares(Middlewares...)
	Server = &http.Server{
		Addr:           *ListeningAddress,
		Handler:        FullServeMux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes << 1,
		ErrorLog:       Logger,
		// Disable HTTP/2
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
}

func main() {
	FullServeMux.HandleFunc("/", HandlerQR)
	Logger.Fatal("ListenAndServe:", Server.ListenAndServe())
}
