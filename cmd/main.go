package main

import (
	"log"
	"net/http"
	"time"

	"github.com/DER-2SH-KA/TramSheduleBotGo/cmd/parser"
)

var (
	Parser = &parser.Parser{
		Client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
)

func main() {

	message, err := Parser.Parse(parser.ROUTE_2, parser.DIRECTION_FROM_CH_SLOBODA, parser.WORK_DAYS)
	if err != nil {
		log.Println(err)
	}

	log.Println(message)
}
