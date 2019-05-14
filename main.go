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
	b1 := bot.NewBot()

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
	b.Hs["messages.new"] = func (msg bot.Message) {
		mg := msg["data"].(map[string]interface{})
		if mg["senderId"] != b.Id {
			fmt.Println("b -",mg["message"].(string))
			if b1.Dialog_id != 0 && b2.Dialog_id != 0 && b.Dialog_id != 0 {
				b1.SendMessage(mg["message"].(string))
				b2.SendMessage(mg["message"].(string))
			}
		}
	}
	b.Hs["dialog.closed"] = func (msg bot.Message) {
		fmt.Println("b - dialog closed")
		b.Dialog_id = 0
		b.StartSearch()
	}

	b1.Hs["auth.successToken"] = func(msg bot.Message) {
		fmt.Println("b1 - success token, great job!")
		mg := msg["data"].(map[string]interface{})
		b1.Id = mg["id"].(float64)
		b1.Dialog_id = 0
		b1.StartSearch()
	}
	b1.Hs["dialog.opened"] = func(msg bot.Message) {
		mg := msg["data"].(map[string]interface{})
		b1.Dialog_id = mg["id"].(float64)
		fmt.Println("b1 dialog opened")
	}
	b1.Hs["messages.new"] = func (msg bot.Message) {
		mg := msg["data"].(map[string]interface{})
		if mg["senderId"] != b1.Id {
			fmt.Println("b1 -",mg["message"].(string))
			if b.Dialog_id != 0 && b2.Dialog_id != 0 && b1.Dialog_id != 0{
				b.SendMessage(mg["message"].(string))
				b2.SendMessage(mg["message"].(string))
			}
		}
	}
	b1.Hs["dialog.closed"] = func (msg bot.Message) {
		fmt.Println("b1 dialog closed")
		b1.Dialog_id = 0
		b1.StartSearch()
	}

	//b.SetCookie("")
	//go b.Connect("")

	//b1.SetCookie("")
	//go b1.Connect("")

	reader := bufio.NewReader(os.Stdin)reader := bufio.NewReader(os.Stdin)
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
				b1.LeaveDialog()
				b2.LeaveDialog()
			case "close":
				b.Close();
				b1.Close();
				b2.Close();
				break
			default:
				wo := ""
				for _, word := range args {
					wo += word+" "
				}
				b.SendMessage(wo)
				b1.SendMessage(wo)
		}
	}
}
