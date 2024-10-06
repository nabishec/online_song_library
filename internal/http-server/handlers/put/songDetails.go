package put

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
)

type SongPutImp interface {
	PutSongDetail(song *model.Song, songDetail *model.SongDetail) error
	AddSongDetail(song *model.Song, songDetail *model.SongDetail) error
}

type Request struct {
	SongData      model.Song       `json:"dataSong" validate:"required"`
	NewSongDetail model.SongDetail `json:"songDetail" validate:"required"`
}

// @Summary      Add Song Detail
// @Tags         songslibrary/song
// @Description  Add the details of a new song to the library.
// @Accept       json
// @Produce      json
// @Param        request body      Request true  "Request with song data and details" Example: {"dataSong": {"song": "Song1", "group": "Group1"}, "songDetail": {"releaseDate": "2022-01-01", "link": "http://example.com", "text": "This is a great song"}}
// @Success      200         {object}  model.Response    "OK"
// @Failure      400         {object}  model.Response       "Bad request"
// @Failure      500         {object}  model.Response       "Failed to add song detail"
// @Router       /songslibrary/song [put]
func SongDetail(log *slog.Logger, songPutImp SongPutImp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.put.songDetail()"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")
			} else {
				log.Error("failed to decode request body", slerr.Err(err))
			}

			w.WriteHeader(http.StatusBadRequest) // 400
			render.JSON(w, r, model.StatusError("bad request"))
			return
		}

		log.Info("request body decoded", slog.Any("songInf: ", req))

		if err := validator.New().Struct(req); err != nil {
			validatorErr := err.(validator.ValidationErrors)
			log.Error("invalid types", slerr.Err(err))

			w.WriteHeader(http.StatusBadRequest) // 400
			render.JSON(w, r, model.StatusError(validatorErr.Error()))
			return
		}

		err = songPutImp.AddSongDetail(&req.SongData, &req.NewSongDetail)
		if err != nil {
			log.Error("failed to add song detail", slerr.Err(err))
			w.WriteHeader(http.StatusInternalServerError) // 500
			render.JSON(w, r, model.StatusError("failed to add song detail"))
			return

		}

		log.Info("song detail changed")
		render.JSON(w, r, model.OK())
	}
}
