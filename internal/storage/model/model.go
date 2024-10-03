package model

type Song struct {
	SongName  string `json:"song" db:"song_name"`
	GroupName string `json:"group" db:"group_name"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate" db:"release_date"`
	Link        string `json:"link" db:"link"`
	Text        string `json:"text" db:"text"`
}
