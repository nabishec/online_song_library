package deletion

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/nabishec/restapi/internal/lib/logger/slerr"
	"github.com/nabishec/restapi/internal/model"
	"github.com/nabishec/restapi/internal/storage"
)

// @Summary      Delete a Song
// @Tags         songdelete/song
// @Description  Delete a song from the library by song name and group name.
// @Produce      json
// @Param        song    query     string  true  "Name of the song"   Example: "Song1"
// @Param        group   query     string  true  "Name of the group"  Example: "Group1"
// @Success      200     {object}  model.Status  "OK"
// @Failure      400     {object}  model.Error    "Bad request"
// @Failure      404     {object}  model.Error    "Song doesn't exist"
// @Failure      500     {object}  model.Error    "Failed deletion of song"
// @Router       /songslibrary/song [delete]

type SongDeletingImp interface {
	DeleteSong(song *model.Song, log *slog.Logger) error
}

func SongDelete(log *slog.Logger, songDeleting SongDeletingImp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.songDelete.SongDelete()"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		songName := r.URL.Query().Get("song")
		groupName := r.URL.Query().Get("group")
		if songName == "" || groupName == "" {
			log.Error("request incomplete")

			w.WriteHeader(http.StatusBadRequest) // 400
			render.JSON(w, r, model.StatusError("request is incomplete"))
			return
		}

		song := &model.Song{
			SongName:  songName,
			GroupName: groupName,
		}

		err := songDeleting.DeleteSong(song, log)
		if errors.Is(err, storage.ErrSongNotFound) {
			log.Info("song doesn't exist", slog.String("song:", song.SongName+
				":"+song.GroupName))

			w.WriteHeader(http.StatusNotFound) // 404
			render.JSON(w, r, model.StatusError("song doesn't exist"))
			return
		}
		if err != nil {
			log.Error("failed delete song", slerr.Err(err))

			w.WriteHeader(http.StatusInternalServerError) // 500
			render.JSON(w, r, model.StatusError("failed deletion of song"))
			return
		}

		log.Info("song deleted", slog.String("song:", song.SongName+":"+song.GroupName))
		render.JSON(w, r, model.OK())
	}

}
