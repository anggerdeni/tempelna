package unsplash

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	unsplashclient "github.com/hbagdi/go-unsplash/unsplash"
	"golang.org/x/oauth2"
)

type Unsplash interface {
	GetImage() (io.Reader, error)
}

type unsplash struct {
	c *unsplashclient.Unsplash
}

func New(accessKey string) Unsplash {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: fmt.Sprintf("Client-ID %s", accessKey)},
	)
	client := oauth2.NewClient(context.Background(), ts)
	return &unsplash{
		c: unsplashclient.New(client),
	}
}

func (u *unsplash) GetImage() (io.Reader, error) {
	photo, _, err := u.c.Photos.Random(&unsplashclient.RandomPhotoOpt{
		Count:         1,
		CollectionIDs: []int{10720924},
	})
	if err != nil {
		return nil, err
	}

	url, _, err := u.c.Photos.DownloadLink(*(*photo)[0].ID)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(bodyBytes)

	return reader, nil
}

func (u *unsplash) GetImageOld() (io.Reader, error) {
	reader, err := os.Open("assets/img.jpg")
	if err != nil {
		return nil, err
	}
	return reader, nil
}
