package storage

import (
	e "bot_dictionary_draft/lib/errors"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
	CheckLinks(userName string) (*Page, error)
}

var ErrNoSavePage = errors.New("No saved pages")

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
