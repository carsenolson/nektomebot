package main

import (
  "./bot"
	"strings"
	"bufio"
	"os"
  "log"
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

  for _, b := range bots {
    b.Bots = &bots;
  }

  for i, b := range bots {
    go b.Connect(tokens[i])
  }

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
