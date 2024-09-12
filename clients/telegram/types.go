package telegram

// в этом файле определяем все типы, с которыми работает клиент

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// на самом деле мы будем получать не только Update,
// а одну большую структуру в которой Update это лишь часть
type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
	// поле message может отсутсвовать, мы хотели бы здесь тогда видеть nil
	// поэтому делаем не просто значение а указатель - *
}

type IncomingMessage struct {
	Text string `json:"text"`
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
}

type User struct {
	Username string `json:"username"`
}

type Chat struct {
	Id int `json:"id"`
}
