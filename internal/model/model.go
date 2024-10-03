package model

type Song struct {
	SongName  string `json:"song" validate:"required,song" db:"song_name"`
	GroupName string `json:"group" validate:"required,group" db:"group_name"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate" validate:"required,releaseDate" db:"release_date"`
	Link        string `json:"link" validate:"required,link" db:"link"`
	Text        string `json:"text" validate:"required,text" db:"text"`
}
