package bot

import (
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"net/http"
	"strconv"
	"fmt"
	"time"
)

type Message map[string]interface{}
type Handler func(msg Message)

type HandlerSet map[string]Handler

type Bot struct {
	Cli *gosocketio.Client
	Headers http.Header
	Hs HandlerSet
	Id, Dialog_id float64
}

func NewBot() (*Bot) {
	bot := new(Bot)
	bot.Headers = make(http.Header)
	handlerSet := make(HandlerSet)
	// Set wanted Headers	
	bot.Headers.Add("Accept-Encoding", "gzip, deflate")
	bot.Headers.Add("Accept-Language", "en-US,en;q=0.9")
	bot.Headers.Add("Cache-Control", "no-cache")
	bot.Headers.Add("Pragma", "no-cache")
	bot.Headers.Add("Origin", "http://nekto.me")
	bot.Headers.Add("Host", "im.nekto.me")
	bot.Headers.Add("User-Agent", "Mozilla/5.0 (X11; Darwin x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Mojave Chromium/73.0.3683.86 Chrome/73.0.3683.86 Safari/537.36")
	// Init default Headers to headerSet
	handlerSet["auth.successToken"] = func (msg Message) {
		fmt.Println("success token, great job!")
		mg := msg["data"].(map[string]interface{})
		bot.Id = mg["id"].(float64)
	}
	handlerSet["error.code"] = func (msg Message) {
		fmt.Println("error, look at this:", msg)
	}
	handlerSet["captcha.verify"] = func (msg Message) {
		fmt.Println("captcha", msg)
	}
	handlerSet["search.success"] = func (msg Message) {
		fmt.Println("searching for stranger...")
	}
	handlerSet["dialog.opened"] = func (msg Message) {
		mg := msg["data"].(map[string]interface{})
		bot.Dialog_id = mg["id"].(float64)
		fmt.Println("dialog opened")
	}
	//handlerSet["dialog.typing"] = func (msg Message) {
	//	fmt.Println("stranger typing")
	//}
	//handlerSet["messages.reads"] = func (msg Message) {
	//	fmt.Println("stranger reads")
	//}
	handlerSet["messages.new"] = func (msg Message) {
		mg := msg["data"].(map[string]interface{})
		if mg["senderId"] != bot.Id {
			fmt.Println(mg["message"].(string))
		}
	}
	handlerSet["dialog.closed"] = func (msg Message) {
		fmt.Println("dialog closed")
	}
	handlerSet["unexpected"] = func (msg Message) {
		fmt.Println("bot doesn't know about this messge type, you should add it to bot.hs", msg)
	}
	bot.Hs = handlerSet
	return bot
}

func (bot *Bot) SetCookie(cookie string) {
	bot.Headers.Set("Cookie", cookie)
}

// Connect to the server and pass auth, token might be just empty string
func (bot *Bot) Connect(token string) error {
	tptr := &transport.WebsocketTransport{
		PingInterval: 20 * time.Second,
		PingTimeout: 10 * time.Second,
		ReceiveTimeout: transport.WsDefaultReceiveTimeout,
		SendTimeout: transport.WsDefaultSendTimeout,
		BufferSize: transport.WsDefaultBufferSize,
		RequestHeader: bot.Headers,
	}
	c, err := gosocketio.Dial(gosocketio.GetUrl("im.nekto.me", 80, false), tptr)
	if err != nil {
		return err
	}

	err = c.On("notice", func(h *gosocketio.Channel, args Message) {
		if bot.Hs[args["notice"].(string)] != nil {
			bot.Hs[args["notice"].(string)](args)
		}
	})
	if err != nil {
		return err
	}

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		fmt.Println("Disconnected from the server")
	})
	if err != nil {
		return err
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		fmt.Println("Connected to the server")
		if token == "" {
			c.Emit("action", map[string]interface{}{"action":"auth.getToken","deviceType": 2})
		} else {
			c.Emit("action", map[string]interface{}{"action":"auth.sendToken","token":token})
		}
	})
	if err != nil {
		return err
	}
	bot.Cli = c
	return nil
}

func (bot *Bot) SendMessage(msg string) {
		t := time.Now()
		bot.Cli.Emit("action", map[string]interface{}{"action":"anon.message","dialogId":bot.Dialog_id,
		"message":msg,"randomId":strconv.FormatFloat(bot.Id, 'f', -1, 64)+"_"+strconv.FormatInt(t.Unix(), 10)})
}

func (bot *Bot) LeaveDialog() {
	bot.Cli.Emit("action", map[string]interface{}{"action":"anon.leaveDialog", "dialogId": bot.Dialog_id})
}

func (bot *Bot) StartSearch() {
	bot.Cli.Emit("action", map[string]interface{}{"action":"search.run","myAge":[]int{18, 21},"mySex":"F","wishAge":*new([]int),"wishSex":"M"})
}

func (bot *Bot) Track() {
	bot.Cli.Emit("action", map[string]interface{}{"action":"online.track","on":true})
}

func (bot *Bot) UnTrack() {
	bot.Cli.Emit("action", map[string]interface{}{"action":"online.track","on":false})
}

func (bot *Bot) Close() {
	bot.Cli.Close()
}

