package event_consumer

import (
	"log"
	"read-adviser-bot/events"
	"time"
)

type EventConsumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int // сколько событий мы будем обрабатывать за 1 раз
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *EventConsumer {
	return &EventConsumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

// метод Start будет в бесконечном цикле получать новые события и обрабатывать их
func (ec *EventConsumer) Start() error {
	for {
		events, err := ec.fetcher.Fetch(ec.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}
		if len(events) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := ec.handleEvents(events); err != nil {
			log.Printf("[ERR] Cant Handle Events: %s", err.Error())
			continue
		}
	}
}

/*
!1 - потеря событий (то есть если не смогли обработать, то Fetcher
пойдет дальше и будет выдавать новые события, те потеряем навсегда обработку пред. события)
$1 - механизм ReTry - делать несколько попыток обработки (а вдруг проблема с сетью)
$2 - мех BackUp - сохранять необработанные события в некоторое
вспомогательное хранилище, чтобы потом добавить их в основное
!2 - если пытаться обработать Каждое события(когда у Всех при обработке будет ошибка)
(например когда отвалится Сеть) с таймутом, то обработка большой пачки займёт много времени
$1 - ввести счётчик ошибок - Еx: >5 подряд ошибок <==> прекращаем обработку пачки
#1 - можно добавить асинхронную обработку Событий (понадобится sync.WaitGroup)
#2 - Подтверждение для Фетчера, т.e. пока не обработает текущую пачку, не сдвинется вперед
*/
func (ec *EventConsumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)
		if err := ec.processor.Process(event); err != nil {
			log.Printf("cant handle event: %s", err.Error())
			continue
		}
	}
	return nil
}
