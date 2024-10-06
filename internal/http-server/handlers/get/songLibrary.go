package get

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/nabishec/restapi/internal/lib/logger/slerr"
	"github.com/nabishec/restapi/internal/model"
)

type SongLibraryImp interface {
	GetSongLibrary(songName string, groupName string, limit int64, offset int64, log *slog.Logger) ([]*model.Song, error)
	CountNumberOfSong(song string, group string) (int64, error)
}

// @Summary      Get Song Library
// @Tags         songslibrary/song
// @Description  Retrieve the song library with pagination options.
// @Produce      json
// @Param        song    query     string  false "Name of the song"   Example: "Song1"
// @Param        group   query     string  false "Name of the group"  Example: "Group1"
// @Param        first   query     int64   false "Number of items to return"  Example: 10
// @Param        after   query     int64   false "Offset from which to return items" Example: 0
// @Success      200     {object}  model.Response      "OK"
// @Failure      400     {object}  model.Response         "Bad request"
// @Failure      404     {object}  model.Response         "No songs matching the request"
// @Failure      500     {object}  model.Response         "Failed to get song library"
// @Router       /songslibrary [get]
func SongsLibrary(log *slog.Logger, songLibraryImp SongLibraryImp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.put.SongsLibrary()"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		songName := r.URL.Query().Get("song")
		groupName := r.URL.Query().Get("group")
		firstStr := r.URL.Query().Get("first")
		afterStr := r.URL.Query().Get("after")

		var first int64
		var err error
		if firstStr == "" {
			first = 10
		} else {
			first, err = strconv.ParseInt(firstStr, 10, 64)
			if err != nil {
				log.Error("failed converting of first:", slerr.Err(err))

				w.WriteHeader(http.StatusBadRequest) // 400
				render.JSON(w, r, model.StatusError("incorrect value of first"))
				return
			}
		}

		var after int64
		if afterStr == "" {
			after = 0
		} else {
			after, err = strconv.ParseInt(afterStr, 10, 64)
			if err != nil {
				log.Error("failed converting of after:", slerr.Err(err))

				w.WriteHeader(http.StatusBadRequest) // 400
				render.JSON(w, r, model.StatusError("incorrect value of after"))
				return
			}
		}

		library, err := songLibraryImp.GetSongLibrary(songName, groupName, first, after, log)
		if err != nil {
			log.Error("failed get library", slerr.Err(err))

			w.WriteHeader(http.StatusInternalServerError) // 500
			render.JSON(w, r, model.StatusError("failed get song library"))
			return
		}
		if len(library) == 0 {
			log.Info("there wasn't single song matching request")

			w.WriteHeader(http.StatusNotFound) // 404
			render.JSON(w, r, model.StatusError("there wasn't single song matching request"))
			return
		}
		resp := paginationLibrary(library, after, songName, groupName, songLibraryImp, log)

		log.Info("song library getted")
		render.JSON(w, r, resp)
	}
}

func paginationLibrary(library []*model.Song, after int64, songName string, groupName string, accesDBFunc SongLibraryImp, log *slog.Logger) model.Response {
	edges := make([]*model.SongEdge, 0, len(library))

	songsNumber, err := accesDBFunc.CountNumberOfSong(songName, groupName)
	if err != nil {
		log.Error("can't count songs", slerr.Err(err))
		songsNumber = 0
	}

	for i, val := range library {
		edges = append(edges, &model.SongEdge{
			Node:   val,
			Cursor: (int64(i) + after + 1),
		})
	}

	endCursor := after + int64(len(library))
	hasNextPage := endCursor < songsNumber
	if songsNumber == 0 {
		hasNextPage = false
	}

	return model.Response{
		Status: "OK",
		SongsLibrary: &model.SongsConnection{
			Edges: edges,
			PageInfo: &model.LibraryPageInfo{
				EndCursor:   &endCursor,
				HasNextPage: hasNextPage,
			},
		},
	}

}
