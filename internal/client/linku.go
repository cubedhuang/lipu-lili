package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cubedhuang/lipu-lili/internal/models"
)

const linkuWordsAPIURL = "https://api.linku.la/v1/words?lang=*"
const linkuSignsAPIURL = "https://api.linku.la/v1/luka_pona/signs?lang=*"

func FetchLinku() (*models.LinkuData, error) {
	var words *models.WordsData
	if err := fetchJSON(linkuWordsAPIURL, &words); err != nil {
		return nil, fmt.Errorf("failed to fetch linku data: %w", err)
	}

	var signs *models.SignsData
	if err := fetchJSON(linkuSignsAPIURL, &signs); err != nil {
		return nil, fmt.Errorf("failed to fetch linku data: %w", err)
	}

	return &models.LinkuData{
		Words: *words,
		Signs: *signs,
	}, nil
}

func fetchJSON(url string, v interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return nil
}
