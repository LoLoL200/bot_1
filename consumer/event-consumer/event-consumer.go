package event_consumer

import (
	event "bot_dictionary_draft/events"
	"log"
	"time"
)

type Consumer struct {
	fatcher   event.Fetcher
	processor event.Processor
	bathSize  int
}

func New(fatcher event.Fetcher, processor event.Processor, bathSize int) Consumer {
	return Consumer{
		fatcher:   fatcher,
		processor: processor,
		bathSize:  bathSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fatcher.Fetch(c.bathSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}
func (c *Consumer) handleEvents(events []event.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())
			continue
		}
	}
	return nil
}
