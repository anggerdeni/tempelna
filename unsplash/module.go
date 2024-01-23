package unsplash

import (
	"io"
	"os"
)

type Unsplash interface {
	GetImage() (io.Reader, error)
}

type unsplash struct{}

func New() Unsplash {
	return &unsplash{}
}

func (u *unsplash) GetImage() (io.Reader, error) {
	reader, err := os.Open("assets/img.jpg")
	if err != nil {
		return nil, err
	}
	return reader, nil
}
