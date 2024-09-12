package main

import (
	"context"
	"flag"
	"log"
	tgClient "read-adviser-bot/clients/telegram"
	event_consumer "read-adviser-bot/consumer/event-consumer"
	tgProcessor "read-adviser-bot/events/telegram"
	"read-adviser-bot/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

// 6570067236:AAGoSAkzATHWnKu2UH5z4wlB9TwU6rFqzjU
func main() {
	//s := files.New(storagePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("Cant connect to SQLite storage: ", err)
	}
	// TODO возвращает context.Background() по дефолту
	// но так мы помечаем себе что ещё точно не определилсь
	// какой именно контекст будем использовать(он может быть разным)
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("Cant init SQLite storage: ", err)
	}
	// и да, context.Backround() это контекст по дефолту
	// Он не имеет ограничений
	// в будущем возможно понадобиться контекст с Time-Out

	eventsProcessor := tgProcessor.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("Service started!")
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	err = consumer.Start()
	if err != nil {
		log.Fatal("Service is stopped!", err)
	}
	// fetcher := fetcher.New(tgClient)
	// processor := processor.New(tgClient)
	// consumer.Start(fetcher, processor)

	// Токен это удостоверение User'a
	// Получаем токен из флагов запуска приложения
	// token = flags.Get(token)
	// tgClient = telegram.New(token)
	// fetcher = new fetcher(tgClient)
	// processor = new processor(tgClient)

	// и Fetcher и Processor будут общаться через API телеграма
	// Fetcher будет отправлять туда запрос чтобы получить новые события
	// Processor будет обрабатывать сообщения и отправлять ответ

	// consumer.Start(fetcher, processor)

	//// инициализация контекста с таймаутом == 15 секунд
	//ctx, cancel := context.WithTimeout(context.Background(),
	//	15*time.Second)
	//// defer cancel() гарантирует что после выхода из функции
	//// или go-routine'ы контекст будет отменён
	//defer cancel()
	//
	//req, err := http.NewRequestWithContext(ctx, http.MethodGet,
	//	"https://example.com", nil)
	//if err != nil {
	//	return nil, fmt.Errorf("Failed to create request with ctx: %w", err)
	//}
	//
	//result, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	return nil, fmt.Errorf("Failed to perform http request: %w", err)
	//}
	//
	//return result, nil
}

// Consumer постоянно получает и обрабатывает события
// Всё что умеет делать Fetcher это получать события
// Processor обрабатывает события
// Client
// Event - событие - всё что получаем из Telegram-bot'a

func mustToken() string {
	// во время запуска программы укажем в виде переменной
	// bot-tg-bot-token 'my_token'
	token := flag.String(
		"tg-bot-token",
		"",
		"token to access telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}

// $ Context используется чтобы из внешней функции
// остановить исполнение внутренней функции
