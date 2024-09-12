package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"read-adviser-bot/lib/er"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, UserName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

var ErrorNoSavedPages = errors.New("Dont have saved pages!")

func (p *Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", er.Wrap("cant calculate hash from URL", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", er.Wrap("cant calculate hash from UserName", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
