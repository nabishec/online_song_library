package get

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/nabishec/restapi/internal/lib/logger/slerr"
	"github.com/nabishec/restapi/internal/model"
	"github.com/nabishec/restapi/internal/storage"
)

type GettingTesxtSongImp interface {
	GetSongText(song *model.Song) (*string, error)
}

func TextSongGet(log *slog.Logger, gettingTesxtSongImp GettingTesxtSongImp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get.textSong.TextSongGet()"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		songName := r.URL.Query().Get("song")
		groupName := r.URL.Query().Get("group")
		if songName == "" || groupName == "" {
			log.Error("request incimplete")

			w.WriteHeader(http.StatusBadRequest) //400
			render.JSON(w, r, model.StatusError("request is incomplete"))
			return
		}

		firstStr := r.URL.Query().Get("first")
		afterStr := r.URL.Query().Get("after")

		var first int
		var err error
		if firstStr == "" {
			first = 2
		} else {
			first, err = strconv.Atoi(firstStr)
			if err != nil {
				log.Error("failed to convert 'first' value", slerr.Err(err))

				w.WriteHeader(http.StatusBadRequest) //400
				render.JSON(w, r, model.StatusError("incorrect value of first"))
				return
			}
		}

		var after int
		if afterStr == "" {
			after = 1
		} else {
			after, err = strconv.Atoi(afterStr)
			if err != nil {
				log.Error("failed to convert 'after' value", slerr.Err(err))

				w.WriteHeader(http.StatusBadRequest) //400
				render.JSON(w, r, model.StatusError("incorrect value of after"))
				return
			}
		}

		song := &model.Song{
			SongName:  songName,
			GroupName: groupName,
		}

		text, err := gettingTesxtSongImp.GetSongText(song)
		if errors.Is(err, storage.ErrSongNotFound) {
			log.Info("song doesn't exist", slog.String("song:", song.SongName+
				":"+song.GroupName))

			w.WriteHeader(http.StatusNotFound) //404
			render.JSON(w, r, model.StatusError("song doesn't exist"))
			return
		}
		if err != nil {
			log.Error("failed getiing text of song", slerr.Err(err))

			w.WriteHeader(http.StatusInternalServerError) //500
			render.JSON(w, r, model.StatusError("failed getting text of song"))
			return
		}

		resp := pagination(text, first, after)

		log.Info("song text retrieved successfully")
		render.JSON(w, r, resp)
	}

}

func pagination(text *string, first int, after int) model.Response {
	var edges []*model.CoupletEdge

	couplets := strings.Split(*text, "\n\n")

	startInd := after - 1
	endIndex := startInd + first

	var endCursor int

	for i, length := startInd, len(couplets); i < length || i < endIndex; i++ {
		edges = append(edges, &model.CoupletEdge{
			Node:   &couplets[i],
			Cursor: i + 1,
		})
		if i < (length-1) || i < (endIndex-1) {
			endCursor = i
		}
	}

	hasNextPage := endCursor < len(couplets)

	return model.Response{
		Status: "OK",
		SongText: &model.TextConnection{
			Edges: edges,
			PageInfo: &model.TextPageInfo{
				EndCursor:   &endCursor,
				HasNextPage: hasNextPage,
			},
		},
	}

}
