package telegram

import (
	"context"
	"errors"
	"log"
	url2 "net/url"
	"read-adviser-bot/lib/er"
	"read-adviser-bot/storage"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *TgProcessor) DoCmd(text string, chatId int, username string) error {
	// удаляем из текста лишние пробелы
	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from '%s'", text, username)

	// нужно добавить следующие команды:
	// save page - http:/././.
	// rnd page - /rng
	// help - /help
	// start - /start

	url, err := url2.Parse(text)
	// команда является save page
	if err == nil && url.Host != "" {
		return p.savePage(chatId, url.String(), username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatId, username)
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	default:
		return p.tg.SendMessage(chatId, msgUnknownCommand)
	}
}

func (p *TgProcessor) savePage(chatId int, pageUrl string, username string) (err error) {
	defer func() {
		err = er.Wrap("Cant do command 'save page' !", err)
	}()

	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}
	exists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}
	if exists {
		return p.tg.SendMessage(chatId, msgAlreadyExists)
	}
	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}
	if err := p.tg.SendMessage(chatId, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *TgProcessor) sendRandom(chatId int, username string) (err error) {
	defer func() { err = er.Wrap("Cant do command 'send random' !", err) }()
	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrorNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrorNoSavedPages) {
		return p.tg.SendMessage(chatId, msgNoSavedPages)
	}
	if err := p.tg.SendMessage(chatId, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(context.Background(), page)
}

func (p *TgProcessor) sendHelp(chatId int) error {
	return p.tg.SendMessage(chatId, msgHelp)
}

func (p *TgProcessor) sendHello(chatId int) error {
	return p.tg.SendMessage(chatId, msgHello)
}
