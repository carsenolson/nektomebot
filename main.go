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
    b.Index = strconv.Itoa(i)
    b.Hs["auth.successToken"] = func(msg bot.Message) {
      fmt.Println("b"+b.Index+" - success token, great job!")
      mg := msg["data"].(map[string]interface{})
      b.Id = mg["id"].(float64)
      b.Dialog_id = 0
      b.StartSearch()
    }
    b.Hs["dialog.opened"] = func(msg bot.Message) {
      mg := msg["data"].(map[string]interface{})
      b.Dialog_id = mg["id"].(float64)
      fmt.Println("b"+b.Index+" - dialog opened")
    }
    b.Hs["dialog.closed"] = func(msg bot.Message) {
      fmt.Println("b"+b.Index+" - dialog closed")
      b.Dialog_id = 0
      b.StartSearch()
    }
    b.Hs["message.new"] = func(msg bot.Message) {
      mg := msg["data"].(map[string]interface{})
      if mg["senderId"] != b.Id {
        fmt.Println(mg["message"].(string))
        for _, bb := range bots {
          if bb.Dialog_id != 0 {
            bb.SendMessage(bb.Index+" "+mg["message"].(string))
          }
        }
      }
    }
    fmt.Println(b.Hs["auth.successToken"])
  }

  //for i, b := range bots {
  //  fmt.Println(b.Index, b.Hs["auth.successToken"], b.Hs["message.new"], tokens[i])
  //}
  //for i, b := range bots {
  //  go b.Connect(tokens[i])
  //}
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		args := strings.Split(text, " ")
		switch string(args[0]) {
			case "track":
				bots[0].Track()
			case "untrack":
				bots[0].UnTrack()
			case "search":
				bots[0].StartSearch()
			case "leave":
				bots[0].LeaveDialog()
			case "close":
				bots[0].Close()
			case "connect":
				bots[0].Connect(os.Args[1])
				break
			default:
				wo := ""
				for _, word := range args {
					wo += word+" "
				}
				bots[0].SendMessage(wo)
		}
	}
}
