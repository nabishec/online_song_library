package model

type Song struct {
	SongName  string `json:"song" validate:"required" db:"song_name"`
	GroupName string `json:"group" validate:"required" db:"group_name"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate" validate:"required" db:"release_date"`
	Link        string `json:"link" validate:"required" db:"link"`
	Text        string `json:"text" validate:"required" db:"text"`
}
