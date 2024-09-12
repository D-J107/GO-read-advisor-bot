package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	// этих двух полей недостаточно для общего пользования. => делаем
	// поле с метаинформацией чтобы понимать например что Event это от Telegram'a
	Meta interface{}
}

// ещё раз - разница между Event и Update - Event более общее,
//а Update понятие конкретно для телеграмма
