package main
import (
	"./bot"
	"fmt"
	"strings"
	"bufio"
	"os"
)

func main() {
	b := bot.NewBot()
	
	b.Hs["auth.successToken"] = func(msg bot.Message) {
		fmt.Println("b - success token, great job!")
		mg := msg["data"].(map[string]interface{})
		b.Id = mg["id"].(float64)
		b.Dialog_id = 0
		b.StartSearch()
	}
	b.Hs["dialog.opened"] = func(msg bot.Message) {
		mg := msg["data"].(map[string]interface{})
		b.Dialog_id = mg["id"].(float64)
		fmt.Println("b - dialog opened")
	}
	b.Hs["dialog.closed"] = func (msg bot.Message) {
		fmt.Println("b - dialog closed")
		b.Dialog_id = 0
		b.StartSearch()
	}
	
	go b.Connect(os.Args[1])

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		args := strings.Split(text, " ")
		switch string(args[0]) {
			case "track":
				b.Track()
			case "untrack":
				b.UnTrack()
			case "search":
				b.StartSearch()
			case "leave":
				b.LeaveDialog()
			case "close":
				b.Close()
			case "connect":		
				b.Connect(os.Args[1])	
				break
			default:
				wo := ""
				for _, word := range args {
					wo += word+" "
				}
				b.SendMessage(wo)
		}
	}
}
