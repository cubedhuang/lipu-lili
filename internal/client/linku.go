package client

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cubedhuang/lipu-lili/internal/models"
)

const linkuAPIURL = "https://api.linku.la/v1/words?lang=*"

func FetchLinku() (models.LinkuData, error) {
	res, err := http.Get(linkuAPIURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data models.LinkuData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}
