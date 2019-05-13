package main

import (
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"net/http"
	"bufio"
	"os"
	"fmt"
	"strings"
	"time"
)

type Message map[string]interface{}

func main() {
	customHeaders := make(http.Header)
	customHeaders.Add("Accept-Encoding", "gzip, deflate")
	customHeaders.Add("Accept-Language", "en-US,en;q=0.9")
	customHeaders.Add("Cache-Control", "no-cache")
	customHeaders.Add("Pragma", "no-cache")
	customHeaders.Add("Origin", "http://nekto.me")
	customHeaders.Add("Host", "im.nekto.me")
	customHeaders.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/73.0.3683.86 Chrome/73.0.3683.86 Safari/537.36")
	//set cookie header for trusted connection 

	tptr := &transport.WebsocketTransport{
		PingInterval: 20 * time.Second,
		PingTimeout: 10 *time.Second,
		ReceiveTimeout: transport.WsDefaultReceiveTimeout,
		SendTimeout: transport.WsDefaultSendTimeout,
		BufferSize: transport.WsDefaultBufferSize,
		RequestHeader: customHeaders,
	}
	c, err := gosocketio.Dial(gosocketio.GetUrl("im.nekto.me", 80, false), tptr)
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("notice", func(h *gosocketio.Channel, args Message) {
		log.Println("--- Got chat message: ", args)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected c")
		c.Emit("action", map[string]interface{}{"action": "auth.getToken", "deviceType":2})
	})
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)

	var args []string
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		args = strings.Split(text, " ")
		fmt.Println("args:",string(args[0]))
		// if you input "track true", the server will send you data about online users in chat	
		if string(args[0]) == "track" {
			c.Emit("action", map[string]string{"action":"online.track","on":string(args[1])})
		}
	}
}
