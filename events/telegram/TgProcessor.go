package telegram

import (
	"errors"
	"read-adviser-bot/clients/telegram"
	"read-adviser-bot/events"
	"read-adviser-bot/lib/er"
	"read-adviser-bot/storage"
)

type TgProcessor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatId   int
	Username string
}

var ErrorUnknownEventType = errors.New("Error: Unknown event type!")
var ErrorUnknownMetaType = errors.New("Error: Unknown meta type!")

func New(client *telegram.Client, storage storage.Storage) *TgProcessor {
	return &TgProcessor{
		tg:      client,
		offset:  0,
		storage: storage,
	}

}

func (p *TgProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.GetUpdates(p.offset, limit)
	if err != nil {
		return nil, er.Wrap("Cant fetch updates", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	//tgProc.offset = tgProc.offset + 1
	// если я правильно понял то нужно пройтись па массиву updates и создать Event на каждый
	// update и все это сохранить в массив []Event'ov
	result := make([]events.Event, 0, len(updates))
	for _, update := range updates {
		result = append(result, makeEvent(update))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return result, nil
}

func (p *TgProcessor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return ErrorUnknownEventType
	}
}

func (p *TgProcessor) processMessage(event events.Event) error {
	meta, err := getMeta(event)
	if err != nil {
		return er.Wrap("Error: Cant process message!", err)
	}
	if err := p.DoCmd(event.Text, meta.ChatId, meta.Username); err != nil {
		return er.Wrap("Error: Cant process message!", err)
	}
	return nil
}

func getMeta(event events.Event) (Meta, error) {
	// так называемый Type Assertion - проверка на тип
	//
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, ErrorUnknownMetaType
	}
	return res, nil
}

func makeEvent(update telegram.Update) events.Event {
	updateType := fetchType(update)

	result := events.Event{
		Type: updateType,
		Text: fetchText(update),
	}

	if updateType == events.Message {
		result.Meta = Meta{
			ChatId:   update.Message.Chat.Id,
			Username: update.Message.From.Username,
		}
	}

	return result
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}
