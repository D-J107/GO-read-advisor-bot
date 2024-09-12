package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"read-adviser-bot/lib/er"
	"strconv"
)

type Client struct {
	host     string // хост api-сервиса телеграмма ( по дефолту api.telegram.org )
	basePath string // префикс с которого начинаются все запросы
	// типа такого: tg-bot.com/bot<token>
	client http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) SendMessage(chatId int, text string) error {
	// подготовим параметры запрсоса
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)
	// теперь выполняем этот запрос
	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return er.Wrap("cant SendMessage", err)
	}
	return nil
}

func (c *Client) GetUpdates(offset int, limit int) (updates []Update, err error) {
	defer func() {
		err = er.Wrap("cant GetUpdates", err)
	}()
	// теперь нужно сформировать параметры запроса
	// сделаем это с помощью пакета URL
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	// теперь нужно отправить запрос, ну то есть doRequest()
	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}
	var res UpdatesResponse
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) doRequest(methodName string, query url.Values) (data []byte, err error) {
	defer func() {
		err = er.Wrap("cant do request", err)
	}()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		// функция для соединения через / чтобы вдруг не было 2 слэша
		Path: path.Join(c.basePath, methodName),
	}

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// теперь передаем в request параметры запроса из аргумента
	// метод Encode приведёт параметры к такому виду, который можно отправить на сервер
	request.URL.RawQuery = query.Encode()

	// отправляем получившийся запрос на сервер
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
