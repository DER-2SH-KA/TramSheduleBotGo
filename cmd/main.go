package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DER-2SH-KA/TramSheduleBotGo/cmd/parser"
	"github.com/DER-2SH-KA/TramSheduleBotGo/cmd/telegram"
)

var (
	token  = os.Getenv("TELEGRAM_TOKEN")
	Parser = &parser.Parser{
		Client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
)

func main() {
	botWithParser := &telegram.Bot{
		Token:             token,
		LongPollingTiming: time.Second * 5,
		Parser:            Parser,
	}

	err := botWithParser.Start()
	if err != nil {
		log.Fatal(err)
	}
}
