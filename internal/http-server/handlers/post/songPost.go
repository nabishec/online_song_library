package post

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/nabishec/restapi/internal/http-server/handlers/decoder"
	"github.com/nabishec/restapi/internal/lib/logger/slerr"
	"github.com/nabishec/restapi/internal/model"
	"github.com/nabishec/restapi/internal/storage"
)

type SongAddingImp interface {
	AddSong(song *model.Song) error
	AddSongDetail(song *model.Song, songDetail *model.SongDetail) error
}

func SongPost(log *slog.Logger, songAdding SongAddingImp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.post.addsong.SongPost()"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		song, errStr := decoder.SongDecoderValJSON(log, r)
		if errStr != nil {

			w.WriteHeader(http.StatusBadRequest) // 400
			render.JSON(w, r, errStr)
			return
		}

		err := songAdding.AddSong(song)
		if errors.Is(err, storage.ErrSongAlreadyExists) {
			log.Info("song already exist", slog.String("song: ", song.SongName+
				":"+song.GroupName))

			w.WriteHeader(http.StatusConflict) //409
			render.JSON(w, r, model.StatusError("song already exist"))

			return
		}
		if err != nil {
			log.Error("failed to add song", slerr.Err(err))

			w.WriteHeader(http.StatusInternalServerError) // 500
			render.JSON(w, r, model.StatusError("failed to add song"))
			return
		}

		log.Info("song added")
		//TODO: Get to api
		render.JSON(w, r, model.OK())
	}
}
