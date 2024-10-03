package postgresql

import (
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/nabishec/restapi/internal/model"
	"github.com/nabishec/restapi/internal/storage"
)

func (r *Database) AddSong(song *model.Song) error {
	const op = "internal.storage.postgresql.AddSong()"

	if _, err := r.foundSongId(song); err == nil {
		return fmt.Errorf("%s:%w", op, storage.ErrMusikAlreadyExists)
	}

	_, err := r.DB.Exec("INSERT INTO songs (song_name, group_name) VALUES ($1, $2)",
		song.SongName, song.GroupName)

	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	//if lastinsertid isnt in postgres
	//err := r.DB.QueryRow("INSERT INTO songs (song_name, group_name) VALUES ($1, $2)  RETURNING id",song.SongName, song.GroupName).Scan(&id)

	// id, err := res.LastInsertId()
	// if err != nil {
	// 	return fmt.Errorf("%s:%w", op, err)
	// }

	return nil
}

func (r *Database) DeleteSong(song *model.Song) error {
	const op = "internal.storage.postgresql.DeleteSong()"

	_, err := r.DB.Exec("DELETE FROM songs WHERE song_name = $1 AND group_name = $2",
		song.SongName, song.GroupName)

	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (r *Database) PutSongDetail(song *model.Song, songDetail *model.SongDetail) error {
	const op = "internal.storage.postgresql.PutSongDetail()"

	songId, err := r.foundSongId(song)
	if err != nil {
		return err
	}
	id, err := r.foundSongDetailId(songId)
	if err != nil {
		return err
	}
	_, err = r.DB.Exec("UPDATE songs_detail SET release_date = $1, link = $2, text = $3 WHERE id = $4",
		songDetail.ReleaseDate, songDetail.Link, songDetail.Text, id)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (r *Database) GetMusicLibrary(song *model.Song, limit int, offset int) ([]*model.Song, error) {
	const op = "internal.storage.postgresql.GetMusicLibrary()"

	var library []*model.Song

	query := "SELECT song_name,group_name FROM songs WHERE TRUE"
	args := []interface{}{}

	if song.SongName != "" {
		query += " AND song_name = $" + strconv.Itoa(len(args)+1)
		args = append(args, song.SongName)
	}
	if song.GroupName != "" {
		query += " AND group_name = $" + strconv.Itoa(len(args)+1)
		args = append(args, song.GroupName)
	}
	query += " LIMIT $" + strconv.Itoa(len(args)+1)
	args = append(args, limit)
	query += " OFFSET $" + strconv.Itoa(len(args)+1)
	args = append(args, offset)

	err := r.DB.Select(&library, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return library, nil
}

func (r *Database) GetSongText(song *model.Song) (*string, error) {
	const op = "internal.storage.postgresql.GetSongText()"

	songId, err := r.foundSongId(song)
	if err != nil {
		return nil, err
	}

	id, err := r.foundSongDetailId(songId)
	if err != nil {
		return nil, err
	}

	var text string
	err = r.DB.QueryRow("SELECT text FROM songs_detail WHERE id = $1", id).Scan(&text)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return &text, nil
}

func (r *Database) AddSongDetail(song *model.Song, songDetail *model.SongDetail) error {
	const op = "internal.storage.postgresql.AddSongDetail()"

	songId, err := r.foundSongId(song)
	if err != nil {
		return err
	}

	if _, err := r.foundSongDetailId(songId); err == nil {
		err = r.PutSongDetail(song, songDetail)
		return err
	}

	_, err = r.DB.Exec("INSERT INTO songs_detail (release_date, link, text, song_id) VALUES ($1, $2, $3, $4)",
		songDetail.ReleaseDate, songDetail.Link, songDetail.Text, songId)

	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (r *Database) foundSongId(song *model.Song) (int64, error) {
	op := "internal.storage.postgresql.foundSongId()"
	var songId int64

	err := r.DB.QueryRow("SELECT id FROM songs WHERE song_name = $1 AND group_name = $2",
		song.SongName, song.GroupName).Scan(&songId)

	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("%s:%w", op, storage.ErrMusikNotFound)
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return songId, nil
}

func (r *Database) foundSongDetailId(songId int64) (int64, error) {
	op := "internal.storage.postgresql.foundSongDetailId()"
	var SongDetailId int64

	err := r.DB.QueryRow("SELECT id FROM songs_detail WHERE song_id = $1",
		songId).Scan(&SongDetailId)

	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("%s:%w", op, storage.ErrMusikDetailNotFound)
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return SongDetailId, nil
}
