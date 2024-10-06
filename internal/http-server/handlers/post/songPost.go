package post

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/nabishec/restapi/internal/clients"
	"github.com/nabishec/restapi/internal/http-server/handlers/decoder"
	"github.com/nabishec/restapi/internal/lib/logger/slerr"
	"github.com/nabishec/restapi/internal/model"
	"github.com/nabishec/restapi/internal/storage"
)

type SongAddingImp interface {
	AddSong(song *model.Song) error
	AddSongDetail(song *model.Song, songDetail *model.SongDetail) error
}

// @Summary      Add Song
// @Tags         songslibrary/song
// @Description  Add a new song to the library and fetch its details from an external API.
// @Accept       json
// @Produce      json
// @Param        songData  body      model.Song       true  "Song Data"      Example: {"songName": "Song1", "groupName": "Group1", "releaseDate": "2022-01-01"}
// @Success      200       {object}  model.Response    "OK"
// @Failure      400       {object}  model.Response       "Bad request"
// @Failure      409       {object}  model.Response       "Song already exists"
// @Failure      500       {object}  model.Response       "Failed to add song"
// @Failure      207       {object}  model.Response       "Failed to get song details"
// @Router       /songslibrary/song [post]
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
		songDetail, err := clients.GetSongDetailsOfExternalApi(song)
		if err != nil {
			log.Error("failed to get song details", slerr.Err(err))

			w.WriteHeader(http.StatusMultiStatus) // 207
			render.JSON(w, r, model.StatusError("failed to get song details"))
			return
		}
		err = songAdding.AddSongDetail(song, songDetail)
		if err != nil {
			log.Error("failed to add song details", slerr.Err(err))

			w.WriteHeader(http.StatusMultiStatus) // 207
			render.JSON(w, r, model.StatusError("failed to add song details"))
			return
		}

		log.Info("song added")
		render.JSON(w, r, model.OK())
	}
}
