package main

import (
	"./bot"
	"fmt"
	"strings"
	"bufio"
	"os"
  "log"
  "strconv"
)

func main() {
  var tokens []string
  var bots []*bot.Bot
	b := bot.NewBot()

  file, err := os.Open("tokens/tokens")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    tokens = append(tokens, scanner.Text())
    bots = append(bots, bot.NewBot())
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  for i, b := range bots {
    b.Hs["auth.successToken"] = func(msg bot.Message) {
      fmt.Println("b"+strconv.Itoa(i)+" - success token, great job!")
      mg := msg["data"].(map[string]interface{})
      b.Id = mg["id"].(float64)
      b.Dialog_id = 0
      b.StartSearch()
    }
    b.Hs["dialog.opened"] = func(msg bot.Message) {
      mg := msg["data"].(map[string]interface{})
      b.Dialog_id = mg["id"].(float64)
      fmt.Println("b"+strconv.Itoa(i)+" - dialog opened")
    }
    b.Hs["dialog.closed"] = func (msg bot.Message) {
      fmt.Println("b"+strconv.Itoa(i)+" - dialog closed")
      b.Dialog_id = 0
      b.StartSearch()
    }
    b.Hs["message.new"] = func(msg bot.Message) {
      mg := msg["data"].(map[string]interface{})
      if mg["senderId"] != b.Id {
        fmt.Println(mg["message"].(string))
        for ii, bot := range bots {
          if bot.Dialog_id != 0 {
            bot.SendMessage(strconv.Itoa(ii)+" "+mg["message"].(string))
          }
        }
      }
    }
  }
  for i, bot := range bots {
    go bot.Connect(tokens[i])
  }

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
