package telegram

import (
	"bot_dictionary_draft/clients/telegram"
	"bot_dictionary_draft/events"
	e "bot_dictionary_draft/lib/errors"
	"bot_dictionary_draft/lib/storage"
	"errors"
)

type ProcessorTG struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}
type Meta struct {
	ChatID   int
	UserName string
}

// Local event types to avoid import cycle with bot_dictionary_draft/events/telegram
type Type int

const (
	Unkdown Type = iota
	Message
)

// type Event struct {
// 	Type Type
// 	Text string
// 	Meta interface{}
// }

// ERROR

var ErrUnknownEventType = errors.New("unknown event type")

func New(client *telegram.Client, storage storage.Storage) *ProcessorTG {
	return &ProcessorTG{
		tg:      client,
		storage: storage,
	}
}

func (p *ProcessorTG) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}
	if len(updates) == 0 {
		return nil, nil
	}
	res := make([]events.Event, 0, len(updates))
	for _, u := range updates {
		res = append(res, event(u))
	}
	p.offset = updates[len(updates)-1].Update_id + 1
	return res, nil
}

func (p *ProcessorTG) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMassage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}
func (p *ProcessorTG) processMassage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.UserName); err != nil {
		return e.Wrap("can't process message", err)
	}
	return nil
}
func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownEventType)
	}
	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.From.UserName,
		}
	}
	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unkdown
	}
	return events.Message
}
func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}
