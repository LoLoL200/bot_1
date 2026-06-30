package telegram

import (
	"bot_dictionary_draft/clients/telegram"
	e "bot_dictionary_draft/lib/errors"
	"bot_dictionary_draft/lib/storage"
	"errors"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCms  = "/help"
	StartCmd = "/start"
	LinkCmd  = "/firstlink"
)

func (p *ProcessorTG) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command '%s'from'%s'", text, username)
	//	add page: http://...
	if isAddCmd(text) {
		// TODO: AddPage()
		return p.savePage(chatID, text, username)
	}
	//	rnd page: /rnd
	//	help: /help
	//	start: /start hi + help
	switch text {
	case LinkCmd:
		return p.sendLink(chatID, username)
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCms:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)

	}
}
func (p *ProcessorTG) sendLink(chatID int, username string) (err error) {
	defer func() {
		err = e.WrapIfErr("can' do command: can't send links", err)
	}()
	page, err := p.storage.CheckLinks(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavePage) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavePage) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(page)
}
func (p *ProcessorTG) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() {
		err = e.WrapIfErr("can't do CMD save", err)
	}()
	sendMsg := NewMessageSender(chatID, p.tg)
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}
	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return sendMsg(msgAlreadyExists)
		//p.tg.SendMessage(chatID, msgAlreadyExists)
	}
	if err := p.storage.Save(page); err != nil {
		return err
	}
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}
func (p *ProcessorTG) sendRandom(chatID int, username string) (err error) {
	defer func() {
		err = e.WrapIfErr("can' do command: can't send random", err)
	}()
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavePage) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavePage) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(page)
}
func (p *ProcessorTG) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}
func (p *ProcessorTG) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}
func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}
func isURL(text string) bool {
	// https://ya.ua or  http://ya.ua
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
