package sqlite

import (
	"context"
	// это интерфейс для взаимодействия со всеми
	// реляционными БД
	"database/sql"
	"fmt"
	"read-adviser-bot/storage"
	// очень важно подключить саму библиотеку(тк она сторонний ресурс)
	_ "github.com/mattn/go-sqlite3"
)

// тип кой будет реализовывать интерфейс
type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage.
// path - путь до папки с БД
func New(path string) (*Storage, error) {
	// sql.Open("sqlite3") - сообщаем что будем работать
	// именно с базой данных sqlite3
	// db это некая сущность для взаимодействия с БД
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Cant open database: %w !", err)
	}
	// db.Ping() - проверка что нам удалось
	// установить соединение с БД
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Cant connect to database: %w", err)
	}
	return &Storage{db: db}, nil
}

// теперь реализуем все методы интерфейса Storage
// Save saves page to storage
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	// подготавливаем запрос
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`
	// выполняем запрос
	// контекст это особенная сущность в GO
	// важная хрень короче, без неё ни 1 большой
	// проект не будет работать
	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)
	// _ это некий Результат(сейчас он не важен)
	// Result был бы важен если бы мы хотети что-то
	// дальше сделать с полученными данными
	if err != nil {
		return fmt.Errorf("Cant save page: %w", err)
	}
	return nil
}

// выше именно Exec(Context) тк Exec исп. тогда когда
// функция ДБ не возвращает поля (insert/delete/update)

// PickRandom picks random page from storage.
func (s *Storage) PickRandom(ctx context.Context, UserName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`
	row := s.db.QueryRowContext(ctx, q, UserName)
	// QueryRowContext вернёт нам сущность Row
	// там нет данных в прямом виде (тк в теории они
	// могут быть избыточными и большими)
	// поэтому нужно использовать функцию Scan
	var url string
	err := row.Scan(&url)
	if err == sql.ErrNoRows {
		// не возвращаем ошибку тк Пользователь ещё ничего не сохранил
		return nil, storage.ErrorNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("Cant pick random page: %w", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: UserName,
	}, nil
}

// здесь именно Query тк функция ДБ (SELECT)
// будет возвращать поля

// Remove removes page from storage.
func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`
	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)
	if err != nil {
		return fmt.Errorf("Cant remove page: %w", err)
	}
	return nil
}

// IsExists check that page is exists.
func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`
	var count int
	if err := s.db.QueryRowContext(ctx, q, p.URL, p.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("Cant check if page exists: %w", err)
	}
	return count > 0, nil
}

// ещё нужно инициализировать нашу БД
// по факту всё что будет делать эта функция -
// это создание таблицы Pages
func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`
	if _, err := s.db.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("Cant create table: %w", err)
	}
	return nil
}
