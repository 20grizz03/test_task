package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type SongDetail struct {
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

func FetchSongDetails(group string, song string) (*SongDetail, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", os.Getenv("API_URL"), group, song)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var details SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return &details, nil
}
