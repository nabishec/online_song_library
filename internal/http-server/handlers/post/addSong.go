package post

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nabishec/restapi/internal/lib/logger/slerr"
	"github.com/nabishec/restapi/internal/model"
	"github.com/nabishec/restapi/internal/storage"
)

type SongAddingImp interface {
	AddSong(song *model.Song) error
	AddSongDetail(song *model.Song, songDetail *model.SongDetail) error
}

func NewSongHandler(log *slog.Logger, songAdding SongAddingImp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.addsong.NewAddSongHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var song model.Song
		err := json.NewDecoder(r.Body).Decode(&song)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")
			} else {
				log.Error("failed to decode request body", slerr.Err(err))
			}
			render.JSON(w, r, model.StatusError("bad request"))
			return
		}
		log.Info("request body decoded", slog.Any("song: ", song.SongName+":"+song.GroupName))

		if err := validator.New().Struct(song); err != nil {
			validatorErr := err.(validator.ValidationErrors)

			log.Error("invalid types", slerr.Err(err))

			render.JSON(w, r, model.StatusError(validatorErr.Error()))

			return
		}

		err = songAdding.AddSong(&song)
		if errors.Is(err, storage.ErrMusikAlreadyExists) {
			log.Info("musik already exists", slog.String("song: ", song.SongName+":"+song.GroupName))

			render.JSON(w, r, model.StatusError("song already exixsts"))

			return
		}
		if err != nil {
			log.Error("failed to add song", slerr.Err(err))

			render.JSON(w, r, model.StatusError("failed to add song"))

			return
		}

		log.Info("song added")
		//TODO: Get to api
		render.JSON(w, r, model.OK())
	}
}
