package main

import (
	"flag"
	"github.com/pr0head/ws_example_client/ws"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var (
	addr = flag.String("addr", "localhost:8080", "http service address")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()

	writeWait := 10 * time.Second
	pongWait := 60 * time.Second

	wbs := ws.NewWebSocket(c, pongWait, writeWait)
	go func() {
		sendByTicker(time.NewTicker(3*time.Second), wbs.SendSetServerStatus)
	}()
	go func() {
		sendByTicker(time.NewTicker(5*time.Second), wbs.SendAddGameChar)
	}()
	go func() {
		sendByTicker(time.NewTicker(7*time.Second), wbs.SendSendGameBalance)
	}()

	go wbs.Listen()

	<-interrupt
	log.Println("interrupt")

	return
}

func sendByTicker(t *time.Ticker, f func() error) {
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if err := f(); err != nil {
				log.Print("send by ticker err", err)
				return
			}
		}
	}
}
