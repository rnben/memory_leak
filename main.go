package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/armon/go-socks5"
)

func main() {
	start8320And8314()

	startSocks5(":1080", "foo", "bar")

	exitChan := make(chan os.Signal, 2)
	signal.Notify(exitChan, os.Interrupt, os.Kill)

	go work()

	select {
	case v := <-exitChan:
		log.Println(v)
	}
}

func startSocks5(listenAddr, user, password string) {
	conf := &socks5.Config{
		AuthMethods: []socks5.Authenticator{
			socks5.UserPassAuthenticator{
				Credentials: socks5.StaticCredentials{user: password},
			},
		},
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	go func() {
		log.Fatalln(server.ListenAndServe("tcp", listenAddr))
	}()
}

var isPanic int64

func start8320And8314() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path, "FROM", r.RemoteAddr)
		fmt.Fprintf(w, "pong\n")
	})

	http.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&isPanic, 1)
	})

	go func() {
		log.Fatalln(http.ListenAndServe(":8314", nil))
	}()

	go func() {
		log.Fatalln(http.ListenAndServeTLS(":8320", "cert.pem", "key.pem", nil))
	}()
}

func work() {
	var name = "a"
	for {
		if isPanic > 0 {
			panic("this is a panic")
		}
		// name += name
		println(len(name))
		time.Sleep(time.Second * 3)
	}
}
