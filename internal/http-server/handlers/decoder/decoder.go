package decoder

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nabishec/restapi/internal/lib/logger/slerr"
	"github.com/nabishec/restapi/internal/model"
)

func SongDecoderValJSON(log *slog.Logger, r *http.Request) (*model.Song, *string) {
	var song model.Song

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
		} else {
			log.Error("failed to decode request body", slerr.Err(err))
		}
		reply := "bad request"
		return nil, &reply
	}
	log.Info("request body decoded", slog.Any("song: ", song.SongName+":"+song.GroupName))

	if err := validator.New().Struct(song); err != nil {
		validatorErr := err.(validator.ValidationErrors)

		log.Error("invalid types", slerr.Err(err))

		reply := validatorErr.Error()
		return nil, &reply

	}

	return &song, nil
}
