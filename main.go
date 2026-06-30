package main

import (
	tgClient "bot_dictionary_draft/clients/telegram"
	event_consumer "bot_dictionary_draft/consumer/event-consumer"
	"os"

	telegram "bot_dictionary_draft/events/telegram"

	"bot_dictionary_draft/lib/storage/files"
	"flag"
	"log"

	"github.com/joho/godotenv"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = ""
	bathSize    = 50
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Panicln("Your .env file don`t have TOKEN")
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath))
	log.Println("service started")
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, bathSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
	//fetcher = fetcher,New(tgClient)
	// processor = processor.New(tgClient)
	//consumer.Start(fetcher,processor)

}

func mustToken() string {

	token := flag.String(
		"tg-bot-token",
		storagePath,
		"token for access to telegram bot")

	flag.Parse()
	if *token == "" {
		*token = os.Getenv("BOT_TOKEN")
	}
	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}
