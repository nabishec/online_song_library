package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/nabishec/restapi/internal/model"
)

func GetSongDetailsOfExternalApi(song *model.Song) (*model.SongDetail, error) {
	const op = "external.serviceapi.GetSongDetailsOfExternalApi()"

	baseURL := os.Getenv("EXTERNAL_API_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("%s:failed url(%s)", op, baseURL)
	}
	baseURL += "/info"

	reqParameters := url.Values{}
	reqParameters.Add("group", song.GroupName)
	reqParameters.Add("song", song.SongName)

	reqURL := fmt.Sprintf("%s?%s", baseURL, reqParameters.Encode())

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to external service: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("external service error")
	case http.StatusBadRequest:
		return nil, fmt.Errorf("incorrect request to external api: request URL(%s)", reqURL)
	}

	var songDetail model.SongDetail

	err = json.NewDecoder(resp.Body).Decode(&songDetail)
	if err != nil {
		return nil, fmt.Errorf("error decoding response of external service:%w", err)
	}

	return &songDetail, nil

}
